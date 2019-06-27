.pragma library
// 此文件存放全局变量和 Polly fill, 不要依赖其他js, 否则会造成依赖循环

// global
var g = {}
// qml 根组件
var root
function init(r) {
    root = r
}

var appState

// DataSaver 使用过的key
var dataSaverKeys = {}

// internal server(go) url
var internalServerUrl = ""


var setTimeout = (function(){
    var timer = Qt.createComponent("../comps/timer.qml")
    return function (cb, time) {
        cb = cb || function (){}
        time = time || 0
        var ins = timer.createObject(root, {
                               interval: time,
                               cb: cb
                           })
        return function cancel() {
            if (ins) {
                ins.destroy()
            }
        }
    }
})()
