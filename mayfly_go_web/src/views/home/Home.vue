<template>
    <div class="home-container personal">
        <el-row :gutter="15">
            <!-- 个人信息 -->
            <el-col :xs="24" :sm="16">
                <el-card shadow="hover" header="个人信息">
                    <div class="personal-user">
                        <div class="personal-user-left">
                            <el-upload class="h100 personal-user-left-upload" action="" multiple :limit="1">
                                <img :src="userInfo.photo" />
                            </el-upload>
                        </div>
                        <div class="personal-user-right">
                            <el-row>
                                <el-col :span="24" class="personal-title mb18"
                                    >{{ currentTime }}，{{ userInfo.name }}，生活变的再糟糕，也不妨碍我变得更好！
                                </el-col>
                                <el-col :span="24">
                                    <el-row>
                                        <el-col :xs="24" :sm="12" class="personal-item mb6">
                                            <div class="personal-item-label">用户名：</div>
                                            <div class="personal-item-value">{{ userInfo.username }}</div>
                                        </el-col>
                                        <el-col :xs="24" :sm="12" class="personal-item mb6">
                                            <div class="personal-item-label">角色：</div>
                                            <div class="personal-item-value">{{ roleInfo }}</div>
                                        </el-col>
                                    </el-row>
                                </el-col>
                                <el-col :span="24">
                                    <el-row>
                                        <el-col :xs="24" :sm="12" class="personal-item mb6">
                                            <div class="personal-item-label">上次登录IP：</div>
                                            <div class="personal-item-value">{{ userInfo.lastLoginIp }}</div>
                                        </el-col>
                                        <el-col :xs="24" :sm="12" class="personal-item mb6">
                                            <div class="personal-item-label">上次登录时间：</div>
                                            <div class="personal-item-value">{{ formatDate(userInfo.lastLoginTime) }}</div>
                                        </el-col>
                                    </el-row>
                                </el-col>
                            </el-row>
                        </div>
                    </div>
                </el-card>
            </el-col>

            <!-- 消息通知 -->
            <el-col :xs="24" :sm="8" class="pl15 personal-info">
                <el-card shadow="hover">
                    <template #header>
                        <span>消息通知</span>
                        <span @click="showMsgs" class="personal-info-more">更多</span>
                    </template>
                    <div class="personal-info-box">
                        <ul class="personal-info-ul">
                            <li v-for="(v, k) in state.msgs as any" :key="k" class="personal-info-li">
                                <a class="personal-info-li-title">{{ `[${getMsgTypeDesc(v.type)}] ${v.msg}` }}</a>
                            </li>
                        </ul>
                    </div>
                </el-card>
            </el-col>
        </el-row>

        <el-row :gutter="20" class="mt20 resource-info">
            <el-col :sm="12">
                <el-card shadow="hover">
                    <template #header>
                        <el-row justify="center">
                            <div class="resource-num pointer-icon" @click="toPage('machine')">
                                <SvgIcon
                                    class="mb5 mr5"
                                    :size="28"
                                    :name="TagResourceTypeEnum.Machine.extra.icon"
                                    :color="TagResourceTypeEnum.Machine.extra.iconColor"
                                />
                                <span class="">{{ state.machine.num }}</span>
                            </div>
                        </el-row>
                    </template>
                    <el-row>
                        <el-col :sm="24">
                            <el-table :data="state.machine.opLogs" :height="state.resourceOpTableHeight" stripe size="small" empty-text="暂无操作记录">
                                <el-table-column prop="createTime" show-overflow-tooltip width="135">
                                    <template #default="scope">
                                        {{ formatDate(scope.row.createTime) }}
                                    </template>
                                </el-table-column>
                                <el-table-column prop="codePath" min-width="400" show-overflow-tooltip>
                                    <template #default="scope">
                                        <TagCodePath :path="scope.row.codePath" />
                                    </template>
                                </el-table-column>
                                <el-table-column width="30">
                                    <template #default="scope">
                                        <el-link @click="toPage('machine', scope.row.codePath)" type="primary" icon="Position"></el-link>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-col>
                    </el-row>
                </el-card>
            </el-col>

            <el-col :sm="12">
                <el-card shadow="hover">
                    <template #header>
                        <el-row justify="center">
                            <div class="resource-num pointer-icon" @click="toPage('db')">
                                <SvgIcon class="mb5 mr5" :size="28" :name="TagResourceTypeEnum.Db.extra.icon" :color="TagResourceTypeEnum.Db.extra.iconColor" />
                                <span class="">{{ state.db.num }}</span>
                            </div>
                        </el-row>
                    </template>
                    <el-row>
                        <el-col :sm="24">
                            <el-table :data="state.db.opLogs" :height="state.resourceOpTableHeight" stripe size="small" empty-text="暂无操作记录">
                                <el-table-column prop="createTime" show-overflow-tooltip min-width="135">
                                    <template #default="scope">
                                        {{ formatDate(scope.row.createTime) }}
                                    </template>
                                </el-table-column>
                                <el-table-column prop="codePath" min-width="380" show-overflow-tooltip>
                                    <template #default="scope">
                                        <TagCodePath :path="scope.row.codePath" />
                                    </template>
                                </el-table-column>
                                <el-table-column width="30">
                                    <template #default="scope">
                                        <el-link @click="toPage('db', scope.row.codePath)" type="primary" icon="Position"></el-link>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-col>
                    </el-row>
                </el-card>
            </el-col>
        </el-row>

        <el-row :gutter="20" class="mt20 resource-info">
            <el-col :sm="12">
                <el-card shadow="hover">
                    <template #header>
                        <el-row justify="center">
                            <div class="resource-num pointer-icon" @click="toPage('redis')">
                                <SvgIcon
                                    class="mb5 mr5"
                                    :size="28"
                                    :name="TagResourceTypeEnum.Redis.extra.icon"
                                    :color="TagResourceTypeEnum.Redis.extra.iconColor"
                                />
                                <span class="">{{ state.redis.num }}</span>
                            </div>
                        </el-row>
                    </template>
                    <el-row>
                        <el-col :sm="24">
                            <el-table :data="state.redis.opLogs" :height="state.resourceOpTableHeight" stripe size="small" empty-text="暂无操作记录">
                                <el-table-column prop="createTime" show-overflow-tooltip min-width="135">
                                    <template #default="scope">
                                        {{ formatDate(scope.row.createTime) }}
                                    </template>
                                </el-table-column>
                                <el-table-column prop="codePath" min-width="380" show-overflow-tooltip>
                                    <template #default="scope">
                                        <TagCodePath :path="scope.row.codePath" />
                                    </template>
                                </el-table-column>
                                <el-table-column width="30">
                                    <template #default="scope">
                                        <el-link @click="toPage('redis', scope.row.codePath)" type="primary" icon="Position"></el-link>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-col>
                    </el-row>
                </el-card>
            </el-col>

            <el-col :sm="12">
                <el-card shadow="hover">
                    <template #header>
                        <el-row justify="center">
                            <div class="resource-num pointer-icon" @click="toPage('mongo')">
                                <SvgIcon
                                    class="mb5 mr5"
                                    :size="28"
                                    :name="TagResourceTypeEnum.Mongo.extra.icon"
                                    :color="TagResourceTypeEnum.Mongo.extra.iconColor"
                                />
                                <span class="">{{ state.mongo.num }}</span>
                            </div>
                        </el-row>
                    </template>
                    <el-row>
                        <el-col :sm="24">
                            <el-table :data="state.mongo.opLogs" :height="state.resourceOpTableHeight" stripe size="small" empty-text="暂无操作记录">
                                <el-table-column prop="createTime" show-overflow-tooltip min-width="135">
                                    <template #default="scope">
                                        {{ formatDate(scope.row.createTime) }}
                                    </template>
                                </el-table-column>
                                <el-table-column prop="codePath" min-width="380" show-overflow-tooltip>
                                    <template #default="scope">
                                        <TagCodePath :path="scope.row.codePath" />
                                    </template>
                                </el-table-column>
                                <el-table-column width="30">
                                    <template #default="scope">
                                        <el-link @click="toPage('mongo', scope.row.codePath)" type="primary" icon="Position"></el-link>
                                    </template>
                                </el-table-column>
                            </el-table>
                        </el-col>
                    </el-row>
                </el-card>
            </el-col>
        </el-row>

        <el-dialog width="900px" title="消息" v-model="msgDialog.visible">
            <el-table border :data="msgDialog.msgs.list" size="small">
                <el-table-column property="type" label="类型" width="60">
                    <template #default="scope">
                        {{ getMsgTypeDesc(scope.row.type) }}
                    </template>
                </el-table-column>
                <el-table-column property="msg" label="消息"></el-table-column>
                <el-table-column property="createTime" label="时间" width="150">
                    <template #default="scope">
                        {{ formatDate(scope.row.createTime) }}
                    </template>
                </el-table-column>
            </el-table>
            <el-row type="flex" class="mt5" justify="center">
                <el-pagination
                    small
                    @current-change="searchMsg"
                    style="text-align: center"
                    background
                    layout="prev, pager, next, total, jumper"
                    :total="msgDialog.msgs.total"
                    v-model:current-page="msgDialog.query.pageNum"
                    :page-size="msgDialog.query.pageSize"
                />
            </el-row>
        </el-dialog>
    </div>
