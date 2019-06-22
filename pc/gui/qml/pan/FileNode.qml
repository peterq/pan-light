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
    property bool isSmallFile: !meta.isdir && meta.size < 256 * 1024
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
        type: (meta.isdir && 'dir') || fileItem.meta.server_filename.split(
                  '.').pop()
    }
    Text {
        id: filenameText
        text: fileItem.meta.server_filename
        anchors.verticalCenter: parent.verticalCenter
        anchors.left: fileIcon.right
        anchors.leftMargin: 5
        width: parent.width * 0.6
        elide: Text.ElideRight
    }

    Text {
        text: meta.isdir ? '-' : Util.humanSize(fileItem.meta.size)
        anchors.verticalCenter: parent.verticalCenter
        anchors.left: filenameText.right
        anchors.leftMargin: 20
        width: 100
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
        var dirMenu = [{
                           "name": '进入',
                           "cb": function () {
                               App.enterPath(fileItem.meta.path)
                           }
                       }, {
                           "name": '添加至快捷导航',
                           "cb": function () {
                               var p = App.prompt('请输入快捷方式名称', function (str) {
                                   if (str === '')
                                       return '请输入名称'
                                   return true
                               }, fileItem.meta.server_filename)
                               p.then(function (name) {
                                   App.addPathCollection({
                                                             "name": name,
                                                             "path": fileItem.meta.path
                                                         })
                               })
                           }
                       }]
        if (fileItem.meta.isdir) {
            fileItem.menus = dirMenu
            return
        }

        var fileMenu = [{
                            "name": '直接下载',
                            "cb": function () {
                                fileItem.clickDownload()
                            }
                        }, {
                            "name": 'vip通道下载',
                            "cb": function () {
                                fileItem.clickDownloadViaVip()
                            },
                            "hide": isSmallFile
                        }, {
                            "name": '分享到资源广场',
                            "cb": function () {
                                fileItem.clickShare()
                            },
                            "hide": isSmallFile
                        }, {
                            "name": '播放',
                            "cb": function () {
                                Util.playVideo(fileItem.meta, false)
                            },
                            "hide": !Util.isVideo(fileItem.meta.server_filename)
                        }, {
                            "name": 'vip通道播放',
                            "cb": function () {
                                Util.playVideo(fileItem.meta, true)
                            },
                            "hide": isSmallFile || !Util.isVideo(
                                        fileItem.meta.server_filename)
                        }]
        fileItem.menus = fileMenu.filter(function (v) {
            return !v.hide
        })
    }
    function clickDownload() {
        console.log('down')
        App.appState.transferComp.addDownload(fileItem.meta)
    }
    function clickDownloadViaVip() {
        console.log('vip down')
        App.appState.transferComp.addDownload(fileItem.meta, true)
    }
    function clickShare() {
        Util.openShare(fileItem.meta)
    }
}
