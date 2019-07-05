<template>
    <div style="display: flex; flex-direction: column;">
        <h2 style="text-align: center">在线聊天</h2>
        <hr>
        <div style="display: flex; align-items: center; padding: 10px 20px;">
            <img :src="$state.userSessionInfo.self.avatar" style="border-radius: 50%;" alt="avatar" width="50"
                 height="50"/>
            <p> {{$state.userSessionInfo.self.nickname}}</p>
            <p style="margin: 20px" v-if="$state.ticket">
                <i class="el-icon-time"></i>
                体验门票号码: {{$state.ticket.order}}
            </p>
        </div>
        <hr>
        <el-tabs v-model="activeRoom" style="padding: 5px; flex: 1;">
            <el-tab-pane v-for="(room, roomName) in $state.roomMap" :key="roomName" :label="getRoomLabel(roomName)"
                         :name="roomName">
                <room :room="room" style="height: 100%"></room>
            </el-tab-pane>
        </el-tabs>
    </div>
</template>

<script>
    import room from './room'

    export default {
        data() {
            return {
                activeRoom: 'room.all.user'
            }
        },
        computed: {},
        components: {
            room
        },
        methods: {
            getRoomLabel(name) {
                if (name === 'room.all.user') {
                    return '全员群'
                }
                return '围观群 ' + name.replace('room.slave.all.user.', '')
            }
        }

    }
</script>

<style>
    .el-tabs__content {
        flex: 1;
    }
    .el-tabs {
        display: flex;
        flex-direction: column;
    }
    .el-tab-pane {
        height: 100%;
    }
</style>

