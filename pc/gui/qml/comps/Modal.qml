import QtQuick 2.0
import QtQuick.Window 2.3
import '../js/util.js' as Util

Window {
    id: ratePopup
    width: 200
    height: 200
    flags: Qt.FramelessWindowHint | Qt.WindowStaysOnTopHint
           | Qt.X11BypassWindowManagerHint | Qt.Tool | Qt.WindowMinimizeButtonHint
    visible: true
    color: 'white'
    modality: Qt.WindowModal
//    modality: Qt.ApplicationModal
    property int cursorIndex: -1
    property int itemHeight: 50

    Component.onCompleted: {

        var p = Util.bridge.cursorPos()
        x = p.x
        y = p.y
    }
}
