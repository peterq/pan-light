import QtQuick 2.11
import "../js/app.js" as App

Rectangle {
    id: root
    color: Qt.rgba(1, 1, 1, .6)
    anchors.fill: parent
    property string error: '错误'
    visible: false

    Connections {
        target: App.appState
        onEnterPathPromiseChanged: {
            var p = App.appState.enterPathPromise
            if (!p) return
            root.visible = false
            p.catch(function(e) {
                root.error = e
                root.visible = true
            })
        }
    }

    Column {
        anchors.centerIn: parent
        Text {
            id: title
            text: '出错了'
            font.pointSize: 15
        }

        Text {
            text: '尝试加载目录: <span style="color:orange">' + App.appState.path + '</span> 时出错. 检查一下网络 ? '
                  + '<br><a href="reload">重新加载</a>&nbsp;&nbsp;&nbsp;&nbsp;<a href="detail" style="margin-left:10px">错误详情</a>'
            textFormat: Text.RichText
            width: Math.min(400, root.width / 2)
            wrapMode: Text.Wrap
            onLinkActivated: {
                root[link]()
            }
        }
    }
    function reload() {
        App.enterPath(App.appState.path)
    }
    function detail() {
        App.alert('错误详情', error, true)
    }
}
