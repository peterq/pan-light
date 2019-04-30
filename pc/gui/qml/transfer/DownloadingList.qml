import QtQuick 2.7
import QtQuick.Controls 2.1
import QtQml.Models 2.2
import "../comps"
import "../js/app.js" as App
import "../js/util.js" as Util

Rectangle {
    id: root
    property var list: []
    property bool isFinish: false
    signal checkFid(var data)
    clip: true
    ListModel {
        id: model
    }
    ListView {
        id: listView
        visible: root.list.length > 0
        anchors.fill: parent
        model: list
        delegate: DownloadItem {
            id: item
            downloadId: modelData
            isFinish: root.isFinish
            Connections {
                target: root
                checkFid: {
                    if (item.meta.fileId === data.fid) {
                        data.exist = true
                    }
                }
            }
        }
        ScrollBar.vertical: ScrollBar {}
    }
    Text {
        anchors.centerIn: parent
        text: '暂时没有' + (isFinish ? '已完成':'下载中') + '的任务'
    }
}
