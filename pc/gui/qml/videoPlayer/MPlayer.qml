import QtQuick 2.0
import QtMultimedia 5.0
import QtQuick.Controls 1.2
import QtQuick.Controls.Styles 1.2
import QtQuick.Layouts 1.1
import QtQuick.Dialogs 1.2
import QtQuick.Window 2.0
import './UIComp'
import "../js/util.js" as Util
import "../js/global.js" as G
import "../js/app.js" as App

ApplicationWindow {
    id: player
    width: 1024
    height: 650
    minimumHeight: 325
    minimumWidth: 512
    visible: true
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
    property string videoTitle: ''

    signal controlsWillHide(var hideEvent)
    signal customerEvent(var event, var data)

    function playVideo(title, source) {
        videoTitle = title
        mediaPlayer.stop()
        mediaPlayer.source = source
        mediaPlayer.play()
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
                var evt = { hide: true }
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
        autoPlay: true
        source: "file:///media/peterq/files/share/video/盗梦空间.Inception.2010.中英字幕.BDrip.AAC.720p.x264-人人影视.mp4"
        onError: {
            console.log('-----------', error, errorString)
        }
        onPaused: {
            player.playing = false
        }
        onPlaying: {
            player.playing = true
        }
    }
    VideoOutput {
        anchors.fill: parent
        source: mediaPlayer
    }
    Tips {}
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
    ControlArea {}
}
