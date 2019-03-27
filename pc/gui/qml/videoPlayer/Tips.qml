import QtQuick 2.0
import './UIComp'
Item {
    anchors.fill: parent

    PlayIcon {
        anchors.centerIn: parent
        width: 40
        height: width
    }

    ForwardBackward {
        anchors.centerIn: parent
        width: 40
        height: width
    }

    VolumeIcon {
        anchors.centerIn: parent
        width: 40
        height: width
    }

    ActionTips {
        anchors.topMargin: 10
        anchors.rightMargin: 10
        anchors.right: parent.right
        anchors.top: parent.top
    }
}
