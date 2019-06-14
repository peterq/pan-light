import QtQuick 2.0
import QtQuick.Window 2.3
import "../js/global.js" as G

Item {
    property int maWidth: 5
    property var mainWindow
    // 右
    MouseArea {
        property point clickPos: "0,0"
        anchors.right: parent.right
        y: maWidth
        width: maWidth
        height: parent.height - 2 * maWidth
        cursorShape: Qt.SizeHorCursor
        onPressed: {
            clickPos = Qt.point(mouseX, mouseY)
        }
        onPositionChanged: {
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            mainWindow.width += delta.x
        }
    }
    // 左
    MouseArea {
        property point clickPos: "0,0"
        property real t: 0
        y: maWidth
        width: maWidth
        height: parent.height - 2 * maWidth
        cursorShape: Qt.SizeHorCursor
        onPressed: {
            clickPos = Qt.point(mouseX, mouseY)
        }
        onPositionChanged: {
            var now = +new Date
            if (now - t < 50)
                return
            t = now
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            if (!(mainWindow.minimumWidth === mainWindow.width
                  && delta.x > 0)) {
                mainWindow.width -= delta.x
                mainWindow.x += delta.x
            }
        }
    }
    // 上
    MouseArea {
        property point clickPos: "0,0"
        property real t: 0
        x: maWidth
        width: parent.width - 2 * maWidth
        height: maWidth
        cursorShape: Qt.SizeVerCursor
        onPressed: {
            clickPos = Qt.point(mouseX, mouseY)
        }
        onPositionChanged: {
            var now = +new Date
            if (now - t < 50)
                return
            t = now
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            if (!(mainWindow.minimumHeight === mainWindow.height
                  && delta.y > 0)) {
                mainWindow.height -= delta.y
                mainWindow.y += delta.y
            }
        }
    }
    // 下
    MouseArea {
        property point clickPos: "0,0"
        property real t: 0
        x: maWidth
        width: parent.width - 2 * maWidth
        height: maWidth
        anchors.bottom: parent.bottom
        cursorShape: Qt.SizeVerCursor
        onPressed: {
            clickPos = Qt.point(mouseX, mouseY)
        }
        onPositionChanged: {
            var now = +new Date
            if (now - t < 50)
                return
            t = now
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            mainWindow.height += delta.y
        }
    }
    // 右上
    MouseArea {
        property point clickPos: "0,0"
        property real t: 0
        anchors.right: parent.right
        width: maWidth
        height: maWidth
        cursorShape: Qt.SizeBDiagCursor
        onPressed: {
            clickPos = Qt.point(mouseX, mouseY)
        }
        onPositionChanged: {
            var now = +new Date
            if (now - t < 50)
                return
            t = now
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            mainWindow.width += delta.x
            if (!(mainWindow.minimumHeight === mainWindow.height
                  && delta.y > 0)) {
                mainWindow.height -= delta.y
                mainWindow.y += delta.y
            }
        }
    }
    // 左下
    MouseArea {
        property point clickPos: "0,0"
        property real t: 0
        anchors.bottom: parent.bottom
        width: maWidth
        height: maWidth
        cursorShape: Qt.SizeBDiagCursor
        onPressed: {
            clickPos = Qt.point(mouseX, mouseY)
        }
        onPositionChanged: {
            var now = +new Date
            if (now - t < 50)
                return
            t = now
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            mainWindow.height += delta.y
            if (!(mainWindow.minimumWidth === mainWindow.width
                  && delta.x > 0)) {
                mainWindow.width -= delta.x
                mainWindow.x += delta.x
            }
        }
    }
    // 左上
    MouseArea {
        property point clickPos: "0,0"
        property real t: 0
        width: maWidth
        height: maWidth
        cursorShape: Qt.SizeFDiagCursor
        onPressed: {
            clickPos = Qt.point(mouseX, mouseY)
        }
        onPositionChanged: {
            var now = +new Date
            if (now - t < 50)
                return
            t = now
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            if (!(mainWindow.minimumWidth === mainWindow.width
                  && delta.x > 0)) {
                mainWindow.width -= delta.x
                mainWindow.x += delta.x
            }
            if (!(mainWindow.minimumHeight === mainWindow.height
                  && delta.y > 0)) {
                mainWindow.height -= delta.y
                mainWindow.y += delta.y
            }
        }
    }
    // 右下
    MouseArea {
        property point clickPos: "0,0"
        property real t: 0
        anchors.bottom: parent.bottom
        anchors.right: parent.right
        width: maWidth
        height: maWidth
        cursorShape: Qt.SizeFDiagCursor
        onPressed: {
            clickPos = Qt.point(mouseX, mouseY)
        }
        onPositionChanged: {
            var now = +new Date
            if (now - t < 100)
                return
            t = now
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            mainWindow.width += delta.x
            mainWindow.height += delta.y
        }
    }

    Component.onCompleted: {
        mainWindow = G.root
    }
}
