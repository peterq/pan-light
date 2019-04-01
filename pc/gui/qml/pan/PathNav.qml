import QtQuick 2.0
import QtQuick.Controls 2.2
import '../comps'
import '../widget'
import '../js/app.js' as App

Item {
    Row {
        id: btns
        spacing: 10
        width: (parent.height * 0.6 + spacing) * 3
        height: parent.height
        Button {
            property bool canReturn: App.appState.accessDirHistoryIndex > 0
            flat: true
            height: parent.height * 0.6
            width: height
            anchors.verticalCenter: parent.verticalCenter
            IconFont {
                width: parent.width
                type: parent.canReturn ? 'return-enable' : 'return'
            }
            onClicked: {
                if (canReturn) {
                    App.backPath()
                }
            }
        }
        Button {
            property bool canEnter: App.appState.accessDirHistoryIndex < App.appState.accessDirHistory.length - 1
            flat: true
            height: parent.height * 0.6
            width: height
            anchors.verticalCenter: parent.verticalCenter
            IconFont {
                width: parent.width
                type: parent.canEnter ? 'enter-enable' : 'enter'
            }
            onClicked: {
                if (canEnter) {
                    App.forwardPath()
                }
            }
        }
        Button {
            flat: true
            height: parent.height * 0.6
            width: height
            anchors.verticalCenter: parent.verticalCenter
            IconFont {
                width: parent.width
                type: 'home'
            }
            onClicked: {
                App.enterPath('/')
            }
        }
    }

    Label {
        anchors.left: btns.right
        anchors.verticalCenter: parent.verticalCenter
        text: '当前路径: ' + App.appState.path
        ToolTip {
            text: parent.text
            show: pathMa.containsMouse
        }
        MouseArea {
            id: pathMa
            anchors.fill: parent
            hoverEnabled: true
        }
    }
}
