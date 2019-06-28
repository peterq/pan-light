import QtQuick 2.0
import QtQuick.Window 2.2
import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3
import "../js/app.js" as App
import "../js/util.js" as Util
Window {
    id: window
    flags: Qt.MSWindowsFixedSizeDialogHint | Qt.WindowTitleHint | Qt.WindowCloseButtonHint
           | Qt.WindowModal | Qt.Dialog
    modality: Qt.ApplicationModal
    title: '设置'
    minimumHeight: height
    minimumWidth: width
    maximumHeight: height
    maximumWidth: width
    visible: true
    width: 550
    height: 300

    Component.onCompleted: {
        visible = true
        requestActivate()
    }

    onVisibleChanged: {
        if (!visible) {
            window.destroy()
        }
    }

    GridLayout {
        columns: 3
        width: parent.width - 20
        y: 20
        anchors.horizontalCenter: parent.horizontalCenter
        Label {
            text: '同时下载任务数'
            width: parent.width * 30
            Layout.alignment: Qt.AlignRight
        }
        TextField {
            id: taskNumberField
            width: parent.width * 30
            validator: IntValidator {
                bottom: 1
                top: 10
            }
            placeholderText: "任务数"
            Component.onCompleted: {
                text = App.appState.settings.maxParallelTaskNumber
            }
            onTextChanged: {
                var n = text * 1
                if (n <= 0) {
                    n = 1
                }
                App.appState.settings.maxParallelTaskNumber = n
            }
        }
        Label {
            text: '1-10之间'
        }
        Label {
            text: '每个任务的线程数'
            width: parent.width * 30
            Layout.alignment: Qt.AlignRight
        }
        TextField {
            width: parent.width * 30
            validator: IntValidator {
                           bottom: 1
                           top: ~~(1024 / App.appState.settings.maxParallelTaskNumber)
                       }
            placeholderText: "并发数1-" + ~~(1024 / App.appState.settings.maxParallelTaskNumber)
            onAccepted: textArea.focus = true
            Component.onCompleted: {
                text = App.appState.settings.maxParallelCorutineNumber
            }
            onTextChanged: {
                var n = text * 1
                if (n <= 0) {
                    n = 1
                }
                App.appState.settings.maxParallelCorutineNumber = n
            }
        }
        Label {
            text: '线程数 x 任务数小于1024'
        }

        Label {
            text: ''
            width: parent.width * 30
            Layout.alignment: Qt.AlignRight
        }
        Button {
            text: '重启以生效'
            onClicked: {
                Util.callGoSync('reboot')
            }
        }
    }
}
