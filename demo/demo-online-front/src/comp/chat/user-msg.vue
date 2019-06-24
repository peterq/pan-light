<template>
    <div>
        <div v-if="!isSelf" style="display: flex; margin-bottom: 20px">
            <img :src="user.avatar" style="border-radius: 50%;" alt="avatar" width="50"
                 height="50">
            <div style="display: flex; flex-direction: column;margin-left: 10px;">
                <span style="font-size: small; color: #aaa">{{user.nickname}}</span>
                <span class="msg-content">{{message.msg}}</span>
            </div>
        </div>
        <div v-else style="display: flex; margin-bottom: 20px; flex-direction: row-reverse;">
            <img :src="user.avatar" style="border-radius: 50%;" alt="avatar" width="50"
                 height="50">
            <div style="display: flex; flex-direction: column;margin-right: 10px; align-items: flex-end">
                <span style="font-size: small; color: #aaa">{{user.nickname}}</span>
                <span class="msg-content" style="background: #00aa49;">{{message.msg}}</span>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        props: ['message'],
        data() {
            return {}
        },
        computed: {
            user() {
                return this.$state.userSessionInfo[this.message.from]
            },
            isSelf() {
                console.log({a:this.message.from, b:this.$state.userSessionInfo.self.sessionId})
                return this.message.from === this.$state.userSessionInfo.self.sessionId
            }
        }
    }
</script>

<style scoped>
    .msg-content {
        background: white;
        border: 1px solid #ddd;
        padding: 10px;
        border-radius: 10px;
        max-width: 260px;
        display: inline-block;
        word-break: break-word;
        word-wrap: break-word;
    }
</style>

