import QtQuick 2.0
import QtQuick.Controls 2.1
import "../../js/app.js" as App
import "../../js/util.js" as Util
import "../../widget"
import "../../comps"

Item {
    id: root
    property var meta
    property int idx
    property var parentList
    width: parent.width
    height: 50

    function showMenu() {
        var menu = []
        if (idx !== 0) {
            menu.push({
                          "name": '置顶',
                          "cb": function () {
                              App.movePathCollectionItem(idx, 0)
                          }
                      })
        }
        menu.push({
                      "name": '进入',
                      "cb": function () {
                          App.enterPath(meta.path)
                      }
                  })
        menu.push({
                      "name": '删除',
                      "cb": function () {
                          App.removePathCollection(idx)
                      }
                  })
        Util.showMenu(menu)
    }
    Rectangle {
        id: itemContent
        width: root.width
        height: root.height
        anchors.horizontalCenter: root.horizontalCenter
        anchors.verticalCenter: root.verticalCenter
        MouseArea {
            id: itemMa
            anchors.fill: parent
            hoverEnabled: true
            acceptedButtons: Qt.LeftButton | Qt.RightButton
            onClicked: {
                if (mouse.button === Qt.RightButton) {
                    root.showMenu()
                } else {
                    App.enterPath(meta.path)
                }
            }
        }
        Rectangle {
            anchors.fill: parent
            color: itemContent.Drag.active ? Qt.darker('#c1dbff') : '#c1dbff'
        }

        Rectangle {
            width: parent.width
            height: 1
            color: Qt.lighter('#00a6ff')
            anchors.bottom: parent.bottom
        }
        Rectangle {
            width: parent.width
            height: 1
            color: Qt.lighter('#00a6ff')
            anchors.top: parent.top
        }

        ToolTip {
            text: meta.path
            show: itemMa.containsMouse
        }
        Item {
            id: orderBtn
            width: 25
            height: width
            anchors.verticalCenter: parent.verticalCenter
            anchors.left: parent.left
            anchors.leftMargin: 5
            visible: itemMa.containsMouse || orderBtnMa.containsMouse
                     || orderBtnMa.held
            MouseArea {
                id: orderBtnMa
                anchors.fill: parent
                hoverEnabled: true
                drag.target: itemContent
                drag.axis: Drag.YAxis // 只允许沿Y轴拖动
                property bool held: false
                onPressed: held = true
                onReleased: held = false

                onMouseYChanged: {
                    if (!drag.active)
                        return
                    var p = orderBtnMa.mapToItem(root.parentList, width / 2,
                                                 height / 2)
                    var targetIndex = root.parentList.indexAt(
                                p.x, p.y + root.parentList.contentY)
                    if (targetIndex > -1)
                        App.movePathCollectionItem(root.idx, targetIndex)
                }
                drag.onActiveChanged: {

                }
            }
            ToolTip {
                show: orderBtnMa.containsMouse && !orderBtnMa.held
                text: '拖动以排序'
            }
            IconFont {
                id: orderIcon
                type: 'order'
                width: parent.width
            }
        }

        Text {
            text: meta.name
            anchors.verticalCenter: parent.verticalCenter
            anchors.left: parent.left
            anchors.leftMargin: 35
            width: parent.width - 50
            elide: Text.ElideRight
        }

        states: [
            State {
                when: itemContent.Drag.active
                ParentChange {
                    target: itemContent
                    parent: root.parentList
                }
                AnchorChanges {
                    target: itemContent
                    anchors.horizontalCenter: undefined
                    anchors.verticalCenter: undefined
                }
            }
        ]

        Drag.active: orderBtnMa.drag.active
    }
}
