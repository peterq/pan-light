import QtQuick 2.0
import QtQuick.Controls 2.3
import "../../js/global.js" as G

Button {
    id: iconBtn
    icon.color: 'white'
    icon.width: width / 1.5
    icon.height: height / 1.5
    display: AbstractButton.IconOnly
    opacity: 0
    scale: 1
    enabled: false
    background: Rectangle {
        color: Qt.rgba(0, 0, 0, .8)
        radius: parent.width / 2
    }

    function ani() {
        iconBtn.opacity = 0.8
        showIconAnimation.restart()
    }

    ParallelAnimation {
        id: showIconAnimation
        SequentialAnimation {
            PauseAnimation {
                duration: 500
            }
            PropertyAnimation {
                id: opacityAnimation
                target: iconBtn
                property: 'opacity'
                from: 0.8
                to: 0
                duration: 1000
            }
        }
        PropertyAnimation {
            id: scaleAnimation
            target: iconBtn
            property: 'scale'
            from: 1
            to: 2.5
            duration: 1800
        }

    }
}
