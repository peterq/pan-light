import QtQuick 2.0

QtObject {
    id: appState
    property var loginSession: null
    property string path: '/'
    property var enterPathPromise: null
    property var fileList: []
    property var accessDirHistory: []
    property int accessDirHistoryIndex: -1
    property var player: null
    property var mainWindow: null

    Component.onCompleted: {
        loginSession = {
            username: '用户名加载中...',
            photo: ''
        }
    }
}
