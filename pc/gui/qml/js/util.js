.pragma library

.import "./global.js" as G
.import "./promise.js" as P
.import QtQuick 2.0 as Q

console.log('--------------util js init------------')

var Promise = P.Promise
var setTimeout = G.setTimeout
var fileSep = Qt.platform === 'windows' ? '\\' : '/'

var dialog = (function(){
    var comp = Qt.createComponent('../comps/dialog.qml')
    return function() {
        comp.createObject(G.root)
    }
})()

function sleep(t) {
    return new Promise(function(resolve) {
        setTimeout(resolve, t)
    })
}

var event = {}
;(function () {
    var map = {}
    event.fire = function(evt, data) {
        if (!map[evt]) return
        var fns = []
        map[evt].forEach(function(fn, idx) {
            try {
                fn(data)
                if (!fn.once) {
                    fns.push(fn)
                }
            } catch (e) {
               console.error(evt, fn, e)
            }
        })
        map[evt] = fns
    }

    event.on = function(evt, fn) {
        map[evt] = map[evt] || []
        map[evt].push(fn)
        return function () {
            var idx = map[evt].findIndex(function(v) {
                return v === fn
            })
            if (idx >= 0) {
                map[evt].splice(idx, 1)
            }
        }
    }

    event.once = function(evt, fn) {
       fn.once = 1
       return event.on(evt, fn)
    }

})()

var openDesktopWidget
var hideDesktopWidget
;(function(){
    var comp = loadComponent(function () {}, "../comps/desktop-widget.qml")
    var ins
    openDesktopWidget = function (){
       if (!ins) {
           ins = comp.createObject(null)
           return
       }
       ins.visible = true
    }
    hideDesktopWidget = function () {
        if (ins) {
            ins.visible = false
        }
    }

})()

var bridge = (function () {
    var comp = loadComponent(function () {}, "../comps/bridge.qml")
    var ins = comp.createObject(G.root)
    ins.goMessage.connect(function(data){
        var obj = JSON.parse(data)
        event.fire('go.' + obj.event, obj)
    })
    event.on('go.fuck', function(data) {
        console.log('fuck you from golang, ', data.name)
    })
    return ins
})()

function loadComponent(cb, url) {
    var comp = Qt.createComponent(url)
    function finishCreation() {
               if (comp.status === Q.Component.Ready) {
                   cb(comp)
               } else if (comp.status === Q.Component.Error) {
                   // Error Handling
                   console.log("Error loading component:", comp.errorString());
               }
           }
   if (comp.status === Q.Component.Ready) {
       cb(comp)
   } else {
       finishCreation()
       comp.statusChanged.connect(finishCreation)
   }
   return comp
}

var callGoAsync
var callGoSync
;(function() {
    var promiseMap = {}
    callGoAsync = function(action, param, chan) {
        param = param || {}
        var callId = action + (+new Date) + ~~(Math.random() * 1e5)
        var promise = new P.Promise(function (resolve, reject, progress) {
            var handler = {
                callId: callId,
                resolve: resolve,
                reject: reject,
                progress: progress
            }
            promiseMap[callId] = handler
            bridge.callAsync(JSON.stringify({chan: !!chan, action: action, param: param, callId: callId}))
        })
        promise.callId = callId
        return promise
    }
    event.on('go.call.ret', function(data) {
        var callId = data.callId
        var handler = promiseMap[callId]
        if (!handler) {
            console.error('can not find handler for call id: ', callId)
            return
        }
        handler[data.type](data[data.type])
        if (data.type !== 'progress') {
            delete promiseMap[callId]
        }
    })
    callGoSync = function(action, param) {
        param = param || {}
        var str = bridge.callSync(JSON.stringify({action: action, param: param}))
        return JSON.parse(str).result
    }
})()

function storageGet(key, $default) {
    var ret = callGoSync('storage.get', { k: key })
    if (ret === '')
        return $default
    return JSON.parse(ret)
}

