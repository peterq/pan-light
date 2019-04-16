import gzip from '../lib/gzip'

export default class Rpc {

    openPromise
    sessionId
    sessionSecret

    constructor (url) {
        this.encrypt = true
        this.encryptKey = 'pan-light'
        this.url = url
        this.eventListener = {}
        this.requestMap = {}
        this.init()
    }

    init () {

        this.openPromise = this.connect()

        this.onRemote('session.new', data => {
            this.sessionId = data.id
            this.sessionSecret = data.secret
        })

        setInterval(() => this.call('ping'), 10e3)
    }

    connect() {
        return new Promise((resolve) => {
            const ws  = this.ws = new WebSocket(this.url)
            if (this.encrypt) {
                ws.binaryType = 'arraybuffer'
            }
            ws.onopen = async evt => (resolve(), this.onWsOpen(evt))
            ws.onerror = evt => this.onWsError(evt)
            ws.onmessage = evt => this.onWsMessage(evt)
            ws.onclose = evt => this.onWsClose(evt)
        })
    }

    reconnnect () {
        this.ws.close()
        return this.connect()
    }


    onWsError (evt) {
        console.log('ws error', evt)
        let cbs = this.eventListener['rt.error'] || []
        cbs.forEach(function (cb) {
            setTimeout(function () {
                cb('realtime.error')
            })
        })
    }

    onWsClose (evt) {
        console.log('ws close', evt)
        let cbs = this.eventListener['realtime.closed'] || []
        cbs.forEach(function (cb) {
            setTimeout(function () {
                cb('realtime.closed')
            })
        })
        setTimeout(() => {
            this.reconnnect()
        }, 5e3)
    }

    onWsOpen () {
        if (this.sessionId) {
            this.wsSend({
                type: 'session.resume',
                sessionId: this.sessionId,
                sessionSecret: this.sessionSecret,
            })
        } else {
            this.wsSend({
                type: 'session.new'
            })
        }
        this.wsSend({
            role: 'user'
        })
    }

    wsSend(data) {
        if (!(data.type === 'call' && data.method === 'user.ping'))
            console.log('ws ->', data)
        let str = JSON.stringify(data)
        let send = this._encrypt(str)
        this.ws.send(send)
    }

    onWsMessage (evt) {
        const data = JSON.parse(this._decrypt(evt.data))
        if (!(
            (data.type === 'event' && data.event === 'ping') ||
            (data.type === 'call.result' && data.data === 'pong')
        )) {
            console.log('ws <-', data)
        }
        if (data.type === 'event') {
            let cbs = this.eventListener['$remote.' + data.event] || []
            cbs.forEach(function (cb) {
                setTimeout(function () {
                    cb(data.payload)
                })
            })
            return
        }
        if (data.type === 'call.result') {
            const { id, success } = data
            if (success) {
                this.requestMap[id].resolve(data.result)
            } else {
                this.requestMap[id].reject(data.error)
            }
            delete this.requestMap[id]
            return
        }
    }

    addEventListener (name, cb) {
        this.eventListener[name] = this.eventListener[name] || []
        this.eventListener[name].push(cb)
    }

    on (name, cb) {
        return this.addEventListener(name, cb)
    }

    onRemote(name, cb) {
        return this.on('$remote.' + name, cb)
    }

    call (method, param = {}) {
        method = 'user.' + method
        const id = 'cb' + new Date().getTime() + ~~(Math.random() * 1e5)
        let done = {}
        let promise = new Promise(function (resolve, reject) {
            done.resolve = resolve
            done.reject = reject
        })
        this.requestMap[id] = done
        this.wsSend({
            type: 'call',
            method,
            param,
            id
        })
        return promise
    }


    _encrypt(str) {
        if (!this.encrypt) {
            return str
        }
        let bin = gzip.zip(str)
        // console.log('发送, 压缩', bin)
        let cipher = new ArrayBuffer(bin.length)
        cipher = new Uint8Array(cipher)
        for (let i = 0; i < cipher.length; i++) {
            cipher[i] = bin[i] ^ this.encryptKey[i % this.encryptKey.length].charCodeAt(0)
        }
        // console.log('发送, 加密', cipher)
        return cipher.buffer
    }

    _decrypt(buf) {
        if (!this.encrypt) {
            return buf
        }
        let plain = Array.prototype.slice.call(new Uint8Array(buf), 0)
        plain.forEach((b, i) => {
            plain[i] ^=  this.encryptKey[i % this.encryptKey.length].charCodeAt(0)
        })
        // console.log('接收, 明文', JSON.stringify(plain))
        let raw = gzip.unzip(plain)
        // console.log('接收, 解压', raw)
        return raw.map(b => String.fromCharCode(b)).join('')
    }

}
