<template>
    <el-container style="height: 100vh;width: 100vw;">
        <el-header>
            pan-light 在线体验
            <el-button style="margin-left: 10px;" @click="clickReturnHome">返回官网首页</el-button>
        </el-header>
        <el-container class="down-con">
            <el-main style="min-width: 800px">
                <vnc style="flex: 1;" v-if="vncShow"
                     :config="connectVnc"></vnc>
                <host-list v-else></host-list>
            </el-main>
            <el-aside width="400px">
                <div v-if="!$state.connected" style="height: 100%; display: flex; justify-content: center; align-items: center">
                    <p>聊天室初始化中...</p>
                </div>
                <chat-main v-else style="height: 100%"></chat-main>
            </el-aside>
        </el-container>
    </el-container>
</template>

<script>
    import vnc from './vnc'
    import hostList from './hostList'
    import {$state} from "../app"
    import chatMain from './chat/chat-main'

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
            clickReturnHome() {
                location.href = location.origin
            }
        },
        watch: {

        },
        components: {vnc, hostList, chatMain}
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
    .down-con {
        height: calc(100vh - 60px);
    }
</style>

