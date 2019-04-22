<template>
    <div ref="vncContainer" style="height: 100%; width: 100%;" class="vnc-con"></div>
</template>

<script>
    import RFB from "../lib/vnc/core/rfb"
    import {init_logging} from "../lib/vnc/core/util/logging"

    export default {
        props: ['config'],
        data() {
            return {}
        },
        created() {
            window.debugObj.vnc = this
            init_logging('debug')
        },
        mounted() {
            const addr = `wss://${this.config.host}/${this.config.slave}/` + (this.config.viewOnly ? 'view' : 'operate')
            console.log(addr)
            let rfb = this.$data.$rfb = new RFB(this.$refs.vncContainer,
                addr, {credentials: {password: this.config.password}})
            rfb.scaleViewport = true
            rfb.addEventListener("connect", e => console.log('connect', e))
            rfb.addEventListener("disconnect", e => console.log('disconnect', e))
            rfb.addEventListener("credentialsrequired", e => console.log('credentialsrequired', e))
            rfb.addEventListener("desktopname", e => console.log('desktopname', e))
            rfb.addEventListener("fail", e => {
                console.log('fail', e)
                alert(e.detail)
            })
        },
        beforeDestroy() {
            this.$data.$rfb.close()
        }
    }
</script>

<style>
</style>
