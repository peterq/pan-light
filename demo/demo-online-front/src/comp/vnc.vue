<template>
    <div ref="vncContainer" class="vnc-con"></div>
</template>

<script>
    import RFB from "../lib/vnc/core/rfb"

    export default {
        data() {
            return {
                rfb: null
            }
        },
        created() {
            window.debugObj.vnc = this
        },
        mounted() {
            let rfb = this.rfb = new RFB(this.$refs.vncContainer, 'wss://asus-test/asus-test.slave.0/view')
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
