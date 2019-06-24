import QtQuick 2.0

Rectangle {
    id: root
    property color successColor: 'green'
    property color failColor: 'red'
    property color warningColor: 'orange'
    property alias duration: hideTimer.interval
    property int padding: 20

    width: parent.width
    height: padding + msgText.implicitHeight
    y: -height

    Text {
        id: msgText
        text: ''
        color: 'white'
        width: parent.width - 20
        wrapMode: Text.WrapAnywhere
        anchors.centerIn: parent
    }

    MouseArea {
        id: ma
        anchors.fill: parent
        hoverEnabled: true
    }

    Timer {
        id: hideTimer
        interval: 2500
        running: false
        onTriggered: {
            if (ma.containsMouse) {
                restart()
                return
            }
            y = -height
        }
    }

    Behavior on y {
        PropertyAnimation {
            duration: 200
        }
    }

    function success(msg) {
        color = successColor
        show(msg)
    }

    function fail(msg) {
        color = failColor
        msg = msg || 'unspecifed error'
        if (typeof msg !== 'string') {
            if (msg.message) msg= msg.message
            else if (msg.toString()) msg = msg.toString()
            else msg = JSON.stringify(msg)
        }
        show(msg)
    }

    function warn(msg) {
        color = warningColor
        show(msg)
    }

    function show(msg) {
        y = 0
        msgText.text = msg
        hideTimer.restart()
    }
}
