import {newPeerConnection} from "./realtime/webRtc"
import RealTime from "./realtime/realtime"
import Vue from "vue"
import State, {dataTemplate} from "./util/state"
import {registerProxyChannelResolver} from "./lib/vnc/core/RtcWebSocket"
import whatJpg from './assets/what.jpeg'
import {setWsFactory} from "./lib/vnc/core/websock"
import RtcWebSocket from "./lib/vnc/core/RtcWebSocket"
import ProxyWebSocket from "./lib/vnc/core/ProxyWebSocket"

// const $rt = new RealTime('ws://localhost:8001/demo/ws')

var connectWs
var p = new Promise(function (resolve) {
    connectWs = resolve
})
export const $rt = new RealTime((location.protocol === 'https:' ? 'wss' : 'ws') + '://' + location.host + '/demo/ws', p)
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

setWsFactory(function (uri, protocols) {
    uri = uri.replace('ws://', '').replace('wss://', '')
    let [hostName, slave] = uri.split('/')
    let host = $state.hosts.find(function (host) {
        return host.name === hostName
    })
    if (!host)
        throw new Error('host not found in host list: ' + hostName)
    if (host.wsAgentUrl) {
        // return new WebSocket('ws://172.17.0.2:5901', protocols)
        // return proxyWs(host.wsAgentUrl + '?slave=' + slave, protocols)
        return new ProxyWebSocket(host.wsAgentUrl + '?slave=' + slave, protocols)
    }
    return new RtcWebSocket(uri, protocols)
})


const connectionRequestMap = {}
console.log(process.env)
if (process.env.NODE_ENV === 'production') {
    setInterval(function () {
        window.debugObj = {}
    }, 3e3)
}

$rt.onRemote("host.candidate.ok", data => {
    const id = data.requestId
    const candidate = data.candidate
    const handler = connectionRequestMap[id]
    handler.pc.continueWithRemote(candidate)
})

function roomHandleUserBroadCast(room, data) {
    if (data.event === 'chat') {
        room.messages = room.messages.concat([{
            id: +new Date,
            type: 'chat',
            msg: data.payload,
            from: data.from
        }])
    }
}

function roomHandleUserTicketTurn(room, data) {

    console.log(data, $state.ticket)
    const {order} = data
    if ($state.ticket && $state.ticket.order === order) {
        $state.ticket.inService = true

        let {host, slave} = data
        $state.connectVnc = {
            host, slave, viewOnly: false,
            password: $state.ticket.ticket
        }

        $event.fire('operate.turn', data)
        console.log(data)
    }
    room.messages = room.messages.concat([{
        id: +new Date,
        type: 'system',
        evt: 'turn',
        sessionId: data.sessionId,
        ticket: data
    }])

}

let cdnPrefix = process.env.NODE_ENV === 'production' ? window.cdnPrefix : ''
$rt.on('room.new', room => {

    async function getSessionInfo(ids) {
        let newOnes = []
        ids.forEach(id => {
            if (!$state.userSessionInfo[id]) {
                newOnes.push(id)
                Vue.set($state.userSessionInfo, id, dataTemplate.deppClone('userSessionInfo'))
            }
        })
        let infoMap = await $rt.call('session.public.info', {sessionIds: newOnes})
        for (let id in infoMap) {
            infoMap[id].avatar = cdnPrefix + infoMap[id].avatar
            console.log(infoMap[id])
            $state.userSessionInfo[id] = infoMap[id]
        }
    }

    room.messages = []
    room.members = []
    room.sendMsg = function (msg) {
        room.broadcast('chat', msg)
        room.messages = room.messages.concat([{
            id: +new Date,
            type: 'chat',
            msg,
            from: $state.userSessionInfo.self.sessionId
        }])
    }

    $rt.call('room.members', {room: room.name})
        .then(members => room.members = members)
        .then(() => {
            getSessionInfo(room.members)
        })
    room.on('leave', () => {
        Vue.delete($state.roomMap, room.name)
    })
    room.onRemote('room.member.join', sessionId => {
        getSessionInfo([sessionId])
        room.messages = room.messages.concat([{
            id: +new Date,
            type: 'system',
            evt: 'join',
            sessionId
        }])
        room.members = room.members.concat([sessionId])
    })
    room.onRemote('room.member.remove', sessionId => {
        room.members = room.members.filter(id => id !== sessionId)
        room.messages = room.messages.concat([{
            id: +new Date,
            type: 'system',
            evt: 'leave',
            sessionId
        }])
    })
    room.onRemote('broadcast.user', data => roomHandleUserBroadCast(room, data))
    Vue.set($state.roomMap, room.name, room)
    // 全员群
    if (room.name === 'room.all.user') {
        room.onRemote('ticket.turn', data => roomHandleUserTicketTurn(room, data))
    }

    // slave 全员群
    if (room.name.indexOf('room.slave.all.user') === 0) {
        room.onRemote('slave.exit', data => {
            if (data.unexpected) {
                $state.$alert('demo 进程意外结束', 'Whoops', {type: 'error'})
            } else {
                $state.$message.info('体验结束')
            }
            $state.connectVnc = null
        })

        room.onRemote('broadcast.slave', data => {
            switch (data.event) {
                case 'operator.leave':
                    $state.$message.info('demo 操作用户已离开')
                    break
            }
        })
    }

})

