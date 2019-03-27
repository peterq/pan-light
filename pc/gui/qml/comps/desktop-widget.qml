import QtQuick 2.9
import QtQuick.Window 2.2
import QtQuick.Controls 1.4 as Controls
import Qt.labs.platform 1.0
import '../js/util.js' as Util


//包含系统托盘的导包
Window {
    id: mainWindow
    //大小根据屏幕计算，宽高比为6:14
    visible: true
    minimumHeight: 50
    minimumWidth: 120
    width: Screen.desktopAvailableWidth / 14
    height: width * 3 / 7
    title: qsTr("tiny monitor")
    x: 100
    y: 100
    //无边框的window flags
    flags: Qt.WA_TranslucentBackground | Qt.WA_TransparentForMouseEvents| Qt.FramelessWindowHint
                      | Qt.WindowSystemMenuHint
                       | Qt.WindowStaysOnTopHint | Qt.X11BypassWindowManagerHint
    //灰色0.9透明度
//    color: Qt.rgba(0,0,0,0)
    color: 'transparent'

    Component.onCompleted: {
        Util.bridge.changeAttribute(mainWindow, Qt.WA_TranslucentBackground, true)
        Util.bridge.changeAttribute(mainWindow, Qt.WA_TransparentForMouseEvents, true)
    }

    Rectangle {
      anchors.fill: parent
      color: Qt.rgba(0,0,0,0)
      border.color: 'red'
      border.width: 2
    }
    Rectangle {
        id: rectangle
        x: 0
        y: 0
        width: mainWindow.height
        height: width
        color: Qt.rgba(0.2, 1.0, 0.0, 0.7)
        Component.onCompleted: {

        }
    }

//    混合动画效果（这里混合x和y轴平移
    ParallelAnimation {
        id: moveAnimation
        running: false
        PropertyAnimation {
            target: mainWindow
            property: 'x'
            easing.type: Easing.Linear
            duration: 100
        }
        PropertyAnimation {
            target: mainWindow
            property: 'y'
            easing.type: Easing.Linear
            duration: 100
        }
    }

    //鼠标可控制区域
    MouseArea {
        property point clickPos: "0,0"
        id: dragRegion
        anchors.fill: rectangle
        drag.minimumX: 0
        drag.maximumX: Screen.desktopAvailableWidth - mainWindow.width
        drag.minimumY: 0
        drag.maximumY: Screen.desktopAvailableHeight - mainWindow.heigh
        onPressed: {
            mainWindow.requestActivate()
            clickPos = Qt.point(mouseX, mouseY)
        }

        onPositionChanged: {
            moveAnimation.stop()
            //鼠标偏移量
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
//            console.log(delta.x + "  " + delta.y)
            mainWindow.x += delta.x
            mainWindow.y += delta.y
            moveAnimation.start()
        }
        //添加右键菜单
        acceptedButtons: Qt.LeftButton | Qt.RightButton // 激活右键（别落下这个）
        onClicked: {
            if (mouse.button === Qt.RightButton) {
                // 右键菜单
                contentMenu.popup()
            }
        }
    }
    //不是托盘的菜单类
    Controls.Menu {
        id: contentMenu
        // 右键菜单
        Controls.MenuItem {
            id:hideItem
            text: qsTr("隐藏")
            onTriggered: {
                if(trayIcon==null){
                    console.log("系统托盘不存在");
                    contentMenu.removeItem(hideItem);
                    return;
                }else{
                    if(trayIcon.available){
                        console.log("系统托盘存在");
                    }else{
                        console.log("系统托盘不存在");
                        contentMenu.removeItem(hideItem)
                    }
                }
                mainWindow.hide()
            }
        }
        Controls.MenuItem {
            text: qsTr("退出")
            onTriggered: Qt.quit()
        }
    }

    //使用系统托盘的菜单组件
    Menu {
        id: systemTrayMenu
        // 右键菜单
        MenuItem {
            text: qsTr("隐藏")
            shortcut: "Ctrl+z"
            onTriggered: mainWindow.hide()
        }
        MenuItem {
            text: qsTr("显示")
            onTriggered: {
               mainWindow.show()
               mainWindow.flags = Qt.WA_TranslucentBackground
                   | Qt.WA_TransparentForMouseEvents| Qt.FramelessWindowHint
                   | Qt.WindowSystemMenuHint
                   | Qt.WindowStaysOnTopHint | Qt.X11BypassWindowManagerHint
            }
        }
        MenuItem {
            text: qsTr("退出")
            onTriggered: Qt.quit()
        }
    }
    //系统托盘
    SystemTrayIcon {
        id:trayIcon
        visible: true
        iconSource: "../assets/images/pan-light-1.png"
        tooltip: "tiny-流量监控软件"
        onActivated: {
            mainWindow.show()
            mainWindow.raise()
            mainWindow.requestActivate()
        }
        menu: systemTrayMenu
    }
}
