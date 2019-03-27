
import QtQuick 2.0
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0

// 桌面悬浮挂件
ApplicationWindow {
    id: widgetWin
    visible: true
    width: 640
    height: 480
    title: '悬浮任务栏'
    color: "transparent"
    flags: Qt.FramelessWindowHint | Qt.WindowSystemMenuHint
            | Qt.WindowStaysOnTopHint | Qt.X11BypassWindowManagerHint
    Rectangle {
        id: r
        width: 300
        height: 200
        //灰色0.9透明度
        color:Qt.rgba(0.5,0.5,0.5,0.9)
        MouseArea {
            id: dragRegion
            anchors.fill: parent
            property point clickPos: "0,0"
            onPressed: {
                clickPos  = Qt.point(mouse.x,mouse.y)
                }
            onPositionChanged: {
                r.x += 1
                //鼠标偏移量
                var delta = Qt.point(mouse.x-clickPos.x, mouse.y-clickPos.y)
                //如果mainwindow继承自QWidget,用setPos
                widgetWin.setX(widgetWin.x+delta.x)
                widgetWin.setY(widgetWin.y+delta.y)
            }
        }
    }
}
