import QtQuick 2.0
import QtQuick.Controls 2.3
import "../../js/app.js" as App

Slider {
    id: volumeBar
    to: 100
    value: Math.round(player.volume * 100)
    property var player: App.appState.player
    property bool show: false
    implicitWidth: Math.max(
                       background ? background.implicitWidth : 0,
                       (handle ? handle.implicitWidth : 0) + leftPadding + rightPadding)
    implicitHeight: Math.max(
                        background ? background.implicitHeight : 0,
                        (handle ? handle.implicitHeight : 0) + topPadding + bottomPadding)
    onMoved: {
        player.volume = volumeBar.value / 100
    }
    background: Rectangle {
        x: volumeBar.leftPadding
        y: volumeBar.topPadding + volumeBar.availableHeight / 2 - height / 2
        implicitWidth: 150
        implicitHeight: 4
        width: volumeBar.availableWidth
        height: implicitHeight
        radius: 2
        color: "#bdbebf"

        Rectangle {
            width: volumeBar.visualPosition * parent.width
            height: parent.height
            color: "red"
            radius: 2
        }
    }

    handle: Rectangle {
        x: volumeBar.leftPadding + volumeBar.visualPosition * (volumeBar.availableWidth - width)
        y: volumeBar.topPadding + volumeBar.availableHeight / 2 - height / 2
        implicitWidth: 16
        implicitHeight: 16
        radius: 8
        color: volumeBar.pressed ? "#f0f0f0" : "#f6f6f6"
        border.color: "#bdbebf"
    }

    Behavior on width {
        PropertyAnimation {
            duration: 800
            easing.type: Easing.OutCubic
        }
    }
}
