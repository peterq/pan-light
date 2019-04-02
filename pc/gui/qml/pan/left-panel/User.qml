import QtQuick 2.0
import QtGraphicalEffects 1.0
import QtQuick.Controls 2.1
import "../../js/app.js" as App
// 头像, 用户名
Row {
    id: user
    width: parent.width
    height: 60
    spacing: 5
    Item {
        width: 10
        height: 1
    }
    Rectangle {
        width: 50
        height: width
        radius: width / 2
        anchors.verticalCenter: parent.verticalCenter
        color: 'orange'
        Image {
            id: avatar
            smooth: true
            visible: false
            anchors.fill: parent
            source: App.appState.loginSession.photo
            antialiasing: true
        }
        Rectangle {
            id: mask
            anchors.fill: parent
            radius: width / 2
            visible: false
            antialiasing: true
            smooth: true
        }
        OpacityMask {
            anchors.fill: avatar
            source: avatar
            maskSource: mask
            antialiasing: true
            visible: avatar.status === Image.Ready
        }
    }

    Text {
        anchors.verticalCenter: parent.verticalCenter
        text: App.appState.loginSession.username
    }
}
