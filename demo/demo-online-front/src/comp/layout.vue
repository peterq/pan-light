<template>
    <el-container style="height: 100vh;width: 100vw;">
        <el-header>
            pan-light 在线体验
            <el-button @click="clickGetTicket" :loading="$state.loading.getTicket">立即体验</el-button>
        </el-header>
        <el-container>
            <el-main>
                <vnc style="flex: 1;" v-if="connectVnc"
                     :host="connectVnc.host" :slave="connectVnc.slave"
                     :view-only="connectVnc.viewOnly"></vnc>
            </el-main>
            <el-aside width="400px">Aside</el-aside>
        </el-container>
    </el-container>
</template>

<script>
    import {getTicket, showError} from "../app"

    export default {
        data() {
            return {
                connectVnc: null
            }
        },
        created() {
        },
        methods: {
            async clickGetTicket() {
                await getTicket().catch(showError)
            },
        },
        components: {}
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

