import QtQuick 2.0
import './UIComp'
import "../js/app.js" as App
Item {
    property var player: App.appState.player

    anchors.fill: parent

    PlayIcon {
        anchors.centerIn: parent
        width: 40
        height: width
    }

    ForwardBackward {
        anchors.centerIn: parent
        width: 40
        height: width
    }

    VolumeIcon {
        anchors.centerIn: parent
        width: 40
        height: width
    }

    ActionTips {
        anchors.topMargin: 10
        anchors.rightMargin: 10
        anchors.right: parent.right
        anchors.top: parent.top
    }
    Item {
        id: topBar
        width: parent.width
        height: 50
        anchors.top: parent.top
        opacity: player.showControls ? 1 : 0
        VideoTitle {}
        Behavior on opacity {
            PropertyAnimation {
                duration: 1000
            }
        }
        states: State {
            name: "hide"
            when: opacity === 0
            PropertyChanges {
                target: topBar
                visible: false
            }
        }
    }

    LoadingTips {
        anchors.centerIn: parent
    }
}