function storageSet(key, value) {
    return callGoSync('storage.set', { k: key, v: JSON.stringify(value || '') })
}

function notifyPromise(promise, data){
    function getCallId(promise) {
        if (promise.callId)
            return promise.callId
        if (promise.parent)
            return getCallId(promise.parent)
        return null
    }
    var asyncCallId = getCallId(promise)
    if (!asyncCallId)
        throw new Error('call id not exist')
    return callGoSync('asyncTaskMsg', {msg: data, asyncCallId: asyncCallId})
}

var alert
;(function () {
    var comp = loadComponent(function () {}, "../comps/Alert.qml")
    alert = function(option){
        option = option || {}
        var defaultOption = {
            parent: G.root,
            title: '请注意',
            msg: '这是一条消息',
            copyButton: false
        }
        for(var k in defaultOption) {
            if (!option.hasOwnProperty(k))
                option[k] = defaultOption[k]
        }
        return new Promise(function(resolve, reject){
            var ins = comp.createObject(option.parent, {
                                            tipText: option.msg,
                                            title: option.title,
                                            closeCb: resolve,
                                            copyButton: option.copyButton
                                        })
        })
    }
})()

var prompt
;(function () {
    var comp = loadComponent(function () {}, "../comps/prompt-window.qml")
    prompt = function(option){
        option = option || {}
        var defaultOption = {
            parent: G.root,
            title: '请输入',
            msg: '这是一条消息',
            checkFunc: function() {return true},
            content: ''
        }
        for(var k in defaultOption) {
            if (!option.hasOwnProperty(k))
                option[k] = defaultOption[k]
        }
        return new Promise(function(resolve, reject){
            function onClose(result) {
                if (result === false) reject('input canceled')
                resolve(result)
            }
            var ins = comp.createObject(option.parent, {
                                            tipText: option.msg,
                                            title: option.title,
                                            closeCb: onClose,
                                            checkFunc: option.checkFunc,
                                            content: option.content
                                        })
        })
    }
})()

var confirm
;(function () {
    var comp = loadComponent(function () {}, "../comps/confirm-window.qml")
    confirm = function(option){
        option = option || {}
        var defaultOption = {
            parent: G.root,
            title: '是否继续',
            msg: '这是一条消息',
        }
        for(var k in defaultOption) {
            if (!option.hasOwnProperty(k))
                option[k] = defaultOption[k]
        }
        return new Promise(function(resolve, reject){
            function onClose(result) {
                if (result === false) reject(result)
                resolve(result)
            }
            var ins = comp.createObject(option.parent, {
                                            tipText: option.msg,
                                            title: option.title,
                                            closeCb: onClose
                                        })
        })
    }
})()

var pickSavePath
;(function () {
    var comp = loadComponent(function () {}, "../comps/select-save-path.qml")
    pickSavePath = function(option){
        option = option || {}
        var defaultOption = {
            parent: G.root,
            title: '选择保存位置',
            defaultFolder: '',
            fileName: '',
        }
        for(var k in defaultOption) {
            if (!option.hasOwnProperty(k))
                option[k] = defaultOption[k]
        }
        return new Promise(function(resolve, reject){
            function onClose(result) {
                if (result === false) reject(result)
                resolve(result)
            }
            var ins = comp.createObject(option.parent, {
                                            defaultFolder: option.defaultFolder,
                                            title: option.title,
                                            fileName: option.fileName,
                                            resolve: resolve,
                                            reject: reject
                                        })
            ins.open()
        })
    }
})()

var showMenu = function(items, parent ) {
    G.appState.globalRightMenu.show(items)
}

function isVideo(f){
    var ext = ['mp4', 'avi', 'rmvb', 'mkv', 'mov', 'wmv']
    var e = f.split('.').pop().toLowerCase()
    return ext.indexOf(e) >= 0
}

