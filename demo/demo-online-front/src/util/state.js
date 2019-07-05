export const dataTemplate = {
    ticket: {
        order: 1,
        ticket: '23',
        inService: false
    },
    roomMap: {
        'room.user.all': {
            name: 'room.user.all',
            members: [],
            messages: [],
        }
    },
    connectVnc: {
        host: '',
        slave: '',
        viewOnly: true
    },
    userSessionInfo: {
        _role: 'user',
        nickname: '',
        avatar: '',
        sessionId: ''
    },
    deppClone(key) {
        return JSON.parse(JSON.stringify(this[key]))
    }
}

function initialData() {
    return {
        connected: false,
        loading: {
            getTicket: false,
        },
        userSessionInfo: {
            self: {...dataTemplate.userSessionInfo},
        },
        ticket: null, // tpl: ticket
        roomMap: {}, // tpl roomMap
        timestamp: 0,
        hosts: [],
        slaveStateMap: {
            wait: '空闲',
            running: '运行中',
            starting: '启动中',
        },
        connectVnc: null
    }
}

export default {
    data() {
        return initialData()
    },
    created() {
        window.debugObj.$state = this
        setInterval(() => {
            this.timestamp++
        }, 100)
    },
    methods: {
        resetData() {
            let data = initialData()
            for (let k in data) {
                this[k] = data[k]
            }
        }
    }
}
