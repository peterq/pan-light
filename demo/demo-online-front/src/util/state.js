const dataTemplate = {
    ticket: {
        order: 1,
        ticket: '23',
        inService: false
    },
    roomMap: {
        'room.user.all': {
            name: 'room.user.all',
            members: ['234']
        }
    }
}

function data() {
    return {
        loading: {
            getTicket: false,
        },
        ticket: null, // tpl: ticket
        roomMap: {} // tpl roomMap
    }
}

export default {
    data() {
        return data()
    },
    created() {
        window.debugObj.$state = this
    }
}
