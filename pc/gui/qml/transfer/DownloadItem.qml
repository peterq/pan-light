import QtQuick 2.0
import "../pan"
import "../comps"
import "../js/app.js" as App
import "../js/util.js" as Util

Item {
    id: root
    height: topRow.height + bottomRow.height
    width: parent.width
    property bool isFinish
    property string downloadId: meta.downloadId
    property var meta
    property int idx
    property string resumeData: ''
    property bool isNewAdd: true
    property string downloadState: ''
    property int progress: 0
    property string errString: ''

    signal taskEvent(string event, var data)

    property string speed: ''

    DataSaver {
        $key: 'download-item-' + root.downloadId
        property alias resumeData: root.resumeData
        property alias isNewAdd: root.isNewAdd
        Component.onCompleted: {
            root.dataSaverOk()
        }
    }

    function dataSaverOk() {
        if (isFinish) return
        if (!isNewAdd) {
            console.log('恢复任务', downloadId)
            var res = Util.callGoSync('download.resume', {
                                          "downloadId": downloadId,
                                          "bin": resumeData
                                      })
        } else {
            isNewAdd = false
            Util.callGoSync('download.start', {
                                "downloadId": downloadId
                            })
        }
        updateState()
    }

    function updateState() {
        downloadState = ''
        var data = Util.callGoSync('download.state', {
                                       "downloadId": downloadId
                                   })
        downloadState = data.state
        progress = data.progress
    }

    Timer {
        id: speedClearTimer
        interval: 2000
        onTriggered: {
            speed = ''
        }
    }

    onTaskEvent: {
        // 更新下载速度
        if (event === 'task.speed') {
            speed = Util.humanSize(data.speed) + '/s'
            speedClearTimer.restart()
            progress = data.progress
            return
        }

        // 下载快照
        if (event === 'task.capture') {
            resumeData = data
            return
        }

        // 下载状态
        if (event === 'task.state') {
            downloadState = data.state
            progress = data.progress
            if (downloadState === 'errored')
                errString = data.error
            if (downloadState === 'completed')
                App.appState.transferComp.itemCompleted(idx)
            return
        }
    }

    function getMenus() {
        root.meta.path = ''
        return []
    }

    Rectangle {
        height: parent.height
        width: Math.min(progress / meta.size, 1) * parent.width
        color: Qt.rgba(140 / 255, 197 / 255, 1, .4)
        visible: !isFinish
    }

    Item {
        id: topRow
        width: parent.width - 20
        height: 50
        x: 10
        FileIcon {
            id: fileIcon
            width: parent.height
            height: width
            anchors.verticalCenter: parent.verticalCenter
            type: root.meta.saveName.split('.').pop()
        }
        Text {
            id: fileNameText
            text: root.meta.saveName
            anchors.verticalCenter: parent.verticalCenter
            anchors.left: fileIcon.right
            anchors.leftMargin: 5
            width: parent.width * 0.4
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

        Row {
            anchors.verticalCenter: parent.verticalCenter
            anchors.right: parent.right
            spacing: 5
            Text {
                text: downloadState
            }
            IconButton {
                iconType: 'start'
                title: '开始'
                color: enabled ? '#409EFF' : 'gray'
                lighter: 1.1
                visible: !isFinish
                enabled: downloadState === 'wait.start'
                onClicked: {
                    Util.callGoSync('download.start', {
                                        "downloadId": downloadId
                                    })
                    updateState()
                }
            }
            IconButton {
                iconType: 'pause'
                title: '暂停'
                color: enabled ? '#409EFF' : 'gray'
                lighter: 1.1
                visible: !isFinish
                enabled: downloadState === 'downloading'
                onClicked: {
                    Util.callGoSync('download.pause', {
                                        "downloadId": downloadId
                                    })
                    updateState()
                }
            }
            IconButton {
                iconType: 'delete'
                title: '删除'
                onClicked: {
                    Util.callGoSync('download.delete', {
                                        "downloadId": downloadId
                                    })
                    App.appState.transferComp.deleteItem(idx, isFinish)
                }
            }
        }
    }

    Item {
        id: bottomRow
        width: parent.width - 20
        x: 10
        visible: !isFinish
        height: visible ? 30 : 0
        y: parent.height / 2

        Text {
            id: progressText
            width: 300
            x: 50
            anchors.verticalCenter: parent.verticalCenter
            text: (progress == 0 ? '' : Util.humanSize(
                                       progress) + '/') + Util.humanSize(
                      meta.size)
        }

        Text {
            id: speedText
            anchors.verticalCenter: parent.verticalCenter
            anchors.left: progressText.right
            text: speed
            color: 'orange'
            width: 300
        }

        Text {
            anchors.left: speedText.right
            visible: downloadState === 'errored'
            text: errString
            color: 'red'
        }
    }

    Rectangle {
        width: parent.width
        height: 1
        color: Qt.lighter('gray')
        anchors.bottom: parent.bottom
    }
}
