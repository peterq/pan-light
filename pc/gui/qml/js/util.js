.pragma library
.import QtQuick 2.0 as Q
.import "global.js" as G
.import "promise.js" as P
var event = {}
;(function () {
    var map = {}
    event.fire = function (evt, data) {
        if (!map[evt]) return
        map[evt].forEach(function (fn, idx) {
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

    event.on = function (evt, fn) {
        map[evt] = map[evt] || []
        map[evt].push(fn)
        return function () {
            var idx = map[evt].findIndex(function (v) {
                return v === fn
            })
            if (idx >= 0) {
                map[evt].splice(idx, 1)
            }
        }
    }

    event.once = function (evt, fn) {
        fn.once = 1
        return event.on(evt, fn)
    }

})()

var bridge = (function () {
    var comp = loadComponent(function () {
    }, "../comps/bridge.qml")
    var ins = comp.createObject(G.root)
    ins.goMessage.connect(function (data) {
        var obj = JSON.parse(data)
        event.fire('go.' + obj.event, obj)
    })
    event.on('go.fuck', function (data) {
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
            console.log("Error loading component:", comp.errorString())
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
;(function () {
    var promiseMap = {}
    callGoAsync = function (action, param, chan) {
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
    event.on('go.call.ret', function (data) {
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
    callGoSync = function (action, param) {
        param = param || {}
        var str = bridge.callSync(JSON.stringify({action: action, param: param}))
        return JSON.parse(str).result
    }
})()