</template>

<script lang="ts" setup>
import { toRefs, reactive, onMounted, computed } from 'vue';
// import * as echarts from 'echarts';
import { formatAxis } from '@/common/utils/format';
import { indexApi } from './api';
import { useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useUserInfo } from '@/store/userInfo';
import { personApi } from '../personal/api';
import { formatDate } from '@/common/utils/format';
import SvgIcon from '@/components/svgIcon/index.vue';
import { TagResourceTypeEnum } from '@/common/commonEnum';
import { resourceOpLogApi } from '../ops/tag/api';
import TagCodePath from '../ops/component/TagCodePath.vue';
import { useAutoOpenResource } from '@/store/autoOpenResource';

const router = useRouter();
const { userInfo } = storeToRefs(useUserInfo());

const state = reactive({
    accountInfo: {
        roles: [],
    },
    msgs: [],
    msgDialog: {
        visible: false,
        query: {
            pageSize: 10,
            pageNum: 1,
        },
        msgs: {
            list: [],
            total: null,
        },
    },
    resourceOpTableHeight: 180,
    defaultLogSize: 5,
    machine: {
        num: 0,
        opLogs: [],
    },
    db: {
        num: 0,
        opLogs: [],
    },
    redis: {
        num: 0,
        opLogs: [],
    },
    mongo: {
        num: 0,
        opLogs: [],
    },
});

