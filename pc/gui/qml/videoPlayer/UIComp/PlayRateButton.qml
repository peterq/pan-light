import QtQuick 2.0
import QtQml.Models 2.2
import QtQuick.Window 2.3
import QtQml 2.11
import "../../js/global.js" as G
import "../../js/app.js" as App

Item {
    id: root
    property var player: App.appState.player
    property var rates: [0.5, 0.75, 1, 1.25, 1.5, 1.75, 2]
    height: parent.height
    width: height * 1.5
    Text {
        text: 'x ' + player.playbackRate
        anchors.fill: parent
        horizontalAlignment: Text.AlignHCenter
        verticalAlignment: Text.AlignVCenter
        color: 'white'
        font.pointSize: 12
    }
    MouseArea {
        id: labelMa
        anchors.fill: parent
        hoverEnabled: true
        propagateComposedEvents: true
        //acceptedButtons: Qt.NoButton
        onEntered: {
            G.setTimeout(function () {
                if (labelMa.containsMouse)
                    ratePopup.visible = true
            }, 300)
        }
        onExited: {
            G.setTimeout(function () {
                if (!labelMa.containsMouse && !optionsMa.containsMouse)
                    ratePopup.visible = false
            }, 300)
        }
    }

    Window {
        id: ratePopup
        width: 200
        height: itemHeight * root.rates.length
        flags: Qt.FramelessWindowHint | Qt.Window | Qt.WindowStaysOnTopHint
               | Qt.X11BypassWindowManagerHint | Qt.Tool
        visible: false
        color: 'transparent'
        property int cursorIndex: -1
        property int itemHeight: 50

        Rectangle {
            anchors.fill: parent
            color: Qt.rgba(0, 0, 0, .6)
            border.color: 'white'
//            radius: 10
            clip: true
            Repeater {
                model: root.rates
                Rectangle {
                    id: rateItem
                    width: ratePopup.width
                    height: ratePopup.itemHeight
                    x: 0
                    y: index * height
                    color: ratePopup.cursorIndex === index ? Qt.rgba(
                                                                1, 1, 1,
                                                                0.6) : Qt.rgba(
                                                                0, 0, 0, 0)
                    property var rate: modelData
                    Text {
                        anchors.fill: parent
                        text: 'x ' + rateItem.rate
                        horizontalAlignment: Text.AlignHCenter
                        verticalAlignment: Text.AlignVCenter
                        color: player.playbackRate === rateItem.rate ? 'orange' : 'white'
                        font.pointSize: 12
                    }
                }
            }
        }

        onVisibleChanged: {
            var p = root.mapToGlobal(root.width / 2, 0)
            x = p.x - width / 2
            y = p.y - height
        }

        Connections {
            target: player
            onControlsWillHide: {
                if (optionsMa.containsMouse)
                    hideEvent.hide = false
            }
        }

        MouseArea {
            id: optionsMa
            hoverEnabled: true
            anchors.fill: parent
            onPositionChanged: {
                ratePopup.cursorIndex = Math.floor(
                            mouse.y / ratePopup.itemHeight)
            }
            onExited: {
                G.setTimeout(function () {
                    if (!labelMa.containsMouse && !optionsMa.containsMouse)
                        ratePopup.visible = false
                }, 300)
                ratePopup.cursorIndex = -1
            }
            onClicked: {
                if (ratePopup.cursorIndex > -1) {
                    ratePopup.visible = false
                    player.setRate(root.rates[ratePopup.cursorIndex])
                }
            }
        }
    }
}
