import {newPeerConnection} from "./realtime/webRtc"
import RealTime from "./realtime/realtime"
import Vue from "vue"
import State from "./util/state"
import {registerProxyChannelResolver} from "./lib/vnc/core/RtcWebSocket"

// const $rt = new RealTime('ws://localhost:8001/demo/ws')
export const $rt = new RealTime((location.protocol === 'https:' ? 'wss' : 'ws') + '://' + location.host + '/demo/ws')
export const $state = new Vue(State)
export const $event = (function () {
    function fire(evt, payload) {
        map.has(evt) && map.get(evt).forEach(fn => {
            setTimeout(() => {
                fn(payload)
            }, 0)
        })
    }

    function on(event, handler) {
        map.has(event) || map.set(event, new Set())
        map.get(event).add(handler)
    }

    function off(event, handler) {
        map.has(event) && map.get(event).delete(handler)
    }

    const map = new Map()
    return {fire, on, off}
})()

registerProxyChannelResolver(async function (uri) {
    uri = uri.replace('ws://', '').replace('wss://', '')
    let [hostName, slave, method] = uri.split('/')
    let host = await connectHost(hostName)
    console.log(host)
    return await host.vncProxyChanel(slave, method === 'view')
})

const connectionRequestMap = {}


$rt.onRemote("host.candidate.ok", data => {
    const id = data.requestId
    const candidate = data.candidate
    const handler = connectionRequestMap[id]
    handler.pc.continueWithRemote(candidate)
})
$rt.onRemote('room.member.join', (sid, room) => {
    $state.roomMap[room] || Vue.set($state.roomMap, room, {
        name: room,
        members: new Set()
    })
    if (sid === $rt.sessionId) {
        // todo get members
    } else {
        $state.roomMap[room].members.add(sid)
    }
})

$rt.onRemote('room.member.join', (sid, room) => {
    $state.roomMap[room] = $state.roomMap[room] || {
        name: room,
        members: new Set()
    }
    if (sid === $rt.sessionId) {
        // todo get members
    } else {
        $state.roomMap[room].members.add(sid)
    }
})

$rt.onRemote('ticket.turn', data => {
    console.log(data, $state.ticket)
    const {order} = data
    if ($state.ticket && $state.ticket.order === order) {
        $state.ticket.inService = true
        $event.fire('operate.turn', data)
        console.log(data)
    }
})

class Host {
    /**
     * @type  RTCDataChannel
     */
    infoChannel

    /**
     * @type RTCPeerConnection
     */
    pc

    channelMap
    channelWaitMap
    proxyWaitMap

    static hostMap = {}

    constructor(name, infoChannel, pc) {
        if (Host.hostMap[name])
            throw new Error(`the host ${name} is already exist`)
        Host.hostMap[name] = this
        this.name = name
        this.pc = pc
        this.pc.ondatachannel = this._onDataChannel.bind(this)
        this.infoChannel = infoChannel
        this.infoChannel.onmessage = this._onInfoMessage.bind(this)
        this.callHostMap = {}
        this.channelMap = {}
        this.channelWaitMap = {}
        this.proxyWaitMap = {}
        this.textDecoder = new TextDecoder("utf-8")
    }

    _onDataChannel(evt) {
        const channel = evt.channel
        // console.log(channel)
        this.channelMap[channel.label] = channel
        this.channelWaitMap[channel.label] &&
        this.channelWaitMap[channel.label].resolve(channel)
    }

    _onInfoMessage(evt) {
        let bin = evt.data
        let msg = this.textDecoder.decode(bin)
        msg = JSON.parse(msg)
        console.log('info <-', msg)
        if (msg.type === 'call.ret') {
            let handler = this.callHostMap[msg.id]
            msg.result && (msg.result.__callId = msg.id)
            handler[msg.success ? 'resolve' : 'reject'](msg[msg.success ? 'result' : 'error'])
        } else if (msg.type === 'proxy.callback') {
            console.log(msg)
            let handler = this.proxyWaitMap[msg.id]
            typeof msg.result === 'string' && (msg.result = new String(msg.result))
            msg.result && (msg.result.__callId = msg.id)
            handler[msg.success ? 'resolve' : 'reject'](msg[msg.success ? 'result' : 'error'])
        } else {
            console.log('info msg', msg)
        }
    }

    async vncProxyChanel(slave, viewOnly) {
        let method = viewOnly ? 'view' : 'operate'
        let data = await this._callHost(method, {slave})
        let channelPromise
        if (this.channelMap[data.channel]) {
            channelPromise = Promise.resolve(this.channelMap[data.channel])
        } else {
            channelPromise = new Promise(resolve => {
                this.channelWaitMap[data.change] = {resolve}
            })
        }
        let channel = await channelPromise
        let callId = data.__callId
        let proxyPromise = new Promise((resolve, reject) => {
            this.proxyWaitMap[callId] = {
                resolve, reject
            }
        })
        await proxyPromise
        console.log(proxyPromise)
        return channel
    }

    _callHost(method, param) {
        const id = Math.random() + ''
        console.log('info -> ', {
            ...param,
            method,
            id,
        })
        this.infoChannel.send(JSON.stringify({
            ...param,
            method,
            id,
        }))
        const handler = {}
        this.callHostMap[id] = handler
        return new Promise((resolve, reject) => {
            handler.resolve = resolve
            handler.reject = reject
        })
    }
}


export async function connectHost(name) {
    if (Host.hostMap[name])
        return Host.hostMap[name]
    const id = Math.random() + ''
    const pc = newPeerConnection(id, function (candidate) {
        $rt.call('connect.host', {
            candidate: JSON.stringify(candidate),
            hostName: name,
            requestId: id
        })
    })
    connectionRequestMap[id] = {
        pc
    }
    let resolve
    let channelPromise = new Promise(res => {
        resolve = res
    })
    pc.ondatachannel = ev => {
        const channel = ev.channel
        if (channel.label === 'info') {
            // console.log(channel)
            resolve(channel)
        }
    }
    let infoChannel = await channelPromise
    return new Host(name, infoChannel, pc)
}


export async function getTicket() {
    if ($state.loading.getTicket || $state.ticket)
        throw new Error('cant repeat')
    $state.loading.getTicket = true
    let t = await $rt.call('ticket.new').finally(() => {
        $state.loading.getTicket = false
    })
    t.inService = false
    $state.ticket = t
    return t
}

export function showError(e) {
    $state.$message.error(e.message || e)
}
