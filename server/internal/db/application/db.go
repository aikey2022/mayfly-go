package application

import (
	"context"
	"fmt"
	"mayfly-go/internal/common/consts"
	"mayfly-go/internal/db/dbm"
	"mayfly-go/internal/db/dbm/dbi"
	"mayfly-go/internal/db/domain/entity"
	"mayfly-go/internal/db/domain/repository"
	tagapp "mayfly-go/internal/tag/application"
	"mayfly-go/pkg/base"
	"mayfly-go/pkg/biz"
	"mayfly-go/pkg/errorx"
	"mayfly-go/pkg/model"
	"mayfly-go/pkg/utils/collx"
	"mayfly-go/pkg/utils/stringx"
	"mayfly-go/pkg/utils/structx"
	"sort"
	"strings"
	"time"
)

type Db interface {
	base.App[*entity.Db]

	// 分页获取
	GetPageList(condition *entity.DbQuery, pageParam *model.PageParam, toEntity any, orderBy ...string) (*model.PageResult[any], error)

	Count(condition *entity.DbQuery) int64

	SaveDb(ctx context.Context, entity *entity.Db, tagIds ...uint64) error

	// 删除数据库信息
	Delete(ctx context.Context, id uint64) error

	// 获取数据库连接实例
	// @param id 数据库id
	//
	// @param dbName 数据库名
	GetDbConn(dbId uint64, dbName string) (*dbi.DbConn, error)

	// 根据数据库实例id获取连接，随机返回该instanceId下已连接的conn，若不存在则是使用该instanceId关联的db进行连接并返回。
	GetDbConnByInstanceId(instanceId uint64) (*dbi.DbConn, error)

	// DumpDb dumpDb
	DumpDb(ctx context.Context, reqParam *DumpDbReq) error
}

type dbAppImpl struct {
	base.AppImpl[*entity.Db, repository.Db]

	dbSqlRepo     repository.DbSql `inject:"DbSqlRepo"`
	dbInstanceApp Instance         `inject:"DbInstanceApp"`
	tagApp        tagapp.TagTree   `inject:"TagTreeApp"`
}

// 注入DbRepo
func (d *dbAppImpl) InjectDbRepo(repo repository.Db) {
	d.Repo = repo
}

// 分页获取数据库信息列表
func (d *dbAppImpl) GetPageList(condition *entity.DbQuery, pageParam *model.PageParam, toEntity any, orderBy ...string) (*model.PageResult[any], error) {
	return d.GetRepo().GetDbList(condition, pageParam, toEntity, orderBy...)
}

func (d *dbAppImpl) Count(condition *entity.DbQuery) int64 {
	return d.GetRepo().Count(condition)
}

func (d *dbAppImpl) SaveDb(ctx context.Context, dbEntity *entity.Db, tagIds ...uint64) error {
	// 查找是否存在
	oldDb := &entity.Db{Name: dbEntity.Name, InstanceId: dbEntity.InstanceId}
	err := d.GetBy(oldDb)

	if dbEntity.Id == 0 {
		if err == nil {
			return errorx.NewBiz("该实例下数据库名已存在")
		}

		resouceCode := stringx.Rand(16)
		dbEntity.Code = resouceCode

		return d.Tx(ctx, func(ctx context.Context) error {
			return d.Insert(ctx, dbEntity)
		}, func(ctx context.Context) error {
			return d.tagApp.RelateResource(ctx, resouceCode, consts.TagResourceTypeDb, tagIds)
		})
	}

	// 如果存在该库，则校验修改的库是否为该库
	if err == nil && oldDb.Id != dbEntity.Id {
		return errorx.NewBiz("该实例下数据库名已存在")
	}

	dbId := dbEntity.Id
	old, err := d.GetById(new(entity.Db), dbId)
	if err != nil {
		return errorx.NewBiz("该数据库不存在")
	}

	oldDbs := strings.Split(old.Database, " ")
	newDbs := strings.Split(dbEntity.Database, " ")
	// 比较新旧数据库列表，需要将移除的数据库相关联的信息删除
	_, delDb, _ := collx.ArrayCompare(newDbs, oldDbs)

	// 先简单关闭可能存在的旧库连接（可能改了关联标签导致DbConn.Info.TagPath与修改后的标签不一致、导致操作权限校验出错）
	for _, v := range oldDbs {
		// 关闭数据库连接
		dbm.CloseDb(dbEntity.Id, v)
	}

	for _, v := range delDb {
		// 删除该库关联的所有sql记录
		d.dbSqlRepo.DeleteByCond(ctx, &entity.DbSql{DbId: dbId, Db: v})
	}

	return d.Tx(ctx, func(ctx context.Context) error {
		return d.UpdateById(ctx, dbEntity)
	}, func(ctx context.Context) error {
		return d.tagApp.RelateResource(ctx, old.Code, consts.TagResourceTypeDb, tagIds)
	})
}

