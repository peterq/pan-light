<template>
    <div id="app" style="width: 100vw;height: 100vh; display: flex;flex-direction: column">
        <div style="flex: 1;">
            <div v-for="(h) in hosts" :key="h.name"
                    style="margin: 50px;
            background: rgba(0, 0, 0, .8);
            width: 200px; height: 200px; color: white;
             display: flex; align-items: center;
              justify-content: center; border-radius: 10px">
                <div style="text-align: center">
                    <p>主机: {{h.name}}</p>
                    <p>运行实例: {{h.slaves.length}}</p>
                </div>
            </div>
        </div>
        <vnc style="flex: 1;" v-if="connectVnc"
             :host="connectVnc.host" :slave="connectVnc.slave"
             :view-only="connectVnc.viewOnly"></vnc>
    </div>
</template>

<script>
    import vnc from './comp/vnc'

    export default {
        data() {
            return {
                hosts: [],
                connectVnc: null
            }
        },
        async created() {
            window.debugObj.app = this
            await this.$rt.openPromise
            this.hosts = await this.getHosts()
            console.log(this.hosts )
        },
        methods: {
            async getHosts() {
                return await this.$rt.call('hosts.info')
            }
        },
        components: {
            vnc
        }
    }
</script>

<style>
    * {
        padding: 0;
        margin: 0;
    }
</style>
