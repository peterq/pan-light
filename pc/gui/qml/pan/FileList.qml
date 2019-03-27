import QtQuick 2.0
import QtQml.Models 2.2
import '../js/app.js' as App
import '../js/util.js' as Util

Rectangle {
    id: root
    property var files: App.appState.fileList
    clip: true

    DelegateModel {
        id: visualModel
        model: root.files
        delegate: Rectangle {
            id: fileItem
            property string prefix: meta.isdir ? '文件夹: ' :'文件  : '
            property var meta: modelData
            property var menus: []
            property bool hover: fileItemMa.containsMouse
            height: 40
            width: root.width
            // hover 底色
            Rectangle {
                visible: fileItem.hover
                anchors.fill: parent
                color: Qt.rgba(140/255, 197/255, 1, .4)
            }
            Text {
                text: fileItem.prefix + fileItem.meta.server_filename
                anchors.verticalCenter: parent.verticalCenter
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
                }
                onDoubleClicked: {
                    if (mouse.button === Qt.RightButton) return
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
                    fileItem.menus = [
                        {name: '直接下载', cb: function (){fileItem.clickDownload()}},
                        {name: 'vip通道下载', cb: function (){fileItem.clickDownloadViaVip()}}
                    ]
                    if (Util.isVideo(fileItem.meta.server_filename)) {
                        fileItem.menus = fileItem.menus.concat([
                            {name: '播放', cb: function (){Util.playVideo(fileItem.meta, false)}},
                            {name: 'vip通道播放', cb: function (){Util.playVideo(fileItem.meta, true)}}
                        ])
                    }
                }
            }
            function clickDownload(){
                console.log('down')
            }
            function clickDownloadViaVip() {
                console.log('vip down')
            }


        }
    }
    ListView {
        id:listView;
        anchors.fill: parent
        model: visualModel
    }
}
