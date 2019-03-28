import QtQuick 2.9
import QtQuick.Window 2.2
import "js/util.js" as Util
import "js/global.js" as G
import "./layout"
import "./login"
import "./videoPlayer"
Window {
    id: mainWindow
    visible: true
    width: 1000
    height: 680
    minimumHeight: 600
    minimumWidth: 900
    title: "hello peterq2"
    // 用来触发窗口重汇
    Rectangle {
       id: re
       width: 0
       height: 0
       z: -1
    }
    Layout{}
    Login{}
    Component.onCompleted: {

        // 初始化js工具
        G.init(mainWindow)
        // Util.openDesktopWidget()
        function getSign() {
            Util.callGoAsync('pan.init')
                .then(function(data){
                    console.log('api init success')
                    Util.event.fire('init.api.ok', data)
                })
                .catch(function(e){
                    console.log('get sign error', e)
                    Util.event.fire('init.not-login')
                    Util.event.once('login.success', function(){
                        console.log('get sign again')
                        getSign()
                    })
                })
        }
        getSign()
    }
    onWidthChanged: {
        Util.setTimeout(function () {
          re.width = mainWindow.width
        }, 1)
    }
    onHeightChanged: {
        Util.setTimeout(function () {
          re.height = mainWindow.height
        }, 1)
    }
}
