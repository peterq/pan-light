<template>
    <div id="app">
        <h1>Hello pan-light</h1>
    </div>
</template>

<script>

    import setRemoteDescription from "./realtime/webRtc"

    export default {
        data() {
            return {
                candidate: null,
            }
        },
        created() {
            this.$event.on('rtc.candidate', (candidate) =>  {
                this.candidate = candidate
                this.connectHost('asus-test')
            })
            this.$rt.onRemote('host.candidate.ok', c => {
                console.log(c)
                setRemoteDescription(c)
            })
        },
        methods: {
            async connectHost(host) {
                if (!this.candidate) return
                await this.$rt.openPromise
                let result = await this.$rt.call('connect.host', {
                    candidate: this.candidate,
                    hostName: host,
                    requestId: 'connect'
                })
                console.log(result)
            }
        }
    }
</script>

<style>
</style>
