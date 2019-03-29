import QtQuick 2.0
import "../../js/app.js" as App

Rectangle {
    property var player: App.appState.player
    width: t1.implicitWidth + 10
    height: width
    radius: width / 15
    color: Qt.rgba(0, 0, 0, .6)
    opacity: player.isLoading ? 1 : 0

    Behavior on opacity {
        PropertyAnimation {
            duration: 200
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

    Item {
        anchors.centerIn: parent
        height: t1.implicitHeight + t2.implicitHeight + 5
        width: parent.width
        Text {
            id: t1
            width: implicitWidth
            height: implicitHeight
            anchors.top: parent.top
            anchors.horizontalCenter: parent.horizontalCenter
            text: 'pan-light奋力加载中...'
            color: 'white'
        }
        Text {
            id: t2
            width: implicitWidth
            height: implicitHeight
            anchors.bottom: parent.bottom
            anchors.horizontalCenter: parent.horizontalCenter
            text: Math.round(100 * player.bufferProgress) + '%'
            color: 'white'
        }
    }
}
