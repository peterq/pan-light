<template>
    <div style="margin: 20px;">
        <div style="text-align: center;" v-if="message.evt === 'join'">
            <span class="system-tip"><span style="color: #00a3fb;">{{user.nickname}}</span> 加入了群聊</span>
        </div>
        <div style="text-align: center;" v-else-if="message.evt === 'leave'">
            <span class="system-tip"><span style="color: #00a3fb;">{{user.nickname}}</span> 退出了群聊</span>
        </div>
        <div style="text-align: center;" v-else-if="message.evt === 'turn'">
            <span class="system-tip">请 {{message.ticket.order}} 号 <span style="color: #00a3fb;">{{user.nickname}}</span> 操作 {{message.ticket.slave}}</span>
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
                if (!this.message.sessionId) return {}
                return this.$state.userSessionInfo[this.message.sessionId]
            }
        }
    }
</script>

<style scoped>
    .system-tip {
        padding: 5px;
        border-radius: 3px;
        background: rgba(255, 255, 255, .2);
        font-size: small;
        color: #aaa;
    }
</style>

