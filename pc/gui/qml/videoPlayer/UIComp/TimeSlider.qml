import QtQuick 2.0
import QtQuick.Controls 2.3
import "../../js/app.js" as App
import "../../js/global.js" as G

Slider {
    id: timeBar
    to: player.duration
    value: waitChangeTo ? changeTo : player.position
    property int changeTo: 0
    property bool waitChangeTo: false
    property var player: App.appState.player
    property bool show: false
    implicitWidth: Math.max(
                       background ? background.implicitWidth : 0,
                       (handle ? handle.implicitWidth : 0) + leftPadding + rightPadding)
    implicitHeight: Math.max(
                        background ? background.implicitHeight : 0,
                        (handle ? handle.implicitHeight : 0) + topPadding + bottomPadding)
    onMoved: {
        waitChangeTo = true
        changeTo = value
        var saved = value
        G.setTimeout(function () {
            if (saved === changeTo) {
                player.seekAbs(saved)
                waitChangeTo = false
            }
        }, 500)
    }
    background: Rectangle {
        x: timeBar.leftPadding
        y: timeBar.topPadding + timeBar.availableHeight / 2 - height / 2
        implicitWidth: 150
        implicitHeight: 4
        width: timeBar.availableWidth
        height: implicitHeight
        radius: 2
        color: "#bdbebf"

        Rectangle {
            width: timeBar.visualPosition * parent.width
            height: parent.height
            color: "red"
            radius: 2
        }
    }

    handle: Rectangle {
        x: timeBar.leftPadding + timeBar.visualPosition * (timeBar.availableWidth - width)
        y: timeBar.topPadding + timeBar.availableHeight / 2 - height / 2
        implicitWidth: 16
        implicitHeight: 16
        radius: 8
        color: timeBar.pressed ? "#f0f0f0" : "#f6f6f6"
        border.color: "#bdbebf"
    }
    MouseArea {
        id: hoverMouseArea
        anchors.fill: parent
        hoverEnabled: true
        propagateComposedEvents: true
        acceptedButtons: Qt.NoButton
        enabled: true
        ToolTip {
            id: tip
            visible: parent.containsMouse
            x: timeBar.pressed ?
                   timeBar.handle.x  - 0.5 * width + 0.5 * timeBar.handle.width
                 : parent.mouseX - 0.5 * width
            text: timeFormat((timeBar.pressed ?
                                  timeBar.handle.x + 0.5 * timeBar.handle.width
                                : parent.mouseX) / parent.width * player.duration)
            contentItem: Text {
                text: tip.text
                font: tip.font
                color: 'white'
            }
            background: Rectangle {
                color: Qt.rgba(0, 0, 0, .5)
                radius: Math.min(tip.width, tip.height) * 0.2
            }
            function timeFormat(time) {
                var sec = Math.floor(time / 1000)
                var hours = Math.floor(sec / 3600)
                var minutes = Math.floor((sec - hours * 3600) / 60)
                var seconds = sec - hours * 3600 - minutes * 60
                var hh, mm, ss
                if (hours.toString().length < 2)
                    hh = "0" + hours.toString()
                else
                    hh = hours.toString()
                if (minutes.toString().length < 2)
                    mm = "0" + minutes.toString()
                else
                    mm = minutes.toString()
                if (seconds.toString().length < 2)
                    ss = "0" + seconds.toString()
                else
                    ss = seconds.toString()
                return hh + ":" + mm + ":" + ss
            }
        }
    }


}
