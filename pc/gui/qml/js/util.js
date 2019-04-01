.pragma library

.import "./global.js" as G
.import "./promise.js" as P
.import QtQuick 2.0 as Q

console.log('--------------util js init------------')

var Promise = P.Promise
var setTimeout = G.setTimeout

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
        map[evt].forEach(function(fn, idx) {
            try {
                fn(data)
                if (fn.once) {
                    map[evt].splice(idx, 1)
                }
            } catch (e) {
               console.error(evt, fn, e)
            }
        })
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
       ins.visbal = true
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
            msg: '这是一条消息'
        }
        for(var k in defaultOption) {
            if (!option.hasOwnProperty(k))
                option[k] = defaultOption[k]
        }
        return new Promise(function(resolve, reject){
            var ins = comp.createObject(option.parent, {
                                            tipText: option.msg,
                                            title: option.title,
                                            closeCb: resolve
                                        })
        })
    }
})()

var showMenu = (function () {
    var comp = loadComponent(function () {}, "../comps/rightClickMenu.qml")
    return function(items, parent) {
        var ins = comp.createObject(parent || G.root, {menus: items})
        ins.popup()
        ins.aboutToHide.connect(function(){
            ins.destroy()
        })
        return ins
    }
})()

function isVideo(f){
    var ext = ['mp4', 'avi', 'rmvb', 'mkv', 'mov']
    var e = f.split('.').pop().toLowerCase()
    return ext.indexOf(e) >= 0
}

var videoAgentLink = (function(){
    var server
    return function(meta) {
        if (!server) {
            server =  callGoSync('env.internal_server_url')
            console.log('internal url', JSON.stringify(server))
        }
        return server + '/videoAgent?fid=' + meta.fs_id
    }
})()

function getFileLink(meta) {
    return callGoAsync('pan.link', {fid: meta.fs_id})
}

function getFileLinkVip(meta) {
    return callGoAsync('pan.link.vip', {fid: meta.fs_id})
}

var playVideo = (function(){
    var comp = loadComponent(function(){},'../videoPlayer/MPlayer.qml')
    var ins
    return function(meta, useVip){
        if (!ins || !ins.playVideo) {
            ins = comp.createObject(G.root)
        }
        (useVip ?
             getFileLinkVip(meta) :
             getFileLink(meta))
        .then(function(link){
            var agentLink = videoAgentLink(meta, useVip)
            console.log('play link', agentLink, link)
            ins.playVideo(meta.server_filename, agentLink)
        })

    }
})()


var tooTip = (function(){
    var comp = loadComponent(function(){},'../comps/tool-tip-window.qml')
    var ins
    return function(meta, useVip){
        if (!ins || !ins.playVideo) {
            ins = comp.createObject(G.root)
        }
        return ins

    }
})()
