func (d *dbAppImpl) Delete(ctx context.Context, id uint64) error {
	db, err := d.GetById(new(entity.Db), id)
	if err != nil {
		return errorx.NewBiz("该数据库不存在")
	}
	dbs := strings.Split(db.Database, " ")
	for _, v := range dbs {
		// 关闭连接
		dbm.CloseDb(id, v)
	}

	return d.Tx(ctx,
		func(ctx context.Context) error {
			return d.DeleteById(ctx, id)
		},
		func(ctx context.Context) error {
			// 删除该库下用户保存的所有sql信息
			return d.dbSqlRepo.DeleteByCond(ctx, &entity.DbSql{DbId: id})
		}, func(ctx context.Context) error {
			var tagIds []uint64
			return d.tagApp.RelateResource(ctx, db.Code, consts.TagResourceTypeDb, tagIds)
		})
}

func (d *dbAppImpl) GetDbConn(dbId uint64, dbName string) (*dbi.DbConn, error) {
	return dbm.GetDbConn(dbId, dbName, func() (*dbi.DbInfo, error) {
		db, err := d.GetById(new(entity.Db), dbId)
		if err != nil {
			return nil, errorx.NewBiz("数据库信息不存在")
		}

		instance, err := d.dbInstanceApp.GetById(new(entity.DbInstance), db.InstanceId)
		if err != nil {
			return nil, errorx.NewBiz("数据库实例不存在")
		}

		// 密码解密
		if err := instance.PwdDecrypt(); err != nil {
			return nil, errorx.NewBiz(err.Error())
		}
		di := toDbInfo(instance, dbId, dbName, d.tagApp.ListTagPathByResource(consts.TagResourceTypeDb, db.Code)...)
		if db.FlowProcdefKey != nil {
			di.FlowProcdefKey = *db.FlowProcdefKey
		}

		checkDb := di.GetDatabase()
		if !strings.Contains(" "+db.Database+" ", " "+checkDb+" ") {
			return nil, errorx.NewBiz("未配置数据库【%s】的操作权限", dbName)
		}

		return di, nil
	})
}

func (d *dbAppImpl) GetDbConnByInstanceId(instanceId uint64) (*dbi.DbConn, error) {
	conn := dbm.GetDbConnByInstanceId(instanceId)
	if conn != nil {
		return conn, nil
	}

	var dbs []*entity.Db
	if err := d.ListByCond(&entity.Db{InstanceId: instanceId}, &dbs, "id", "database"); err != nil {
		return nil, errorx.NewBiz("获取数据库列表失败")
	}
	if len(dbs) == 0 {
		return nil, errorx.NewBiz("实例[%d]未配置数据库, 请先进行配置", instanceId)
	}

	// 使用该实例关联的已配置数据库中的第一个库进行连接并返回
	firstDb := dbs[0]
	return d.GetDbConn(firstDb.Id, strings.Split(firstDb.Database, " ")[0])
}

