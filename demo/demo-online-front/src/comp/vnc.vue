<template>
    <div ref="vncContainer" class="vnc-con"></div>
</template>

<script>
    import RFB from "../lib/vnc/core/rfb"
    import {init_logging} from "../lib/vnc/core/util/logging"

    export default {
        props: ['host', 'slave', 'viewOnly'],
        data() {
            return {
                rfb: null
            }
        },
        created() {
            window.debugObj.vnc = this
            init_logging('debug')
        },
        mounted() {
            let rfb = this.rfb = new RFB(this.$refs.vncContainer,
                `wss://${this.host}/${this.slave}/` + this.viewOnly ? 'view' : 'operate', {credentials: {password: ''}})
            rfb.addEventListener("connect", e => console.log('connect', e))
            rfb.addEventListener("disconnect", e => console.log('disconnect', e))
            rfb.addEventListener("credentialsrequired", e => console.log('credentialsrequired', e))
            rfb.addEventListener("desktopname", e => console.log('desktopname', e))
            rfb.addEventListener("fail", e => {
                console.log('fail', e)
                alert(e.detail)
            })
        }
    }
</script>

<style>
</style>
