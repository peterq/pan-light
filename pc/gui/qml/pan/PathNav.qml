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
    Row {
        anchors.left: btns.right
        anchors.verticalCenter: parent.verticalCenter
        height: parent.height
        spacing: 0
        Repeater {
            model: App.appState.pathInfo
            Label {
                id: dirname
                text: modelData.name
                elide: Text.ElideMiddle
                width: Math.min(implicitWidth, 200) + sep.implicitWidth
                anchors.verticalCenter: parent.verticalCenter
                color: pathMa.containsMouse ? '#5c9fff' : 'black'
                visible: index === 0 || App.appState.pathInfo.length < 5 || index >= App.appState.pathInfo.length - 3
                ToolTip {
                    text: modelData.name
                    show: pathMa.containsMouse
                }
                MouseArea {
                    id: pathMa
                    width: parent.width - sep.implicitWidth
                    height: parent.height
                    hoverEnabled: true
                    cursorShape: Qt.PointingHandCursor
                    onClicked: {
                        App.enterPath(modelData.path)
                    }
                }
                Label {
                    id: sep
                    visible: index !== App.appState.pathInfo.length - 1
                    text: App.appState.pathInfo.length >= 5 && index === 0 ? ' > ··· > ' : ' > '
                    anchors.right: parent.right
                    anchors.verticalCenter: parent.verticalCenter
                    color: '#3887ff'
                }
            }
        }
    }
}