$rt.onRemote('session.new', async session => {
    $state.resetData()
    let infoMap = await $rt.call('session.public.info', {sessionIds: [session.id]})
    infoMap[session.id].avatar = cdnPrefix + infoMap[session.id].avatar
    $state.userSessionInfo.self = {...infoMap[session.id], sessionId: session.id}
    $state.connected = true
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

export async function startApp(fn) {
    if (!await canStart())
        return
    connectWs(true)
    fn()
}


const roomWaitMap = {}

export async function getRoom(name) {
    let room = $rt.getRoom(name)
    if (room)
        return room
    return await new Promise(res => {
        roomWaitMap[name] = roomWaitMap[name] || []
        roomWaitMap.push(res)
    })
}

$rt.on('room.new', function (room) {
    console.log('new room', room)
    ;(roomWaitMap[room.name] || []).forEach(cb => cb(room))
})

function isOpenDev() {
    return new Promise(function (resolve) {
        var element = new Image()
        Object.defineProperty(element, 'id', {
            get: function () {
                resolve(true)
                return 0
            }
        })
        element.toString = function () {
            return 'hello pan-light'
        }
        console.log(element)
        console.clear()
        setTimeout(resolve, 100)
    })
}

function nestDebugger(depth) {
    if (depth === 1)
        return 'debugger'
    return `eval(${JSON.stringify(nestDebugger(depth - 1))});debugger`
}

function isOpenByDebugger() {
    let start = new Date()
    eval(nestDebugger(10))
    let end = new Date()
    return end - start > 100
}

// 禁止调试
async function canStart() {
    if (localStorage.getItem('debug') === 'pan-light')
        return true
    let open = await isOpenDev()
    if (!open) { // 没有打开开发工具, 等待未来打开
        var element = new Image()
        Object.defineProperty(element, 'id', {
            get: function () {
                setTimeout(() => {
                    fuckDebug()
                    location.reload()
                })
                return 0
            }
        })
        console.log(element)
    } else { // 打开了开发工具, 等待关闭
        fuckDebug()
        while (isOpenByDebugger()) {
            await new Promise(res => setTimeout(res, 500))
        }
        location.reload()
    }
    return !open
}

function fuckDebug() {
    console.clear()
    document.body.innerHTML = `<h1>偷窥人家可是不好的哦</h1>`
    consoleImage((cdnPrefix || location.origin) + '/demo' + whatJpg, 240, 240)
    console.log('%c想要演示系统源码? 快去点个star啦, 超过 200 star 开源此在线演示系统 https://github.com/peterq/pan-light', 'font-size:24px;color:#0a0')
    console.log('如果你想现在拿到源码, 你可以尝试分析一下: web端和服务端通信规则, 以及远程桌面的实现原理; 把分析结果发送到邮箱 me@peterq.cn , 我会回复源码哦. ps:难度其实不是特别大哦.')
    console.log('%c请求各位大佬不要对我的服务器进行压测, 阿里云最低配机器, 穷.', 'font-size:18px;')
}

function consoleImage(url, w, h) {
    console.log("%c+", `font-size: 1px; padding: ${~~(h / 2)}px ${~~(w / 2)}px; background: url(${url}) no-repeat; background-size: ${w}px ${h}px; color: transparent;`)
}
