import QtQuick 2.0
import QtQuick.Controls 2.4
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
    property var listComp
    property string resumeData: ''
    property bool isNewAdd: true
    property string downloadState: ''
    property real progress: 0
    property string errString: ''
    property bool isQueued: false

    signal taskEvent(string event, var data)
    signal doStart
    signal doPause

    onDoStart: {
        if (!startBtn.enabled)
            return
        isQueued = true
        listComp.enqueue(downloadId)
    }

    onDoPause: {
        if (!pauseBtn.enabled)
            return
        if (isQueued) {
            listComp.dequeue(downloadId)
            isQueued = false
        } else {
            pause()
        }
        listComp.checkQueue()
    }

    property string speed: ''
    property int speedInt: 0

    function pause() {
        Util.callGoSync('download.pause', {
                            "downloadId": downloadId
                        })
        updateState()
    }

    function start() {
        isQueued = false
        Util.callGoSync('download.start', {
                            "downloadId": downloadId
                        })
        updateState()
    }

    Connections {
        target: App.appState.transferComp
        onSumSpeed: {
            if (isFinish)
                return
            data.speed += speedInt
        }
    }

    DataSaver {
        $key: 'download-item-' + root.downloadId
        property alias resumeData: root.resumeData
        property alias isNewAdd: root.isNewAdd
        Component.onCompleted: {
            root.dataSaverOk()
        }
    }

    function dataSaverOk() {
        if (isFinish)
            return
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
        interval: 1100
        onTriggered: {
            speed = ''
            speedInt = 0
        }
    }

    onTaskEvent: {
        // 更新下载速度
        if (event === 'task.speed') {            
            speed = Util.humanSize(data.speed) + '/s'
            speedInt = data.speed
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
        IconFont {
            type: 'vip'
            width: 30
            visible: !isFinish && meta.useVip
            anchors.verticalCenter: parent.verticalCenter
            anchors.left: fileNameText.left
            anchors.leftMargin: 10 + Math.min(fileNameText.width,
                                              fileNameText.implicitWidth)
        }

        MouseArea {
            hoverEnabled: true
            anchors.fill: parent
            acceptedButtons: Qt.LeftButton | Qt.RightButton
            onClicked: {

            }
        }

        PromiseDialog {
            id: comfirmDeleteDialog
            w: 500
            h: 200
            title: '确认删除'
            onClickConfirm: function () {
                result = checkBox.checked
                return true
            }
            contentItem: Column {
                Text {
                    text: '删除后不可恢复! 确认删除?'
                }
                CheckBox {
                    id: checkBox
                    visible: isFinish
                    checked: true
                    text: '同时删除本地文件'
                }
            }
        }

        Row {
            anchors.verticalCenter: parent.verticalCenter
            anchors.right: parent.right
            spacing: 5
            Text {
                text: downloadState
            }
            Text {
                text: '等待中...'
                color: 'orange'
                visible: isQueued
            }
            IconButton {
                id: startBtn
                iconType: 'start'
                title: '开始'
                color: enabled ? '#409EFF' : 'gray'
                lighter: 1.1
                visible: !isFinish
                enabled: !isQueued && downloadState === 'wait.start'
                onClicked: {
                    doStart()
                }
            }
            IconButton {
                id: pauseBtn
                iconType: 'pause'
                title: '暂停'
                color: enabled ? '#409EFF' : 'gray'
                lighter: 1.1
                visible: !isFinish
                enabled: downloadState === 'downloading' || isQueued
                onClicked: {
                    doPause()
                }
            }

            IconButton {
                iconType: 'folder'
                title: '打开所在文件夹'
                onClicked: {
                    var dir = meta.savePath.split(Util.fileSep)
                    dir.pop()

                    dir = dir.join(Util.fileSep)
                    var ret = Qt.openUrlExternally('file://' + dir)
                    console.log(ret)
                }
            }

            IconButton {
                iconType: 'delete'
                title: '删除'
                onClicked: {
                    comfirmDeleteDialog.open().then(function (deleteFile) {
                        Util.callGoSync('download.delete', {
                                            "downloadId": downloadId,
                                            "path": meta.savePath,
                                            "deleteFile": isFinish ? deleteFile : false
                                        })
                        App.appState.transferComp.deleteItem(idx, isFinish)
                    })
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
