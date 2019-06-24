import QtQuick 2.0
import QtMultimedia 5.0
import QtQuick.Controls 1.2
import QtQuick.Controls.Styles 1.2
import QtQuick.Layouts 1.1
import QtQuick.Dialogs 1.2
import QtQuick.Window 2.0
import "./UIComp"
import "../js/util.js" as Util
import "../js/global.js" as G
import "../js/app.js" as App

ApplicationWindow {
    id: player
    width: 1024
    height: 650
    minimumHeight: 325
    minimumWidth: 512
    visible: false
    title: 'pan-light video player'

    DataSaver {
        $key: 'video.player'
        property alias volume: mediaPlayer.volume
        property alias x: player.x
        property alias y: player.y
        property alias width: player.width
        property alias height: player.height
    }

    property bool playing: false

    property alias muted: mediaPlayer.muted
    property alias volume: mediaPlayer.volume
    property alias duration: mediaPlayer.duration
    property alias position: mediaPlayer.position
    property bool showControls: false
    property bool canHideControls: true
    property string videoTitle: 'pan-light video player'
    property int videoRotation: 0
    property real playbackRate: mediaPlayer.playbackRate
    property alias bufferProgress: mediaPlayer.bufferProgress
    property bool isLoading: [MediaPlayer.Loading, MediaPlayer.Stalled].indexOf(
        mediaPlayer.status) > -1
    signal controlsWillHide(var hideEvent)
    signal customerEvent(var event, var data)

    function playVideo(title, source) {
        visible = true
        requestActivate()
        videoTitle = title
        videoRotation = 0
        mediaPlayer.stop()
        if (source.loadingLinkText) {
            loadingLinkTip.text = source.loadingLinkText
            loadingLinkTip.visible = true
            source.then(function(link) {
                mediaPlayer.stop()
                mediaPlayer.source = link
                mediaPlayer.play()
            }).catch(function (err) {
                mediaPlayer.showError("解析链接错误:" + err)
            }).finally(function() {
                loadingLinkTip.visible = false
            })
        } else {
            mediaPlayer.stop()
            mediaPlayer.source = source
            mediaPlayer.play()
        }
    }

    function rotateVideo() {
        var pos = mediaPlayer.position
        videoRotation = (videoRotation + 90) % 360
        G.setTimeout(function () {
            seekAbs(pos)
        }, 500)
    }

    function setRate(rate) {
        mediaPlayer.playbackRate = rate
    }

    function play() {
        playing = true
        mediaPlayer.play()
    }

    function pause() {
        playing = false
        mediaPlayer.pause()
    }

    function tooglePlay() {
        return playing ? pause() : play()
    }

    function toogleMute() {
        muted = !muted
    }

    function seekAbs(pos) {
        mediaPlayer.seek(pos)
    }

    property int lastScreenVisibility
    property bool isFullScreen: visibility === Window.FullScreen
    function toggleFullScreen() {
        if (visibility != Window.FullScreen) {
            lastScreenVisibility = visibility
            visibility = Window.FullScreen
        } else {
            visibility = lastScreenVisibility
        }
    }

    property real hideControlsCheck: 0
    function hideControlsLater() {
        var rand = Math.random()
        hideControlsCheck = rand
        player.showControls = true
        G.setTimeout(function () {
            if (hideControlsCheck === rand) {
                var evt = {
                    "hide": true
                }
                player.controlsWillHide(evt)
                if (evt.hide)
                    player.showControls = false
            }
        }, 4e3)
    }

    onClosing: {
        player.destroy()
    }
    Component.onCompleted: {
        App.appState.player = player
    }
    Rectangle {
        id: background
        anchors.fill: parent
        color: 'black'
    }
    MediaPlayer {
        id: mediaPlayer
        property var errShowPromise: Util.Promise.resolve()
        source: ''
        onError: {
            console.log('-----------', error, errorString)
            // 同时2个弹窗有bug, 直接卡死, 改造成队列模式
            showError(errorString)
        }
        onPaused: {
            player.playing = false
        }
        onPlaying: {
            player.playing = true
        }
        function showError(errorString) {
            errShowPromise = errShowPromise.finally(function () {
                return Util.alert({
                                      "parent": player,
                                      "title": '播放器错误',
                                      "msg": errorString
                                  })
            })
        }
    }
    VideoOutput {
        id: vo
        x: (player.width - width) / 2
        y: (player.height - height) / 2
        width: (player.videoRotation / 90) % 2 ? parent.height : parent.width
        height: (player.videoRotation / 90) % 2 ? parent.width : parent.height
        source: mediaPlayer
        transform: Rotation {
            origin.x: vo.width / 2
            origin.y: vo.height / 2
            angle: player.videoRotation
        }
    }
    Tips {
    }
    MouseArea {
        id: controlControlsArea
        anchors.fill: parent
        hoverEnabled: true
        propagateComposedEvents: true
        enabled: true
        onPositionChanged: {
            player.hideControlsLater()
        }
        onDoubleClicked: {
            mouse.accepted = false
            player.toggleFullScreen()
        }
        cursorShape: player.showControls ? Qt.ArrowCursor : Qt.BlankCursor
    }
    ControlArea {
    }

    Text {
        id: loadingLinkTip
        anchors.centerIn: parent
        color: 'white'
        font.pointSize: 20
        visible: false
    }
}
