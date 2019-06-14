import QtQuick 2.0
import QtGraphicalEffects 1.0
import QtQuick.Window 2.2
import "../comps"
import "../js/global.js" as G

Item {
    property Component content
    property int shadeWidth: isMax ? 0 : 10
    property bool isMax: (mainWindow.visibility === Window.Maximized)
    property var mainWindow
    width: parent.width
    height: parent.height


    Resize {
        width: contentContainer.width + 2 * maWidth
        height: contentContainer.height + 2 * maWidth
        anchors.centerIn: parent
    }
    Rectangle {
        id: contentContainer
        border.color: 'gray'
        property int subInt: 2 * shadeWidth * 1.1
        width: parent.width - subInt
        height: parent.height - subInt
        x: subInt / 2
        y: subInt / 2

        Loader {
            width: parent.width - 2
            height: parent.height - 2
            anchors.centerIn: parent
            focus: true
            sourceComponent: content
        }

        Row {
            anchors.right: parent.right
            anchors.rightMargin: 5
            y: 5
            spacing: 10
            IconButton {
                iconType: 'more-down'
                title: '更多'
                width: 20
            }
            IconButton {
                iconType: 'min'
                title: '最小化'
                width: 20
                onClicked: {
                    G.root.visibility = Window.Minimized
                }
            }
            IconButton {
                iconType: 'max'
                title: '最大化'
                visible: (mainWindow.visibility !== Window.Maximized)
                width: 20
                onClicked: {
                    G.root.visibility = Window.Maximized
                }
            }
            IconButton {
                iconType: 'normal-size'
                title: '恢复'
                visible: (mainWindow.visibility !== Window.Windowed)
                onClicked: {
                    G.root.visibility = Window.Windowed
                }
            }
            IconButton {
                iconType: 'close'
                title: '关闭窗口'
                color: 'red'
                width: 20
                onClicked: {
                    G.root.visible = false
                }
            }
        }
    }

    //    DropShadow {
    //            anchors.fill: contentContainer
    //            radius: shadeWidth
    //            samples: shadeWidth
    //            spread: 0.1
    //            color: "#80000000"
    //            source: contentContainer
    //    }
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
    Component.onCompleted: {
        mainWindow = G.root
    }
}
