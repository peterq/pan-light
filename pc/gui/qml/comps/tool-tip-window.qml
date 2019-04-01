import QtQuick 2.0
import QtQuick.Window 2.3
import QtQuick.Controls 2.2
import '../js/util.js' as Util

Window {
    id: popup
    width: contentArea.width
    height: contentArea.height
    flags: Qt.FramelessWindowHint | Qt.WindowStaysOnTopHint
           | Qt.X11BypassWindowManagerHint | Qt.Tool
    visible: false
    color: 'red'
    property real showId

    Rectangle {
        id: contentArea
        color: 'black'
        width: label.width + 10
        height: label.height + 10
        Label {
            id: label
            width: Math.min(label.implicitWidth, 300)
            height: label.implicitHeight
            anchors.centerIn: parent
            wrapMode: Text.Wrap
            color: 'white'
        }
    }

    function show(str) {
        label.text = str
        var p = Util.bridge.cursorPos()
        var x = p.x + 15
        var y = p.y + 15
        if (x + popup.width > Screen.desktopAvailableWidth)
            x = Screen.desktopAvailableWidth - popup.width
        if (y + popup.height > Screen.desktopAvailableHeight)
            y = p.y - 10 - popup.height
        popup.x = x
        popup.y = y
        popup.visible = true
        showId = Math.random()
        return showId
    }

    function hide(id) {
       if (id === showId) {
           visible = false
       }
    }
}
