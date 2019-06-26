<template>
    <div style="flex: 1; flex-wrap: wrap;">
        <el-button v-show="hosts.length > 0" @click="clickGetTicket"
                   type="primary"
                   :disabled="$state.loading.getTicket" v-if="!$state.ticket"
                   :loading="$state.loading.getTicket">立即体验
        </el-button>
        <div style="display: flex">
            <div v-for="(h) in hosts" :key="h.name"
                 style="margin: 50px;
            background: rgba(0, 0, 0, .8);
            width: 200px; height: 200px; color: white;
             display: flex; align-items: center;
              justify-content: center; border-radius: 10px">
                <div style="text-align: center">
                    <p>主机: {{h.name}}</p>
                    <p>运行实例: {{h.slaves.length}}</p>
                    <br>
                    <el-button type="primary" @click="clickView(h)">查看</el-button>
                </div>
            </div>
        </div>
        <el-alert class="no-host" v-if="$state.connected && hosts.length === 0" style="font-size: 30px"
                title="抱歉, 当前没有主机可提供远程演示服务, 请联系管理员至少开启一台主机"
                type="warning">
        </el-alert>
    </div>
</template>

<script>
    import {openDialog} from "../util/dialogManger"
    import InstanceList from "./dialog/instance-list"
    import {getTicket, showError} from "../app"

    export default {
        data() {
            return {
                hosts: []
            }
        },
        async created() {
            await this.$rt.openPromise
            this.hosts = await this.getHosts()
            console.log(this.hosts)
            console.log(this, this.$state.timestamp)
            this.$watch(() => this.$state.timestamp, async t => {
                if (t % 50) return
                this.hosts = await this.getHosts()
            })
        },
        methods: {
            async getHosts() {
                let hosts = await this.$rt.call('hosts.info')
                this.$state.hosts = hosts
                return hosts || []
            },
            async clickView(host) {
                await openDialog(InstanceList, host.name).getPromise()
            },
            async clickGetTicket() {
                await getTicket().catch(showError)
            },
        },
    }
</script>

<style >
    * {
        padding: 0;
        margin: 0;
    }
    .no-host .el-alert__title {
        font-size: 16px;
    }
</style>
