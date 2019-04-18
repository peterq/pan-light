<template>
    <layout></layout>
</template>

<script>
    import layout from './comp/layout'

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
            layout
        }
    }
</script>

<style>
    * {
        padding: 0;
        margin: 0;
    }
</style>
