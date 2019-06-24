import QtQuick 2.0
import QtQuick.Window 2.3
import '../js/global.js' as G
Item {
    property var mainWindow
    MouseArea {
        property point clickPos: "0,0"
        anchors.fill: parent
//        drag.minimumX: 0
//        drag.maximumX: Screen.desktopAvailableWidth - mainWindow.width
//        drag.minimumY: 0
//        drag.maximumY: Screen.desktopAvailableHeight - mainWindow.heigh
        onPressed: {
            mainWindow.requestActivate()
            clickPos = Qt.point(mouseX, mouseY)
        }
        onPositionChanged: {
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            mainWindow.x += delta.x
            mainWindow.y += delta.y
        }
    }
    Component.onCompleted: {
        mainWindow = G.root
    }
}
