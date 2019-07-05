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

export default class ProxyWebSocket {
    constructor(uri, protocols) {

        this._ws = new WebSocket(uri, protocols)
        this._ws.onclose = (evt) => {
            console.log(evt)
            this.close(evt.code, evt.reason)
        }
        this._ws.onerror = (evt) => {
            this._error()
        }
        this._ws.onmessage = (evt) => {
            console.log(evt.data)
            this._ws.onmessage = this.onmessage
            this._open()
        }
        this._ws.onopen = (evt) => {

        }
        this._ws.binaryType = this.binaryType = "arraybuffer"

        if (!protocols || typeof protocols === 'string') {
            this.protocol = protocols
        } else {
             this.protocol = protocols[0]
        }
        this.readyState = ProxyWebSocket.CONNECTING
        this.__is_fake = true
    }


    close(code, reason) {
        this.readyState = ProxyWebSocket.CLOSED
        if (this.onclose) {
            this.onclose(make_event("close", {'code': code, 'reason': reason, 'wasClean': true}))
        }
        this._ws.close()
    }

    send(data) {
        // console.log('rtc web socket send', data)
        if (this.protocol === 'base64') {
            data = Base64.decode(data)
        } else {
            data = new Uint8Array(data)
        }
        this._ws.send(data)
    }


    _open() {
        this.readyState = ProxyWebSocket.OPEN
        if (this.onopen) {
            this.onopen(make_event('open'))
            console.log('rtc web socket open')
        }
    }

    _error() {
        this.readyState = ProxyWebSocket.OPEN
        if (this.onerror) {
            this.onerror(make_event('error'))
            console.log('rtc web socket error')
        }
    }

}


ProxyWebSocket.OPEN = WebSocket.OPEN
ProxyWebSocket.CONNECTING = WebSocket.CONNECTING
ProxyWebSocket.CLOSING = WebSocket.CLOSING
ProxyWebSocket.CLOSED = WebSocket.CLOSED

ProxyWebSocket.__is_fake = true

