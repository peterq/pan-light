.pragma library

.import "./util.js" as Util
.import "./global.js" as G

// 响应式状态机
var appState = (function () {
    var comp = Qt.createComponent('./appState.qml')
    return comp.createObject()
})()

G.appState = appState

// 加载首页
Util.event.once('init.api.ok', function (loginSession) {
    appState.loginSession = loginSession
    enterPath('/')
    console.log('app.js ', JSON.stringify(appState.loginSession))
    Util.api("refresh-token").then(function (token) {
        console.log('refresh token ok')
    }).catch(function (err) {
        console.log('refresh token error:', err.message)
        Util.callGoAsync("api.login").then(function (token) {
            console.log('登录 pan-light server 成功')
        })
    })
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

function prompt(msg, checkFunc, content) {
    return new Util.Promise(function (resolve, reject) {
        appState.alertPromise = appState.alertPromise.finally(function () {
            return Util.prompt({
                                   "parent": appState.mainWindow,
                                   "title": '请输入',
                                   "msg": msg,
                                   "checkFunc": checkFunc,
                                   "content": content
                               }).then(resolve, reject)
        })
    })
}

function syncPathCollection() {
    var arr = []
    for (var i = 0; i < appState.pathCollectionModel.count; i++) {
        arr.push(appState.pathCollectionModel.get(i))
    }
    appState.pathCollection = arr
}

function addPathCollection(option) {
    appState.pathCollectionModel.append(option)
    syncPathCollection()
}

function clearPathCollection(option) {
    appState.pathCollectionModel.clear()
    appState.pathCollection = []
}

function removePathCollection(index) {
    appState.pathCollectionModel.remove(index, 1)
    syncPathCollection()
}

function movePathCollectionItem(from, to) {
    appState.pathCollectionModel.move(from, to, 1)
    syncPathCollection()
}
