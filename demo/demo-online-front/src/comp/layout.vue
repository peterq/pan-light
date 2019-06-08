<template>
    <el-container style="height: 100vh;width: 100vw;">
        <el-header>
            pan-light 在线体验
        </el-header>
        <el-container>
            <el-main style="min-width: 800px">
                <vnc style="flex: 1;" v-if="vncShow"
                     :config="connectVnc"></vnc>
                <host-list v-else></host-list>
            </el-main>
            <el-aside width="400px">Aside</el-aside>
        </el-container>
    </el-container>
</template>

<script>
    import vnc from './vnc'
    import hostList from './hostList'
    import {$state} from "../app"

    export default {
        data() {
            return {
                connectVnc: null,
                vncShow: null,
            }
        },
        created() {
            this.$watch(() => $state.connectVnc, async (v) => {
                this.connectVnc = v
                if (!v) {
                    this.vncShow = false
                } else {
                    this.vncShow = false
                    await this.$nextTick()
                    this.vncShow = true
                }
            })
        },
        methods: {
        },
        watch: {

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