const { msgDialog } = toRefs(state);

const roleInfo = computed(() => {
    if (state.accountInfo.roles.length == 0) {
        return '';
    }
    return state.accountInfo.roles.map((val: any) => val.roleName).join('、');
});

// 当前时间提示语
const currentTime = computed(() => {
    return formatAxis(new Date());
});

// 页面加载时
onMounted(() => {
    initData();
    getAccountInfo();

    getMsgs().then((res) => {
        state.msgs = res.list;
    });
});

const showMsgs = async () => {
    state.msgDialog.query.pageNum = 1;
    searchMsg();
    state.msgDialog.visible = true;
};

const searchMsg = async () => {
    state.msgDialog.msgs = await getMsgs();
};

const getMsgTypeDesc = (type: number) => {
    if (type == 1) {
        return '登录';
    }
    if (type == 2) {
        return '通知';
    }
};

const getAccountInfo = async () => {
    state.accountInfo = await personApi.accountInfo.request();
};

const getMsgs = async () => {
    return await personApi.getMsgs.request(state.msgDialog.query);
};

// 初始化数字滚动
const initData = async () => {
    resourceOpLogApi.getAccountResourceOpLogs
        .request({ resourceType: TagResourceTypeEnum.MachineAuthCert.value, pageSize: state.defaultLogSize })
        .then((res: any) => {
            state.machine.opLogs = res.list;
        });

    resourceOpLogApi.getAccountResourceOpLogs.request({ resourceType: TagResourceTypeEnum.DbName.value, pageSize: state.defaultLogSize }).then((res: any) => {
        state.db.opLogs = res.list;
    });

    resourceOpLogApi.getAccountResourceOpLogs.request({ resourceType: TagResourceTypeEnum.Redis.value, pageSize: state.defaultLogSize }).then((res: any) => {
        state.redis.opLogs = res.list;
    });

    resourceOpLogApi.getAccountResourceOpLogs.request({ resourceType: TagResourceTypeEnum.Mongo.value, pageSize: state.defaultLogSize }).then((res: any) => {
        state.mongo.opLogs = res.list;
    });

    indexApi.machineDashbord.request().then((res: any) => {
        state.machine.num = res.machineNum;
    });

    indexApi.dbDashbord.request().then((res: any) => {
        state.db.num = res.dbNum;
    });

    indexApi.redisDashbord.request().then((res: any) => {
        state.redis.num = res.redisNum;
    });

    indexApi.mongoDashbord.request().then((res: any) => {
        state.mongo.num = res.mongoNum;
    });
};

