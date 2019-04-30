import QtQuick 2.0
import QtQuick.Controls 2.1
import "../comps"
import "../js/app.js" as App
import "../js/util.js" as Util

Rectangle {
    id: root
    property var appState: App.appState
    property bool isFinish: false
    signal checkFid(var data)
    clip: true
    ListModel {
        id: listModel
    }
    ListView {
        id: listView
        visible: listModel.count > 0
        anchors.fill: parent
        model: listModel
        delegate: DownloadItem {
            id: item
            meta: listModel.get(index)
            isFinish: root.isFinish
            Connections {
                target: root
                onCheckFid: {
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
        visible: listModel.count === 0
        text: '暂时没有' + (isFinish ? '已完成':'下载中') + '的任务'
    }

    Component.onCompleted: {
        Util.arrToListModel(list(), listModel)
//        updateList(Util.listModelClear(listModel))
    }

    function list() {
        if (isFinish)
            return appState.completedList
        return appState.downloadingList
    }

    function updateList(l) {
        if (isFinish)
            appState.completedList = l
        else appState.downloadingList = l
    }

    function add(data) {
        updateList(Util.listModelAdd(listModel, data))
    }
}
