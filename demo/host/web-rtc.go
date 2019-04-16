package host

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pion/webrtc/v2"
	"github.com/pkg/errors"
	"io"
	"log"
	"sync"
)

var webRtcApi *webrtc.API
var webRtcConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs:           []string{"turn:peterq.cn:1425"},
			Username:       "pan_light_turn",
			Credential:     "pan_light_turn",
			CredentialType: webrtc.ICECredentialTypePassword,
		},
	},
}

func init() {
	s := webrtc.SettingEngine{}
	s.DetachDataChannels()
	webRtcApi = webrtc.NewAPI(webrtc.WithSettingEngine(s))
}

type p2p struct {
	peerConnection    *webrtc.PeerConnection
	infoChannel       *webrtc.DataChannel
	ctx               context.Context
	cancel            context.CancelFunc
	infoChannelWriter io.Writer
	sendInfoLock      sync.Mutex
}

func handleNewUser(cand string, sessionId string) {
	remoteSd := webrtc.SessionDescription{}
	err := json.Unmarshal([]byte(cand), &remoteSd)
	if err != nil {
		log.Println("cand 解码错误", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	peerConnection, err := webRtcApi.NewPeerConnection(webRtcConfig)
	if err != nil {
		log.Println(err)
		return
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})
	infoChannel, err := peerConnection.CreateDataChannel("info", nil)
	if err != nil {
		log.Println("创建info channel错误", err)
	}
	p := &p2p{
		peerConnection: peerConnection,
		infoChannel:    infoChannel,
		ctx:            ctx,
		cancel:         cancel,
	}
	infoChannel.OnOpen(p.handleInfoChannel)

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		log.Println("获取本地sd错误", err)
		return
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		log.Println("设置本地sd错误", err)
		return
	}

	err = peerConnection.SetRemoteDescription(remoteSd)
	if err != nil {
		log.Println("设置远端sd错误", err)
		return
	}
	func() {
		host.p2pMapLock.Lock()
		defer host.p2pMapLock.Unlock()
		host.p2pMap[sessionId] = p
	}()
	<-ctx.Done()
	func() {
		host.p2pMapLock.Lock()
		defer host.p2pMapLock.Unlock()
		delete(host.p2pMap, sessionId)
	}()
	err = peerConnection.Close()
	if err != nil {
		log.Println("关闭p2p错误", err)
	}
}

func (p *p2p) handleInfoChannel() {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	defer p.cancel()
	readWriter, err := p.infoChannel.Detach()
	if err != nil {
		err = errors.Wrap(err, "detach info channel 错误")
		return
	}
	p.infoChannelWriter = readWriter
	buf := make([]byte, 2048)
	for {
		n, err := readWriter.Read(buf)
		if err != nil {
			err = errors.Wrap(err, "read info channel 错误")
			return
		}
		var data gson
		err = json.Unmarshal(buf[:n], &data)
		if err != nil {
			err = errors.Wrap(err, "json decode错误")
			return
		}
		go p.handleInfoMsg(data)
	}
}

func (p *p2p) sendInfo(data gson) error {
	p.sendInfoLock.Lock()
	defer p.sendInfoLock.Unlock()
	bin, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "json encode错误")
	}
	n, err := p.infoChannelWriter.Write(bin)
	if err != nil {
		return errors.Wrap(err, "写入数据错误")
	}
	if n != len(bin) {
		return errors.New("写入数据不完整")
	}
	return nil
}

func (p *p2p) handleInfoMsg(data gson) {
	method := data["method"].(string)
	id := data["id"].(string)
	var ret gson
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if err != nil {
			p.sendInfo(gson{
				"type":    "call.ret",
				"id":      id,
				"error":   err.Error(),
				"success": false,
			})
		} else {
			p.sendInfo(gson{
				"type":    "call.ret",
				"id":      id,
				"result":  ret,
				"success": true,
			})
		}
	}()
	if method == "view" {
		slaveName := data["slave"].(string)
		holder, ok := host.holderMap[slaveName]
		if !ok {
			err = errors.New("slave 不存在")
			return
		}
		c := "proxy.view" + slaveName
		viewChanel, err := p.peerConnection.CreateDataChannel(c, nil)
		if err != nil {
			err = errors.Wrap(err, "创建view channel 错误")
			return
		}
		rw, err := viewChanel.Detach()
		if err != nil {
			err = errors.Wrap(err, "view channel detach 错误")
			return
		}
		go holder.VncProxy(rw, func(err error) {
			if err != nil {
				p.sendInfo(gson{
					"type":    "proxy.callback",
					"success": false,
					"error":   err.Error(),
					"id":      id,
				})
			} else {
				p.sendInfo(gson{
					"type":    "proxy.callback",
					"success": true,
					"id":      id,
				})
			}
		}, true)
		ret = gson{
			"channel": c,
		}
		return
	}
}
