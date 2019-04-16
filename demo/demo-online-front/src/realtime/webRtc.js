let pc = new RTCPeerConnection({
    iceServers: [
        {
            urls: 'turn:peterq.cn:1425',
            username: "pan_light_turn",
            credential: "pan_light_turn"
        }
    ]
})
window.debugObj.pc = pc
pc.createDataChannel('init')
pc.ondatachannel = e => console.log(e)

pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState, e)

pc.onicecandidate = event => {
    if (event.candidate === null) {
        // console.log('candidate is ok', pc.localDescription)
        window.$event.fire('rtc.candidate', pc.localDescription)
    }
}

pc.onnegotiationneeded = e => {
    // console.log("onnegotiationneeded", e)
    window.debugObj.onnegotiationneeded = e
    pc.createOffer().then(d => pc.setLocalDescription(d)).catch(console.log.bind(console))
}

export default function setRemoteDescription(sd)
{
    pc.setRemoteDescription(new RTCSessionDescription(sd))
}
