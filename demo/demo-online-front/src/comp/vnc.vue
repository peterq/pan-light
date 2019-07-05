<template>
    <div style="height: 100%; width: 100%; display: flex; flex-direction: column;">
        <div style="margin-bottom: 10px;">
            <el-button @click="clickGetTicket"
                       type="primary"
                       :disabled="$state.loading.getTicket" v-if="!$state.ticket"
                       :loading="$state.loading.getTicket">我也要体验
            </el-button>

            <el-button style="margin-right: 10px"
                       @click="$state.connectVnc = null"
                       type="danger"
                       :loading="$state.loading.getTicket">退出
            </el-button>
            <span v-show="!err">{{msg}}</span>
            <span style="color: red">{{err && ("出错了: " + err)}}</span>
        </div>
        <div ref="vncContainer" style="flex: 1;"></div>
    </div>
</template>

<script>
    import RFB from "../lib/vnc/core/rfb"
    import {init_logging} from "../lib/vnc/core/util/logging"
    import {getTicket, showError} from "../app"

    export default {
        props: ['config'],
        data() {
            return {
                msg: '远程桌面连接中...',
                err: ''
            }
        },
        created() {
            window.debugObj.vnc = this
            init_logging('debug')
            this.$rt.call('join.slave', {slave: this.config.slave})
            if (this.config.viewOnly) {
                this.config.password = 'peter.q.is.so.cool'
            }
        },
        methods: {
            async clickGetTicket() {
                await getTicket().catch(showError)
            },
        },
        mounted() {
            const addr = `wss://${this.config.host}/${this.config.slave}/` + (this.config.viewOnly ? 'view' : 'operate')
            console.log(addr)
            let rfb = this.$data.$rfb = new RFB(this.$refs.vncContainer,
                addr, {credentials: {password: this.config.password}})
            rfb.scaleViewport = true
            rfb.addEventListener("connect", e => this.msg = `远程桌面已连接, ${this.config.viewOnly ? '旁观' : '操作'}模式`)
            rfb.addEventListener("disconnect", e => console.log('disconnect', e))
            rfb.addEventListener("credentialsrequired", e => console.log('credentialsrequired', e))
            rfb.addEventListener("desktopname", e => console.log('desktopname', e))
            rfb.addEventListener("fail", e => {
                console.log('fail', e)
                this.err = e.detail
            })
        },
        beforeDestroy() {
            console.log(this.$data.$rfb)
            this.$data.$rfb.disconnect()
            this.$rt.call('leave.slave', {slave: this.config.slave})
        }
    }
</script>

<style>
</style>
