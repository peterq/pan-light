import Base64 from './base64.js';

// PhantomJS can't create Event objects directly, so we need to use this
function make_event(name, props) {
    const evt = document.createEvent('Event');
    evt.initEvent(name, true, true);
    if (props) {
        for (let prop in props) {
            evt[prop] = props[prop];
        }
    }
    return evt;
}

let pc
(function () {
    pc = new RTCPeerConnection({
        iceServers: [
            {
                urls: 'turn:peterq.cn:1425',
                username: "pan_light_turn",
                credential: "pan_light_turn"
            }
        ]
    })
    console.log('peer connection ', pc)

    let sendChannel = pc.createDataChannel('init')

    pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState)
    pc.onicecandidate = event => {
        if (event.candidate === null) {
            console.log('candidate is ok')
            console.log(btoa(JSON.stringify(pc.localDescription)))
        }
    }

    pc.onnegotiationneeded = e => {
        console.log("onnegotiationneeded", e)
        pc.createOffer().then(d => pc.setLocalDescription(d)).catch(console.log.bind(console))
    }

    window.startSession = (sd) => {
        if (sd === '') {
            return alert('Session Description must not be empty')
        }

        try {
            pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(sd))))
        } catch (e) {
            alert(e)
        }
    }
})()

export default class RtcWebSocket {
    constructor(uri, protocols) {

        console.log('web rtc -> web socket', uri, protocols)
        let sendChannel = pc.createDataChannel('proxy')
        this.sendChannel = sendChannel
        console.log(sendChannel)
        sendChannel.onclose = () => {
            console.log('sendChannel has closed')
            this.close(0, 'data channel closed')
        }
        sendChannel.onopen = () => {
            console.log('sendChannel has opened')
            this._open()
        }
        sendChannel.onmessage = e => {
            // console.log(`Message from DataChannel '${sendChannel.label}' payload '${e.data}'`)
            this._receive_data(e.data)
        }

        this.url = uri;
        this.binaryType = "arraybuffer";
        this.extensions = "";

        if (!protocols || typeof protocols === 'string') {
            this.protocol = protocols;
        } else {
            this.protocol = protocols[0];
        }

        this._send_queue = new Uint8Array(20000);

        this.readyState = RtcWebSocket.CONNECTING;
        this.bufferedAmount = 0;

        this.__is_fake = true;
    }

    close(code, reason) {
        this.readyState = RtcWebSocket.CLOSED;
        if (this.onclose) {
            this.onclose(make_event("close", { 'code': code, 'reason': reason, 'wasClean': true }));
        }
    }

    send(data) {
        if (this.protocol == 'base64') {
            data = Base64.decode(data);
        } else {
            data = new Uint8Array(data);
        }
        // this._send_queue.set(data, this.bufferedAmount);
        // this.bufferedAmount += data.length;
        this.sendChannel.send(data)
        console.log('rtc web socket send', data)
    }

    // _get_sent_data() {
    //     const res = new Uint8Array(this._send_queue.buffer, 0, this.bufferedAmount);
    //     this.bufferedAmount = 0;
    //     return res;
    // }

    _open() {
        this.readyState = RtcWebSocket.OPEN;
        if (this.onopen) {
            this.onopen(make_event('open'));
            console.log('rtc web socket open')
        }
    }

    _receive_data(data) {
        // Break apart the data to expose bugs where we assume data is
        // neatly packaged
        this.onmessage(make_event("message", { 'data': data }));
        console.log('rtc web socket on message', data)
    }
}

RtcWebSocket.OPEN = WebSocket.OPEN;
RtcWebSocket.CONNECTING = WebSocket.CONNECTING;
RtcWebSocket.CLOSING = WebSocket.CLOSING;
RtcWebSocket.CLOSED = WebSocket.CLOSED;

RtcWebSocket.__is_fake = true;

RtcWebSocket.replace = () => {
    if (!WebSocket.__is_fake) {
        const real_version = WebSocket;
        // eslint-disable-next-line no-global-assign
        WebSocket = RtcWebSocket;
        RtcWebSocket.__real_version = real_version;
    }
};

RtcWebSocket.restore = () => {
    if (WebSocket.__is_fake) {
        // eslint-disable-next-line no-global-assign
        WebSocket = WebSocket.__real_version;
    }
};