func (d *dbAppImpl) DumpDb(ctx context.Context, reqParam *DumpDbReq) error {
	writer := newGzipWriter(reqParam.Writer)
	defer writer.Close()
	dbId := reqParam.DbId
	dbName := reqParam.DbName
	tables := reqParam.Tables

	dbConn, err := d.GetDbConn(dbId, dbName)
	if err != nil {
		return err
	}
	writer.WriteString("\n-- ----------------------------")
	writer.WriteString("\n-- 导出平台: mayfly-go")
	writer.WriteString(fmt.Sprintf("\n-- 导出时间: %s ", time.Now().Format("2006-01-02 15:04:05")))
	writer.WriteString(fmt.Sprintf("\n-- 导出数据库: %s ", dbName))
	writer.WriteString("\n-- ----------------------------\n\n")

	dbMeta := dbConn.GetMetaData()
	if len(tables) == 0 {
		ti, err := dbMeta.GetTables()
		biz.ErrIsNil(err)
		tables = make([]string, len(ti))
		for i, table := range ti {
			tables[i] = table.TableName
		}
	}

	// 查询列信息，后面生成建表ddl和insert都需要列信息
	columns, err := dbMeta.GetColumns(tables...)
	biz.ErrIsNil(err)

	// 以表名分组，存放每个表的列信息
	columnMap := make(map[string][]dbi.Column)
	for _, column := range columns {
		columnMap[column.TableName] = append(columnMap[column.TableName], column)
	}

	// 按表名排序
	sort.Strings(tables)

	quoteSchema := dbMeta.QuoteIdentifier(dbConn.Info.CurrentSchema())
	dumpHelper := dbMeta.GetDumpHelper()
	dataHelper := dbMeta.GetDataHelper()

	// 遍历获取每个表的信息
	for _, tableName := range tables {
		quoteTableName := dbMeta.QuoteIdentifier(tableName)

		writer.TryFlush()
		// 查询表信息，主要是为了查询表注释
		tbs, err := dbMeta.GetTables(tableName)
		biz.ErrIsNil(err)
		if err != nil || tbs == nil || len(tbs) <= 0 {
			panic(errorx.NewBiz(fmt.Sprintf("获取表信息失败：%s", tableName)))
		}
		tabInfo := dbi.Table{
			TableName:    tableName,
			TableComment: tbs[0].TableComment,
		}

		// 生成表结构信息
		if reqParam.DumpDDL {
			writer.WriteString(fmt.Sprintf("\n-- ----------------------------\n-- 表结构: %s \n-- ----------------------------\n", tableName))
			tbDdlArr := dbMeta.GenerateTableDDL(columnMap[tableName], tabInfo, true)
			for _, ddl := range tbDdlArr {
				writer.WriteString(ddl + ";\n")
			}
		}

		// 生成insert sql，数据在索引前，加速insert
		if reqParam.DumpData {
			writer.WriteString(fmt.Sprintf("\n-- ----------------------------\n-- 表记录: %s \n-- ----------------------------\n", tableName))

			dumpHelper.BeforeInsert(writer, quoteTableName)
			// 获取列信息
			quoteColNames := make([]string, 0)
			for _, col := range columnMap[tableName] {
				quoteColNames = append(quoteColNames, dbMeta.QuoteIdentifier(col.ColumnName))
			}

			_ = dbConn.WalkTableRows(ctx, quoteTableName, func(row map[string]any, _ []*dbi.QueryColumn) error {
				rowValues := make([]string, len(columnMap[tableName]))
				for i, col := range columnMap[tableName] {
					rowValues[i] = dataHelper.WrapValue(row[col.ColumnName], dataHelper.GetDataType(string(col.DataType)))
				}

				beforeInsert := dumpHelper.BeforeInsertSql(quoteSchema, quoteTableName)
				insertSQL := fmt.Sprintf("%s INSERT INTO %s (%s) values(%s)", beforeInsert, quoteTableName, strings.Join(quoteColNames, ", "), strings.Join(rowValues, ", "))
				writer.WriteString(insertSQL + ";\n")
				return nil
			})

			dumpHelper.AfterInsert(writer, tableName, columnMap[tableName])
		}

		indexs, err := dbMeta.GetTableIndex(tableName)
		biz.ErrIsNil(err)

		// 过滤主键索引
		idxs := make([]dbi.Index, 0)
		for _, idx := range indexs {
			if !idx.IsPrimaryKey {
				idxs = append(idxs, idx)
			}
		}

		if len(idxs) > 0 {
			// 最后添加索引
			writer.WriteString(fmt.Sprintf("\n-- ----------------------------\n-- 表索引: %s \n-- ----------------------------\n", tableName))
			sqlArr := dbMeta.GenerateIndexDDL(idxs, tabInfo)
			for _, sqlStr := range sqlArr {
				writer.WriteString(sqlStr + ";\n")
			}
		}

	}

	return nil
}

func toDbInfo(instance *entity.DbInstance, dbId uint64, database string, tagPath ...string) *dbi.DbInfo {
	di := new(dbi.DbInfo)
	di.InstanceId = instance.Id
	di.Id = dbId
	di.Database = database
	di.TagPath = tagPath

	structx.Copy(di, instance)
	return di
}