var videoAgentLink = (function(){
    var server
    return function(fid) {
        if (!server) {
            server =  callGoSync('env.internal_server_url')
            console.log('internal url', JSON.stringify(server))
        }
        return server + '/videoAgent?fid=' + fid
    }
})()

function getFileLink(meta) {
    return callGoAsync('pan.link', {fid: 'direct.' + meta.fs_id})
}

function getFileLinkVip(meta) {
    return callGoAsync('pan.link', {fid: 'vip.' + meta.fs_id})
}

var playVideoByLink
var playVideo = (function(){
    var comp = loadComponent(function(){},'../videoPlayer/MPlayer.qml')
    var ins
    playVideoByLink = function(name, link) {
        if (!ins || !ins.playVideo) {
            ins = comp.createObject(G.root)
        }
        ins.playVideo(name, link)
    }
    return function(meta, useVip){
        if (!ins || !ins.playVideo) {
            ins = comp.createObject(G.root)
        }
        var fid = (useVip ? 'vip.' : 'direct.') + meta.fs_id
        var linkPromise = callGoAsync('pan.link', {fid: fid})
        .then(function(link){
            var agentLink = videoAgentLink(fid)
            console.log('play link', agentLink, link)
            return agentLink
        })
        linkPromise.loadingLinkText = '正在解析播放链接'
        ins.playVideo(meta.server_filename, linkPromise)
    }
})()


var tooTip = (function(){
    var comp = loadComponent(function(){},'../comps/tool-tip-window.qml')
    var ins
    return function(){
        if (!ins) {
            ins = comp.createObject(G.root)
        }
        return ins

    }
})()

function listModelToArr(model) {
    var arr = []
    for (var i = 0; i < model.count; i++) {
        arr.push(model.get(i))
    }
    return arr
}

function arrToListModel(arr, model) {
    model.clear()
    arr.forEach(function (item) {
        model.append(item)
    })
    return model
}

function listModelAdd(model, data) {
    model.append(data)
    return listModelToArr(model)
}

function listModelClear(model) {
    model.clear()
    return []
}

function listModelRemove(model, index) {
    model.remove(index, 1)
    return listModelToArr(model)
}

function listModelMove(model, from, to) {
    model.move(from, to, 1)
    return listModelToArr(model)
}

function humanSize(size) {
    var i, unit = ['B', 'KB', 'MB', 'GB']
    for (i = 0; i < unit.length - 1; i++) {
        if (size < 1024)
            break
        size /= 1024
    }
    return size.toFixed(2) + unit[i]
}

var openSetting = (function () {
    var comp = loadComponent(function(){},'../pages/setting-window.qml')
    var ins
    return function(){
        if (!ins || !ins.visible) {
            ins = comp.createObject(G.root)
        }
        return ins
    }
})()

var openAbout = (function () {
    var comp = loadComponent(function(){},'../pages/about-window.qml')
    var ins
    return function(){
        if (!ins || !ins.visible) {
            ins = comp.createObject(G.root)
        }
        return ins
    }
})()

var openFeedback = (function () {
    var comp = loadComponent(function(){},'../pages/feedback-window.qml')
    var ins
    return function(){
        if (!ins || !ins.visible) {
            ins = comp.createObject(G.root)
        }
        return ins
    }
})()

var openShare = (function () {
    var comp = loadComponent(function(){},'../pages/share-window.qml')
    var ins
    return function(meta){
        if (!ins || !ins.visible) {
            ins = comp.createObject(G.root, {meta: meta})
        }
        return ins
    }
})()

function api(name, param) {
    param = param || {}
    return callGoAsync('api.call', {name: name, param: param})
}

function digital(i) {
  return i < 10 ? '0' + i: i
}

function unixTime(t) {
    var d = new Date(t * 1000)
    var date = [d.getFullYear(),d.getMonth()+1, d.getDate()].map(digital).join('-')
    var time = [d.getHours(), d.getMinutes(), d.getSeconds()].map(digital).join(':')
    return date + ' ' + time
}

function exit() {
    callGoSync('exit')
}

