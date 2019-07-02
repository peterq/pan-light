import QtQuick 2.0
import '../videoPlayer'
import '../js/app.js' as App
import '../js/util.js' as Util
import '../comps'
import './left-panel'
Rectangle {
    id: root

    signal active

    // 左侧信息栏
    LeftPanel {
        id: leftPanel
        height: parent.height
        width: 250
        // 右侧border
        Rectangle {
            color: 'gray'
            width: 2
            height: parent.height
            anchors.right: parent.right
        }
    }
    // 头部加列表
    Rectangle {
        anchors.left: leftPanel.right
        width: root.width - leftPanel.width
        height: parent.height
        PathNav {
            id: pathNav
            width: parent.width
            height: 40
            // 下侧border
            Rectangle {
                color: 'gray'
                width: parent.width
                height: 2
                anchors.bottom: parent.bottom
            }
        }
        FileList {
            width: parent.width
            height: parent.height - pathNav.height
            anchors.top: pathNav.bottom
            color: "#fff"

            Rectangle {
                id: loading
                anchors.fill: parent
                visible: false
                color: Qt.rgba(1,1,1,.6)
                Component {
                    id: iconComp
                    IconFont {
                        type: 'loading'
                        width: Math.min(loading.width, loading.height) * 0.3
                    }
                }
                Loader {
                    id: iconLoader
                    focus: true
                    anchors.centerIn: parent
                }
                Component.onCompleted: {
                    // 监听进入path, 延时500ms显示加载动画
                    App.appState.enterPathPromiseChanged.connect(function() {
                        if (!App.appState.enterPathPromise) {
                            loading.visible = false
                            iconLoader.sourceComponent = null
                            return
                        }
                        var p = App.appState.enterPathPromise
                        Util.sleep(500)
                        .then(function(){
                            if (App.appState.enterPathPromise === p) {
                                iconLoader.sourceComponent = iconComp
                                loading.visible = true
                            }
                        })
                    })
                }
            }
            LoadDirError{}
        }
    }
}
