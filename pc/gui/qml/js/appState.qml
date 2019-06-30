import QtQuick 2.0
import "../comps"
import "util.js" as Util

Item {
    id: appState
    property var loginSession: null
    property string path: '/'
    property var pathInfo: [{
            "path": '/',
            "name": '全部文件'
        }]
    property var enterPathPromise: null
    property var fileList: []
    property var accessDirHistory: []
    property int accessDirHistoryIndex: -1
    property var player: null
    property var mainWindow: null
    property var alertPromise: Util.Promise.resolve()
    property var pathCollection: []
    property var pathCollectionModel: null
    property var downloadingList: []
    property var completedList: []
    property var transferComp: null
    property var downloadingListComp: null
    property var floatWindow: null
    property var globalRightMenu: null
    property alias settings: settings

    DataSaver {
        $key: 'app-state'
        property alias pathCollection: appState.pathCollection
        property alias completedList: appState.completedList
        property alias downloadingList: appState.downloadingList
    }

    DataSaver {
        id: settings
        $key: 'app-settings'
        property string defaultDownloadPath: ''
        property string lastDownloadPath: ''
        property int maxParallelTaskNumber: 3
        property int maxParallelCorutineNumber: 64

        Component.onCompleted: {
            Util.callGoSync('config', {
                                "maxParallelCorutineNumber": maxParallelCorutineNumber
                            })
        }
    }

    onPathChanged: {
        if (path === '/') {
            pathInfo = [{
                            "path": '/',
                            "name": '全部文件'
                        }]
            return
        }
        var dirs = path.split('/')
        var arr = []
        var p = ''
        dirs.forEach(function (dir, idx) {
            p += idx === 1 ? dir : '/' + dir
            arr.push({
                         "path": p,
                         "name": idx === 0 ? '全部文件' : dir
                     })
        })
        pathInfo = arr
    }

    Component.onCompleted: {
        loginSession = {
            "username": '用户名加载中...',
            "photo": ''
        }
    }
}
