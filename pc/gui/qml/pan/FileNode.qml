import QtQuick 2.0
import "../js/app.js" as App
import "../js/util.js" as Util

Item {
    id: fileItem
    property string prefix: meta.isdir ? '文件夹: ' : '文件  : '
    property var meta
    property int idx
    property var menus: []
    property bool hover: fileItemMa.containsMouse
    height: 50
    width: root.width


    function handlePressEnter() {
        if (meta.isdir)
            App.enterPath(meta.path)
        else {
            handlePressMenu()
        }
    }

    function handlePressMenu() {
        var p = fileItem.mapToGlobal(width / 2, height / 2)
        Util.bridge.setCursorPos(p.x, p.y)
        Util.showMenu(menus)
    }
    // hover 底色
    Rectangle {
        visible: fileItem.hover
        anchors.fill: parent
        color: Qt.rgba(140 / 255, 197 / 255, 1, .4)
    }
    FileIcon {
        id: fileIcon
        width: parent.height
        height: width
        anchors.verticalCenter: parent.verticalCenter
        anchors.leftMargin: 10
        type: (meta.isdir && 'dir') || fileItem.meta.server_filename.split('.').pop()
    }
    Text {
        text: fileItem.meta.server_filename
        anchors.verticalCenter: parent.verticalCenter
        anchors.left: fileIcon.right
        anchors.leftMargin: 5
        anchors.right: parent.right
        anchors.rightMargin: 5
        elide: Text.ElideRight
    }

    MouseArea {
        id: fileItemMa
        hoverEnabled: true
        anchors.fill: parent
        acceptedButtons: Qt.LeftButton | Qt.RightButton
        onClicked: {
            if (mouse.button === Qt.RightButton && fileItem.menus.length) {
                Util.showMenu(fileItem.menus)
            }
            App.appState.mainWindow.customerEvent('node.click', {
                                                      "index": fileItem.idx
                                                  })
        }
        onDoubleClicked: {
            if (mouse.button === Qt.RightButton)
                return
            if (fileItem.meta.isdir) {
                App.enterPath(fileItem.meta.path)
            }
        }
    }
    Rectangle {
        width: parent.width
        height: 1
        color: Qt.lighter('gray')
        anchors.bottom: parent.bottom
    }
    Component.onCompleted: {
        if (!fileItem.meta.isdir) {
            fileItem.menus = [{
                                  "name": '直接下载',
                                  "cb": function () {
                                      fileItem.clickDownload()
                                  }
                              }, {
                                  "name": 'vip通道下载',
                                  "cb": function () {
                                      fileItem.clickDownloadViaVip()
                                  }
                              }]
            if (Util.isVideo(fileItem.meta.server_filename)) {
                fileItem.menus = fileItem.menus.concat([{
                                                            "name": '播放',
                                                            "cb": function () {
                                                                Util.playVideo(
                                                                            fileItem.meta,
                                                                            false)
                                                            }
                                                        }, {
                                                            "name": 'vip通道播放',
                                                            "cb": function () {
                                                                Util.playVideo(
                                                                            fileItem.meta,
                                                                            true)
                                                            }
                                                        }])
            }
        } else {
            fileItem.menus.push({
                                    "name": '进入',
                                    "cb": function () {
                                        App.enterPath(fileItem.meta.path)
                                    }
                                })
        }
    }
    function clickDownload() {
        console.log('down')
    }
    function clickDownloadViaVip() {
        console.log('vip down')
    }
}
