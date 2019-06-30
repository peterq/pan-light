import QtQuick 2.9
import QtQuick.Window 2.2
import QtQuick.Controls 1.4 as Controls
import Qt.labs.platform 1.0
import QtGraphicalEffects 1.0
import "../js/util.js" as Util
import "../js/global.js" as G
import "../js/app.js" as App

Window {
    id: root
    visible: true
    width: contentContainer.width + 20
    height: contentContainer.height + 20
    title: 'pan-light float'
    x: Screen.desktopAvailableWidth - width
    y: 100
    //无边框的window flags
    flags: Qt.WA_TranslucentBackground | Qt.WA_TransparentForMouseEvents
           | Qt.FramelessWindowHint | Qt.WindowSystemMenuHint
           | Qt.WindowStaysOnTopHint | Qt.X11BypassWindowManagerHint
    color: 'transparent'

    Component.onCompleted: {
        App.appState.floatWindow = root
    }

    onVisibleChanged: {
        if (visible) {
            if (x < 0) x = 100
            if (y < 0) y = 100
            x = Math.min(x, Screen.desktopAvailableWidth - width)
            y = Math.min(y, Screen.desktopAvailableHeight - height)
        }
    }

    DataSaver {
        $key: 'window.float'
        property alias x: root.x
        property alias y: root.y
        property alias visible: root.visible
    }

    Rectangle {
        id: contentContainer
        width: 130
        height: 35
        border.color: 'gray'
        // radius: 5
        clip: true
        anchors.centerIn: parent

        Rectangle {
            id: logo
            x: 1
            height: parent.height - 2
            anchors.verticalCenter: parent.verticalCenter
            width: height * 1.1
            color: '#38f'
            IconFont {
                type: 'baidu-cloud'
                width: parent.height * 0.8
                color: 'white'
                anchors.centerIn: parent
            }
        }
        Rectangle {
            anchors.left: logo.right
            anchors.verticalCenter: parent.verticalCenter
            height: parent.height - 2
            width: parent.width - 2 - logo.width
            color: '#a2dcf4'
            Text {
                text: ''
                color: '#34658a'
                anchors.centerIn: parent
                Timer {
                    interval: 1000
                    running: root.visible
                    triggeredOnStart: true
                    repeat: true
                    onTriggered: {
                        var data = {
                            "speed": 0
                        }
                        App.appState.transferComp.sumSpeed(data)
                        parent.text = Util.humanSize(data.speed) + '/s'
                    }
                }
            }
        }
        OpacityMask {
            anchors.fill: logo
            source: logo
            maskSource: contentContainer
            visible: true
            antialiasing: true
        }

        MouseArea {
            property point clickPos: "0,0"
            anchors.fill: parent
            onPressed: {
                root.requestActivate()
                clickPos = Qt.point(mouseX, mouseY)
            }
            onPositionChanged: {
                var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
                root.x += delta.x
                root.y += delta.y
            }
            acceptedButtons: Qt.LeftButton | Qt.RightButton
            onClicked: {
                if (mouse.button === Qt.RightButton) {
                    contentMenu.popup()
                }
            }
        }
    }

    DropShadow {
        anchors.fill: contentContainer
        horizontalOffset: -5
        verticalOffset: -5
        radius: 12.0
        samples: 25
        color: "#20000000"
        spread: 0.0
        source: contentContainer
    }
    DropShadow {
        anchors.fill: contentContainer
        horizontalOffset: 5
        verticalOffset: 5
        radius: 12.0
        samples: 25
        color: "#20000000"
        spread: 0.0
        source: contentContainer
    }

    // 右键菜单
    Controls.Menu {
        id: contentMenu
        Controls.MenuItem {
            id: hideItem
            text: '隐藏悬浮窗'
            onTriggered: {
                root.hide()
            }
        }
        Controls.MenuItem {
            text: '显示主界面'
            visible: !G.root.visible
            onTriggered: {
                G.root.show()
                G.root.raise()
                G.root.requestActivate()
            }
        }
        Controls.MenuItem {
            text: '退出程序'
            onTriggered: Util.exit()
        }
    }

    // 系统托盘菜单
    Menu {
        id: systemTrayMenu
        MenuItem {
            visible: root.visible
            text: '隐藏悬浮窗'
            onTriggered: root.hide()
        }
        MenuItem {
            text: '显示悬浮窗'
            visible: !root.visible
            onTriggered: {
                root.show()
            }
        }
        MenuItem {
            text: '显示主界面'
            onTriggered: {
                G.root.show()
                G.root.raise()
                G.root.requestActivate()
            }
        }
        MenuItem {
            text: '退出程序'
            onTriggered: Util.exit()
        }
    }
    // 系统托盘
    SystemTrayIcon {
        id: trayIcon
        visible: true
        iconSource: "../assets/images/pan-light-1.png"
        tooltip: "pan-light tray"
        onActivated: {
            root.show()
            root.raise()
            root.requestActivate()
        }
        menu: systemTrayMenu
    }
}
