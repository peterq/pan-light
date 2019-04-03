import QtQuick 2.0

import QtQuick.Controls 2.1
import "../../comps"
import "../../widget"
import "../../js/util.js" as Util

// 网盘用量
Column {
    id: usage
    width: parent.width
    property real total: 0
    property real used: 0
    property bool loading: false
    Row {
        spacing: 5
        Item {
            height: 1
            width: 10
        }
        Text {
            text: '网盘空间: ' + (loading ? '加载中...' : usage.humanSize(
                                           usage.used) + '/' + usage.humanSize(
                                           usage.total))
        }
        Button {
            id: iconBtn
            width: 25
            height: width
            anchors.verticalCenter: parent.verticalCenter
            ToolTip {
                show: iconBtn.hovered
                text: '刷新'
            }
            IconFont {
                id: icon
                type: 'refresh'
                width: parent.width
                color: iconBtn.hovered ? 'red' : Qt.lighter('red')
            }
            display: AbstractButton.IconOnly
            background: Item {
            }
            transform: [
                Rotation {
                    id: rotationAni
                    origin.x: iconBtn.width / 2
                    origin.y: iconBtn.height / 2
                }
            ]

            PropertyAnimation {
                target: rotationAni
                property: 'angle'
                running: loading
                from: 0
                to: -360
                duration: 1000
                loops: Animation.Infinite
                easing.type: Easing.Linear
            }
            onClicked: {
                if (loading)
                    return
                getUsage()
            }
        }
    }
    function humanSize(size) {
        var i, unit = ['B', 'KB', 'MB', 'GB']
        for (i = 0; i < unit.length - 1; i++) {
            if (size < 1024)
                break
            size /= 1024
        }
        return '' + Math.round(size) + unit[i]
    }
    function getUsage() {
        loading = true
        return Util.callGoAsync('pan.usage').then(function (data) {
            total = data.total
            used = data.used
        }).finally(function () {
            loading = false
        })
    }
    Component.onCompleted: {
        getUsage()
    }
}
