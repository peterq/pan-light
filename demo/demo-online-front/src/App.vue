<template>
    <div id="app" style="width: 100vw;height: 100vh; display: flex;flex-direction: column">
        <h1>Hello pan-light</h1>
        <vnc style="flex: 1;"></vnc>
    </div>
</template>

<script>
    import vnc from './comp/vnc'

    export default {
        data() {
            return {}
        },
        async created() {
        },
        methods: {},
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
