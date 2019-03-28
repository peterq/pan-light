import QtQuick 2.0
import QtGraphicalEffects 1.0
import "./UIComp"
import "../js/app.js" as App

Item {
    id: controls
    anchors.bottom: parent.bottom
    width: parent.width
    height: 60
    property var player: App.appState.player
    opacity: player.showControls ? 1 : 0

    function mouseInControls() {
        function walkItem(item) {
            if (typeof item.hovered !== 'undefined'
                    || typeof item.containsMouse !== 'undefined') {
                var mouseIn = item.hovered || item.containsMouse
                // console.log('mouse in', item.hovered, item.containsMouse)
                if (mouseIn)
                    return true
            }
            if (typeof item.children === 'undefined')
                return false
            for (var i = 0; i < item.children.length; i++) {
                var b = walkItem(item.children[i])
                if (b)
                    return true
            }
            return false
        }
        return walkItem(controls)
    }

    // 阻止隐藏控制区
    Connections {
        target: player
        onControlsWillHide: {
            console.log('will hide')
            if (controls.mouseInControls())
                hideEvent.hide = false
        }
    }

    Component.onCompleted: {
        function walkItem(item) {
            if (item.onHoveredChanged) {
                item.onHoveredChanged.connect(function () {
                    if (!controls.mouseInControls()) {
                        player.hideControlsLater()
                    }
                })
            }
            if (typeof item.children === 'undefined')
                return false
            for (var i = 0; i < item.children.length; i++) {
                walkItem(item.children[i])
            }
        }
        walkItem(controls)
    }

    // 检测鼠标是否在控制区
    MouseArea {
        id: handleMouseInControlArea
        anchors.fill: controls
        hoverEnabled: true
        propagateComposedEvents: true
        acceptedButtons: Qt.NoButton
        enabled: true
    }

    LinearGradient {
        width: parent.width
        height: parent.height * 1.2
        gradient: Gradient {
            GradientStop {
                position: 0.0
                color: Qt.rgba(0, 0, 0, 0)
            }
            GradientStop {
                position: 0.4
                color: Qt.rgba(0, 0, 0, 0.4)
            }
            GradientStop {
                position: 1.0
                color: Qt.rgba(0, 0, 0, 0.4)
            }
        }
        start: Qt.point(0, 0)
        end: Qt.point(0, height)
    }
    // 按钮区
    Item {
        id: btns
        width: parent.width - 30
        height: parent.height * 0.7
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.bottom: parent.bottom
        ButtonPlay {
            id: playButton
            icon.height: parent.height / 1.25
            icon.width: parent.height / 1.25
            anchors.verticalCenter: parent.verticalCenter
        }

        VolumeButton {
            id: volumeButton
            anchors.left: playButton.right
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            icon.height: parent.height / 1.25
            icon.width: parent.height / 1.25
        }
        VolumeSlider {
            id: volumeSlider
            anchors.left: volumeButton.right
            anchors.top: parent.top
            anchors.bottom: parent.bottom
            height: parent.height
            visible: mouseAreaVolumeArea.containsMouse || volumeButton.hovered
            width: visible ? implicitWidth : 0
        }
        MouseArea {
            id: mouseAreaVolumeArea
            anchors.bottom: parent.bottom
            anchors.left: volumeButton.left
            anchors.right: volumeSlider.right
            anchors.top: parent.top
            height: parent.height
            width: volumeButton.width + (volumeSlider.visible ? volumeSlider.width : 0)
            hoverEnabled: true
            propagateComposedEvents: true
            acceptedButtons: Qt.NoButton
            onWheel: {
                var delta = wheel.angleDelta.y / 120 //一刻滚轮代表正负120度，所以除以120等于1或者-1
                var v = player.volume + delta * 0.1
                if (v > 1)
                    v = 1
                else if (v < 0)
                    v = 0
                player.volume = v
            }
        }
        TimeText {
            id: timeText
            height: parent.height / 1.25
            width: parent.height / 1.25
            anchors.verticalCenter: parent.verticalCenter
            anchors.left: volumeSlider.right
        }
        FullScreenButton {
            id: fullScreenButton
            icon.height: parent.height / 1.25
            icon.width: parent.height / 1.25
            anchors.verticalCenter: parent.verticalCenter
            anchors.right: parent.right
        }

        OpenFileButton {
            id: openFileButton
            icon.height: parent.height / 1.8
            icon.width: parent.height / 1.8
            anchors.verticalCenter: parent.verticalCenter
            anchors.right: fullScreenButton.left
        }
        RotateButton {
            id: rotateButton
            icon.height: parent.height / 1.8
            icon.width: parent.height / 1.8
            anchors.verticalCenter: parent.verticalCenter
            anchors.right: openFileButton.left
        }
        PlayRateButton {
            id: playRateButton
            anchors.verticalCenter: parent.verticalCenter
            anchors.right: rotateButton.left
        }
    }

    TimeSlider {
        id: timeSlider
        width: parent.width - 30
        height: parent.height * 0.3
        anchors.bottom: btns.top
        anchors.horizontalCenter: parent.horizontalCenter
    }

    Behavior on opacity {
        PropertyAnimation {
            duration: 1000
        }
    }
    states: State {
        name: "hide"
        when: opacity === 0
        PropertyChanges {
            target: controls
            visible: false
        }
    }
}
