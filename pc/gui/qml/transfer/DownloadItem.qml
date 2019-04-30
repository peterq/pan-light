import QtQuick 2.0
import "../pan"
import "../comps"
import "../js/app.js" as App
import "../js/util.js" as Util

Item {
    id: root
    height: 50
    width: parent.width
    property bool isFinish
    property string downloadId: meta.downloadId
    property var meta

    property string resumeData: ''

    property bool isNewAdd: true

    Component.onCompleted: {
        if (!isNewAdd) {
            console.log('恢复任务', downloadId)
            Util.callGoSync('download.resume', {
                                "downloadId": downloadId,
                                "bin": resumeData
                            })
        }
    }

    DataSaver {
        $key: 'download-item-' + root.downloadId
        property alias resumeData: root.resumeData
        property alias isNewAdd: root.isNewAdd
    }

    function getMenus() {
        root.meta.path = ''
        return []
    }
    FileIcon {
        id: fileIcon
        width: parent.height
        height: width
        anchors.verticalCenter: parent.verticalCenter
        anchors.leftMargin: 10
        type: root.meta.saveName.split('.').pop()
    }
    Text {
        text: root.meta.saveName
        anchors.verticalCenter: parent.verticalCenter
        anchors.left: fileIcon.right
        anchors.leftMargin: 5
        anchors.right: parent.right
        anchors.rightMargin: 5
        elide: Text.ElideRight
    }

    MouseArea {
        hoverEnabled: true
        anchors.fill: parent
        acceptedButtons: Qt.LeftButton | Qt.RightButton
        onClicked: {
            if (mouse.button === Qt.RightButton)
                Util.showMenu(getMenus())
        }
    }
    Rectangle {
        width: parent.width
        height: 1
        color: Qt.lighter('gray')
        anchors.bottom: parent.bottom
    }
}
