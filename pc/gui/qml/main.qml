import QtQuick 2.1
import QtQuick.Window 2.0
import "js/util.js" as Util

Window {
    visible: true
    width: 360
    height: 360

    MouseArea {
        anchors.fill: parent
        onClicked: {
            Qt.quit()
        }
    }

    Text {
        id: text
        text: "hello pan-light"
        anchors.centerIn: parent
    }

    Component.onCompleted: {
        // 初始化js工具
        Util.callGoAsync('wait', {
                             "time": 3
                         }).then(function (s) {
                             text.text = s + Util.callGoSync('time')
                         }, null, function (s) {
                             text.text = s
                         })
    }
}
