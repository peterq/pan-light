import QtQuick 2.0
import QtGraphicalEffects 1.0
Item {
    id: root
    property alias type: img.type
    width: 50
    height: width
    property var color
    Image {
        id: img
        property string type: 'error'
        anchors.fill: parent
        fillMode: Image.PreserveAspectFit
        source: '../assets/images/icons/'+ type +'.svg'
        visible: !parent.color
    }

    ColorOverlay {
        anchors.fill: img
        source: img
        color: root.color || ''
        visible: !!parent.color
    }


    Item {
        id: name
        transform: [
            Rotation {
                id: rotationAni
                origin.x: width / 2
                origin.y: height / 2
                angle: 45
            }

        ]
        PropertyAnimation {
            id: rotationAniCtrl
            target: rotationAni
            property: 'angle'
            running: false
            from: -35
            to: 360 + from
            duration: 2000
            loops: Animation.Infinite
            easing.type: Easing.InOutQuint
        }
    }

    Component.onCompleted: {
        if (type == 'loading') {
            transform.push(rotationAni)
            rotationAniCtrl.running = true
        }
    }
}
