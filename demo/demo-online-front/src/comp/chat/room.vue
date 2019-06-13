<template>
    <div style="display: flex; flex-direction: column;">
        <el-scrollbar ref="scroll" class="msg-con" style="height: 100%;"
                      :native="false"
                      tag="section">
            <div v-for="(item) in room.messages" :key="item.id">
                <system-msg v-if="item.type === 'system'" :message="item"></system-msg>
                <user-msg v-else-if="item.type==='chat'" :message="item"></user-msg>
            </div>
        </el-scrollbar>
        <div style="height: 80px; display: flex;">
            <el-input
                    style="height: 100%"
                    type="textarea"
                    placeholder="在此输入消息"
                    v-model="inputMsg">
            </el-input>
            <el-button type="primary" @click="clickSend" style="height: 100%; border-radius: 0">
                发 送
            </el-button>
        </div>
    </div>
</template>

<script>

    import UserMsg from "./user-msg"
    import SystemMsg from "./system-msg"
    export default {
        props: ['room'],
        data() {
            return {
                inputMsg: 'hello world'
            }
        },
        mounted() {
            console.log(this.$refs.scroll)
        },
        computed: {},
        methods: {
            clickSend() {
                if (!this.inputMsg) return
                this.room.sendMsg(this.inputMsg)
                this.inputMsg = ''
            }
        },
        components: {SystemMsg, UserMsg},
        watch: {
            async ['room.messages']() {
                await this.$nextTick()
                let div = this.$refs.scroll.$el.querySelector('.el-scrollbar__wrap')
                div.scrollTop = div.scrollHeight
                this.$refs.scroll.update()
            }
        }

    }
</script>

<style>
    .el-textarea__inner {
        height: 100%;
    }

    .msg-con {
        height: calc(100% - 80px);
    }

    .el-scrollbar__wrap {
        overflow-x: hidden !important;
    }
</style>

