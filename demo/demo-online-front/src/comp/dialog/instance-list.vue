<template>
    <el-dialog title="实例列表" :visible.sync="visible" width="800px" @close="reject('close')" v-loading="loading">
        <el-table :data="host.slaves">
            <el-table-column prop="slaveName" label="名称" align="center"></el-table-column>
            <el-table-column prop="state" label="状态" align="center">
                <template slot-scope="scope">
                    {{$state.slaveStateMap[scope.row.state] || '未知'}}
                </template>
            </el-table-column>
            <el-table-column prop="visitorCount" label="在线人数" align="center"></el-table-column>
            <el-table-column prop="state" label="启动时间" align="center">
                <template slot-scope="scope">
                    {{scope.row.state === 'running' ? formatUnix(scope.row.startTime) : '-'}}
                </template>
            </el-table-column>
            <el-table-column prop="state" label="结束时间" align="center">
                <template slot-scope="scope">
                    {{scope.row.state === 'running' ? formatUnix(scope.row.endTime) : '-'}}
                </template>
            </el-table-column>
            <el-table-column prop="op" label="操作" align="center">
                <template slot-scope="scope">
                    <el-button type="primary"
                               @click="clickView(scope.row)"
                               :disabled="!['running', 'starting'].includes(scope.row.state)">围观</el-button>
                </template>
            </el-table-column>
        </el-table>
    </el-dialog>
</template>

<script>
    import moment from 'moment'
    export default {
        data: function () {
            return {
                hostName: '',
                loading: false,
                host: {
                    slaves: []
                }
            }
        },
        methods: {
            async onOpen(hostName) {
                this.hostName = hostName
                this.loading = true
                await this.getHostDetail().finally(() => this.loading = false)
                this.$watch(() => this.$state.timestamp, async t => {
                    if (t % 50) return
                    this.hosts = await this.getHostDetail()
                })
            },
            async getHostDetail() {
                let host = await this.$rt.call('host.detail', {hostName: this.hostName})
                host.slaves.sort((a, b) => (a.slaveName > b.slaveName) ? 1 : -1)
                this.host = host
            },
            async clickView(slave) {
                this.$state.connectVnc = {host: this.hostName, slave: slave.slaveName, viewOnly: true}
                this.resolve()
            },
            formatUnix(t) {
                return moment.unix(t).format('YYYY-MM-DD HH:mm:ss')
            }
        },
        components: {}
    }
</script>
<style lang="scss" scoped>
</style>
