import QtQuick 2.0
import QtGraphicalEffects 1.0
import "../../js/app.js" as App

Item {
    property var player: App.appState.player
    width: parent.width
    height: parent.height
    LinearGradient {
        width: parent.width
        height: parent.height * 1.2
        gradient: Gradient {
            GradientStop {
                position: 0.0
                color: Qt.rgba(0, 0, 0, 0.6)
            }
            GradientStop {
                position: 0.4
                color: Qt.rgba(0, 0, 0, 0.6)
            }
            GradientStop {
                position: 1.0
                color: Qt.rgba(0, 0, 0, 0)
            }
        }
        start: Qt.point(0, 0)
        end: Qt.point(0, height)
    }

    Text {
        id: title
        text: player.videoTitle
        color: 'white'
        font.pointSize: 20
        height: parent.height - 10
        width: parent.width - 30
        anchors.centerIn: parent
        elide: Text.ElideRight
    }

}
