import QtQuick 2.0
import QtQuick.Dialogs 1.2
Dialog {
    visible: true
    title: "Blue sky dialog"
    modality: Qt.ApplicationModal
    id: dia

    contentItem: Rectangle {
        color: "lightskyblue"
        implicitWidth: 400
        implicitHeight: 100
        Text {
            text: "Hello blue sky!"
            color: "navy"
            anchors.centerIn: parent
        }
    }
    Timer  {
        interval: 1000
        running: true
        onTriggered: {

        }
    }
}

