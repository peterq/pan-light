import Base64 from './base64.js'

// PhantomJS can't create Event objects directly, so we need to use this
function make_event(name, props) {
    const evt = document.createEvent('Event')
    evt.initEvent(name, true, true)
    if (props) {
        for (let prop in props) {
            evt[prop] = props[prop]
        }
    }
    return evt
}

let getProxyChannel = function () {
    return Promise.reject('please register the proxy channel resolver')
}

export function registerProxyChannelResolver(fn) {
    getProxyChannel = fn
}

export default class RtcWebSocket {
    constructor(uri, protocols) {

        // console.log('web rtc -> web socket', uri, protocols)

        getProxyChannel(uri)
            .then((proxyChannel) => {
                this.proxyChannel = proxyChannel
                console.log('proxyChannel', proxyChannel)
                proxyChannel.onclose = () => {

                    console.log('sendChannel has closed')
                    this.close(0, 'data channel closed')
                }
                // proxyChannel.onopen = () => {
                //     console.log('sendChannel has opened')
                //     this._open()
                // }
                proxyChannel.onmessage = e => {
                    // console.log(`Message from DataChannel '${proxyChannel.label}' payload '${e.data}'`)
                    this._receive_data(e.data)
                }
                this._open()
            })
            .catch(error => {
                console.log('get proxy channel', error)
                this._error()
                this.close(-1, new Error(error).message)
            })

        this.url = uri
        this.binaryType = "arraybuffer"
        this.extensions = ""

        if (!protocols || typeof protocols === 'string') {
            this.protocol = protocols
        } else {
            this.protocol = protocols[0]
        }

        this._send_queue = new Uint8Array(20000)

        this.readyState = RtcWebSocket.CONNECTING
        this.bufferedAmount = 0

        this.__is_fake = true
    }


    close(code, reason) {
        this.readyState = RtcWebSocket.CLOSED
        if (this.onclose) {
            this.onclose(make_event("close", {'code': code, 'reason': reason, 'wasClean': true}))
        }
        this.proxyChannel.close()
    }

    send(data) {
        // console.log('rtc web socket send', data)
        if (this.protocol == 'base64') {
            data = Base64.decode(data)
        } else {
            data = new Uint8Array(data)
        }
        this.proxyChannel.send(data)
    }

    _open() {
        this.readyState = RtcWebSocket.OPEN
        if (this.onopen) {
            this.onopen(make_event('open'))
            console.log('rtc web socket open')
        }
    }

    _error() {
        this.readyState = RtcWebSocket.OPEN
        if (this.onerror) {
            this.onerror(make_event('error'))
            console.log('rtc web socket error')
        }
    }

    _receive_data(data) {
        // Break apart the data to expose bugs where we assume data is
        // neatly packaged
        this.onmessage(make_event("message", {'data': data}))
        // console.log('rtc web socket on message', data)
    }
}

RtcWebSocket.OPEN = WebSocket.OPEN
RtcWebSocket.CONNECTING = WebSocket.CONNECTING
RtcWebSocket.CLOSING = WebSocket.CLOSING
RtcWebSocket.CLOSED = WebSocket.CLOSED

RtcWebSocket.__is_fake = true

RtcWebSocket.replace = () => {
    if (!WebSocket.__is_fake) {
        const real_version = WebSocket
        // eslint-disable-next-line no-global-assign
        WebSocket = RtcWebSocket
        RtcWebSocket.__real_version = real_version
    }
}

RtcWebSocket.restore = () => {
    if (WebSocket.__is_fake) {
        // eslint-disable-next-line no-global-assign
        WebSocket = WebSocket.__real_version
    }
}
