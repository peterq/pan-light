.pragma library
.import "./util.js" as Util

// 响应式状态机
var appState = (function () {
    var comp = Qt.createComponent('./appState.qml')
    return comp.createObject()
})()

// 加载首页
Util.event.once('init.api.ok', function (loginSession) {
    appState.loginSession = loginSession
    enterPath('/')
    console.log('app.js ', JSON.stringify(appState.loginSession))
})

// 历史记录后退
function backPath() {
    var dir = appState.accessDirHistory[appState.accessDirHistoryIndex - 1]
    if (!dir)
        return Util.Promise.reject('cant go back')
    var p = enterPath(dir, true)
    appState.accessDirHistoryIndex--
    return p
}

// 历史记录前进
function forwardPath() {
    var dir = appState.accessDirHistory[appState.accessDirHistoryIndex + 1]
    if (!dir)
        return Util.Promise.reject('cant go forward')
    var p = enterPath(dir, true)
    appState.accessDirHistoryIndex++
    return p
}

// 加载目录
function enterPath(path, dontRecordHistory) {
    if (appState.enterPathPromise)
        throw new Error('previous enter path promise is not finished')
    var promise = Util.callGoAsync('pan.ls', {
                                       "path": path
                                   })
    appState.enterPathPromise = promise
    promise.finally(function () {
        appState.enterPathPromise = null
    })
    Util.event.fire('pan.ls.promise', promise)
    appState.fileList = []
    appState.path = path
    promise.then(function (list) {
        appState.fileList = list
    })

    // 历史记录操作
    if (dontRecordHistory
            || path === appState.accessDirHistory[appState.accessDirHistoryIndex])
        return
    if (appState.accessDirHistoryIndex !== appState.accessDirHistory.length - 1) {
        for (var i = 0; i < appState.accessDirHistory.length
             - appState.accessDirHistoryIndex; i++) {
            appState.accessDirHistory.pop()
        }
    }
    appState.accessDirHistory.push(path)
    appState.accessDirHistoryIndex = appState.accessDirHistory.length - 1
}

function alert(title, msg, copyButton) {
    appState.alertPromise = appState.alertPromise.finally(function () {
        return Util.alert({
                              "parent": appState.mainWindow,
                              "title": title,
                              "msg": msg,
                              "copyButton": copyButton
                          })
    })
}
