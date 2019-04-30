import QtQuick 2.0
import "../js/app.js" as App
import "../js/util.js" as Util

Item {
    id: root
    HeaderBar {
        id: headerBar
    }
    Rectangle {
        width: parent.width
        height: 2
        color: 'gray'
        anchors.bottom: headerBar.bottom
    }
    DownloadList {
        id: downloadingList
        visible: headerBar.currentTab == '下载中'
        isFinish: false
        anchors.top: headerBar.bottom
        width: parent.width
        height: parent.height - headerBar.height
    }
    DownloadList {
        visible: headerBar.currentTab == '已完成'
        isFinish: true
        anchors.top: headerBar.bottom
        width: parent.width
        height: parent.height - headerBar.height
    }
    Component.onCompleted: {
        App.appState.transferComp = root
    }

    function addDownload(meta) {
        return Util.Promise.resolve().then(function () {
            var evt = {
                "fid": meta.fs_id + '',
                "exist": false
            }
            downloadingList.checkFid(evt)
            if (evt.exist)
                return Util.confirm()
            return true
        }).then(function () {
            if (App.appState.settings.defaultDownloadPath)
                return App.appState.settings.defaultDownloadPath
            return Util.pickSavePath({
                                         "fileName": meta.server_filename,
                                         "defaultFolder": App.appState.settings.lastDownloadPath
                                     })
        }).then(function (savePath) {
            savePath = savePath.toString()
            savePath = savePath.replace('file://', '')
            var id = Util.callGoSync('download.new', {
                                       "fid": meta.fs_id + '',
                                       "savePath": savePath
                                   })
            var obj = JSON.parse(JSON.stringify(meta))
            obj.downloadId = id
            var sep = Qt.platform.os == "windows" ? '\\' : '/'
            var t = String.prototype.split.call(savePath, sep)
            obj.saveName = t.pop()
            obj.savePath = savePath
            downloadingList.add(obj)
        })
    }
}
