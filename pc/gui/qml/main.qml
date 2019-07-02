import QtQuick 2.9
import QtQuick.Window 2.2
import Qt.labs.platform 1.0
import "js/util.js" as Util
import "js/global.js" as G
import "js/app.js" as App
import "./layout"
import "./login"
import "./videoPlayer"
import "./widget"
import "./comps"
Window {
    id: mainWindow
    visible: true
    width: 1000
    height: 680
    minimumHeight: 600
    minimumWidth: 900
    title: "pan-light"
    signal customerEvent(string event, var data)
    flags: Qt.FramelessWindowHint
    color: 'transparent'
    visibility: Window.Windowed

    DataSaver {
        $key: 'window.main'
        property alias x: mainWindow.x
        property alias y: mainWindow.y
        property alias width: mainWindow.width
        property alias height: mainWindow.height
        property string firstStart: '1'
        Component.onCompleted: {
            if (firstStart === '1') {
                firstStart = ''
                mainWindow.x = (Screen.desktopAvailableWidth - mainWindow.width) / 2
                mainWindow.y = (Screen.desktopAvailableHeight - mainWindow.height) / 2
            }
        }
    }

    Component {
        id: layoutComp
        Layout {}
    }
    VirtualFrame {
        x: 0
        y: 0
        content: Item {
            anchors.fill: parent
            Text {
                text: 'pan-light 初始化中...'
                font.pointSize: 20
                anchors.centerIn: parent
            }
            Loader {
                anchors.fill: parent
                id: layoutLoader
                focus: true
            }

            Login{}
            Component.onCompleted: {
                // 初始化js工具
                G.init(mainWindow)
                App.appState.mainWindow = mainWindow
                function getSign() {
                    Util.callGoAsync('pan.init')
                        .then(function(data){
                            console.log('api init success')
                            layoutLoader.sourceComponent = layoutComp
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
        }

    }
}
