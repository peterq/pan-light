import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQml.Models 2.2
import "../js/app.js" as App
import "../js/util.js" as Util

Rectangle {
    id: root
    property var files: App.appState.fileList
    clip: true
    ListView {
        id: listView
        anchors.fill: parent
        model: files
        delegate: FileNode {
            meta: modelData
            idx: index
        }
        highlight: highlightComp
        highlightFollowsCurrentItem: true
        focus: true
        keyNavigationEnabled: true
        Keys.enabled: true
        Keys.onReturnPressed: {
            currentItem.handlePressEnter()
        }
        Keys.onMenuPressed: {
            currentItem.handlePressMenu()
        }
        Keys.onRightPressed: {
            App.forwardPath()
        }
        Keys.onLeftPressed: {
            App.backPath()
        }

        Connections {
            target: App.appState.mainWindow
            onCustomerEvent: {
                if (event != 'node.click') return
                listView.focus = true
                listView.highlightItem.shouldShow = true
                listView.currentIndex = data.index
            }
        }
        ScrollBar.vertical: ScrollBar {}
    }

    Component {
        id: highlightComp
        Rectangle {
            property var currentItem: listView.currentItem
            property bool shouldShow: false
            width:  currentItem &&
                    currentItem.width || 0
            height: currentItem &&
                    currentItem.height || 0
            color: Qt.rgba(140 / 255, 197 / 255, 1, .9)
            visible: listView.focus && shouldShow
            onCurrentItemChanged: {
                // 切换目录, 默认不显示
                if (currentItem === null) {
                    shouldShow = false
                } else if (listView.currentIndex !== 0) {
                    // currentItem 发生变化说明通过键盘操作了
                    shouldShow = true
                }
            }
        }
    }
}
