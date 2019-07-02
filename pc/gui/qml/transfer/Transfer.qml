import QtQuick 2.0
import "../js/app.js" as App
import "../js/util.js" as Util

Item {
    id: root

    signal sumSpeed(var data)
    signal active

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
        id: downloadedList
        visible: headerBar.currentTab == '已完成'
        isFinish: true
        anchors.top: headerBar.bottom
        width: parent.width
        height: parent.height - headerBar.height
    }
    Component.onCompleted: {
        App.appState.transferComp = root
    }

    function addDownload(meta, useVip) {
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
            var newFid = (useVip ? 'vip' : 'direct') + '.' + meta.fs_id
            var id = Util.callGoSync('download.new', {
                                         "fid": newFid,
                                         "savePath": savePath
                                     })
            var obj = JSON.parse(JSON.stringify(meta))
            obj.downloadId = id
            var sep = Qt.platform.os == "windows" ? '\\' : '/'
            var t = String.prototype.split.call(savePath, sep)
            obj.newFid = newFid
            obj.saveName = t.pop()
            obj.savePath = savePath
            obj.useVip = !!useVip
            downloadingList.add(obj)
        })
    }

    function addDownloadShare(md5, sliceMd5, fileSize, fileName) {
        var fid = ['share', md5, sliceMd5, fileSize].join('.')
        return Util.Promise.resolve().then(function () {
            return Util.pickSavePath({
                                         "fileName": fileName,
                                         "defaultFolder": App.appState.settings.lastDownloadPath
                                     })
        }).then(function (savePath) {
            var id = Util.callGoSync('download.new', {
                                         "fid": fid,
                                         "savePath": savePath
                                     })
            var obj = {
                "size": fileSize,
                "downloadId": id,
                "newFid": fid,
                "savePath": savePath,
                "useVip": true,
                "fileName": fileName,
                "isShare": true
            }
            var sep = Qt.platform.os == "windows" ? '\\' : '/'
            var t = String.prototype.split.call(savePath, sep)
            obj.saveName = t.pop()
            downloadingList.add(obj)
        })
    }

    function deleteItem(idx, isFinish) {
        var c = isFinish ? downloadedList : downloadingList
        c.remove(idx)
    }

    function itemCompleted(idx) {
        var data = JSON.parse(JSON.stringify(downloadingList.get(idx)))
        downloadingList.remove(idx)
        downloadedList.add(data)
    }
}
