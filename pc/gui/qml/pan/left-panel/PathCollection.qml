import QtQuick 2.0
import QtQuick.Controls 2.1
import "../../js/app.js" as App
import "../../comps"
import "../../widget"

Column {
    id: root
    width: parent.width
    property var collection: App.appState.pathCollection

    // title
    Rectangle {
        color: '#a8defb'
        height: 40
        width: parent.width
        Text {
            text: '快速导航'
            anchors.centerIn: parent
        }

        Rectangle {
            width: parent.width
            height: 1
            color: '#00a6ff'
            anchors.top: parent.top
        }
        Rectangle {
            width: parent.width
            height: 1
            color: '#00a6ff'
            anchors.bottom: parent.bottom
        }
        Button {
            id: iconBtn
            width: 25
            height: width
            anchors.verticalCenter: parent.verticalCenter
            anchors.right: parent.right
            anchors.rightMargin: 10
            ToolTip {
                show: iconBtn.hovered
                text: '清空'
            }
            IconFont {
                id: icon
                type: 'delete'
                width: parent.width
                color: iconBtn.hovered ? 'red' : Qt.lighter('red')
            }
            display: AbstractButton.IconOnly
            background: Item {
            }
            onClicked: {
                App.clearPathCollection()
            }
        }
    }

    Rectangle {
        width: parent.width
        height: parent.height - 40
        color: '#c3eaff'
        clip: true
        Rectangle {
            width: parent.width
            height: 1
            color: '#00a6ff'
            anchors.bottom: parent.bottom
        }
        // 为空提示
        Text {
            anchors.centerIn: parent
            visible: collection.length === 0
            width: parent.width * 0.9
            wrapMode: Text.Wrap
            horizontalAlignment: Text.AlignHCenter
            text: '还没有内容, 你可以右键单击右侧文件夹, 将其添加至此处'
        }
        // list view
        ListView {
            id: listView
            anchors.fill: parent
            visible: collection.length > 0
            model: listModel

            delegate: PathCollectionItem {
                id: collectionItem
                meta: listModel.get(index)
                idx: index
                parentList: listView
            }

            displaced: Transition {
                NumberAnimation {
                    properties: "x,y"
                    duration: 200
                }
            }
            move: Transition {
                NumberAnimation {
                    properties: "x,y"
                    duration: 200
                }
            }

            ScrollBar.vertical: ScrollBar {
            }
        }

        ListModel {
            id: listModel
            Component.onCompleted: {
                App.appState.pathCollectionModel = listModel
                collection.forEach(function (option) {
                    listModel.append(option)
                })
            }
        }
    }
}
