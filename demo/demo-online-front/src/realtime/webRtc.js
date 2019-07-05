export function newPeerConnection(requestId, onLocalCandidate) {
    let pc = new RTCPeerConnection({
        iceServers: [
            {
                urls: 'turn:peterq.cn:1425',
                username: "pan_light_turn",
                credential: "pan_light_turn"
            }
        ]
    })
    window.debugObj.p2p = window.debugObj.p2p || {}
    window.debugObj.p2p[requestId] = pc
    pc.createDataChannel('init')

    // pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState, e)

    pc.onicecandidate = event => {
        if (event.candidate === null) {
            onLocalCandidate(pc.localDescription)
        }
    }

    pc.onnegotiationneeded = e => {
        // console.log("onnegotiationneeded", e)
        window.debugObj.onnegotiationneeded = e
        pc.createOffer().then(d => pc.setLocalDescription(d)).catch(console.log.bind(console))
    }
    pc.continueWithRemote = function (sd) {
        return pc.setRemoteDescription(new RTCSessionDescription(sd))
    }
    return pc
}
