<template>
    <el-container style="height: 100vh;width: 100vw;">
        <el-header>
            pan-light 在线体验
            <el-button @click="clickGetTicket" :loading="$state.loading.getTicket">立即体验</el-button>
        </el-header>
        <el-container>
            <el-main>
                <vnc style="flex: 1;" v-if="vncShow"
                     :config="connectVnc"></vnc>
                <host-list v-else></host-list>
            </el-main>
            <el-aside width="400px">Aside</el-aside>
        </el-container>
    </el-container>
</template>

<script>
    import {getTicket, showError} from "../app"
    import vnc from './vnc'
    import hostList from './hostList'

    const dataTemplate = {
        connectVnc: {
            host: '',
            slave: '',
            viewOnly: true
        }
    }
    export default {
        data() {
            return {
                connectVnc: null,
                vncShow: null,
            }
        },
        created() {
            this.$event.on('operate.turn', ({host, slave}) => {
                this.connectVnc = {
                    host, slave, viewOnly: false,
                    password: this.$state.ticket.ticket
                }
            })
        },
        methods: {
            async clickGetTicket() {
                await getTicket().catch(showError)
            },
        },
        watch: {
            async connectVnc(v) {
                if (!v) {
                    this.vncShow = false
                } else {
                    this.vncShow = false
                    if (v.viewOnly) {
                        this.connectVnc.password = 'peter.q.is.so.cool'
                    }
                    await this.$nextTick()
                    this.vncShow = true
                }
            }
        },
        components: {vnc, hostList}
    }
</script>

<style scoped>
    .el-aside {
        background-color: #D3DCE6;
    }

    .el-header, .el-footer {
        background-color: #B3C0D1;
        color: #333;
        text-align: center;
        line-height: 60px;
    }
</style>