const toPage = (item: any, codePath = '') => {
    let path;
    switch (item) {
        case 'personal': {
            router.push('/personal');
            break;
        }
        case 'mongo': {
            useAutoOpenResource().setMongoCodePath(codePath);
            path = '/mongo/mongo-data-operation';
            break;
        }
        case 'machine': {
            useAutoOpenResource().setMachineCodePath(codePath);
            path = '/machine/machines-op';
            break;
        }
        case 'db': {
            useAutoOpenResource().setDbCodePath(codePath);
            path = '/dbms/sql-exec';
            break;
        }
        case 'redis': {
            useAutoOpenResource().setRedisCodePath(codePath);
            path = '/redis/data-operation';
            break;
        }
    }

    router.push({ path });
};
</script>

<style scoped lang="scss">
@import '@/theme/mixins/index.scss';

.personal {
    .personal-user {
        height: 130px;
        display: flex;
        align-items: center;

        .personal-user-left {
            width: 100px;
            height: 130px;
            border-radius: 3px;

            ::v-deep(.el-upload) {
                height: 100%;
            }

            .personal-user-left-upload {
                img {
                    width: 100%;
                    height: 100%;
                    border-radius: 3px;
                }

                &:hover {
                    img {
                        animation: logoAnimation 0.3s ease-in-out;
                    }
                }
            }
        }

        .personal-user-right {
            flex: 1;
            padding: 0 15px;

            .personal-title {
                font-size: 18px;
                @include text-ellipsis(1);
            }

            .personal-item {
                display: flex;
                align-items: center;
                font-size: 13px;

                .personal-item-label {
                    color: gray;
                    @include text-ellipsis(1);
                }

                .personal-item-value {
                    @include text-ellipsis(1);
                }
            }
        }
    }

    .personal-info {
        .personal-info-more {
            float: right;
            color: gray;
            font-size: 13px;

            &:hover {
                color: var(--el-color-primary);
                cursor: pointer;
            }
        }

        .personal-info-box {
            height: 130px;
            overflow: hidden;

            .personal-info-ul {
                list-style: none;

                .personal-info-li {
                    font-size: 13px;
                    padding-bottom: 10px;

                    .personal-info-li-title {
                        display: inline-block;
                        @include text-ellipsis(1);
                        color: grey;
                        text-decoration: none;
                    }

                    & a:hover {
                        color: var(--el-color-primary);
                        cursor: pointer;
                    }
                }
            }
        }
    }
}

.resource-info {
    text-align: center;

    ::v-deep(.el-card__header) {
        padding: 2px 20px;
    }

    .resource-num {
        font-weight: 700;
        font-size: 2vw;
    }
}

.home-container {
    overflow-x: hidden;

    .home-card-item {
        width: 100%;
        height: 103px;
        background: gray;
        border-radius: 4px;
        transition: all ease 0.3s;
        cursor: pointer;

        &:hover {
            box-shadow: 0 2px 12px 0 rgb(0 0 0 / 10%);
            transition: all ease 0.3s;
        }
    }

    .home-card-item-box {
        display: flex;
        align-items: center;
        position: relative;
        overflow: hidden;

        &:hover {
            i {
                right: 0px !important;
                bottom: 0px !important;
                transition: all ease 0.3s;
            }
        }

        i {
            position: absolute;
            right: -10px;
            bottom: -10px;
            font-size: 70px;
            transform: rotate(-30deg);
            transition: all ease 0.3s;
        }

        .home-card-item-flex {
            padding: 0 20px;
            color: white;

            .home-card-item-title,
            .home-card-item-tip {
                font-size: 13px;
            }

            .home-card-item-title-num {
                font-size: 2vw;
            }

            .home-card-item-tip-num {
                font-size: 13px;
            }
        }
    }
}
</style>
