import QtQuick 2.0
import QtQuick.Window 2.2
import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3
import "../js/app.js" as App
import "../js/util.js" as Util
import "../widget"
Window {
    id: window
    flags: Qt.MSWindowsFixedSizeDialogHint | Qt.WindowTitleHint | Qt.WindowCloseButtonHint
           | Qt.WindowModal | Qt.Dialog
    modality: Qt.ApplicationModal
    title: '反馈'
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

    TopIndicator {
        id: topIndicator
        z: 2
    }

    Column {
        y: 10
        width: parent.width * 0.8
        anchors.horizontalCenter: parent.horizontalCenter
        spacing: 10
        TextArea {
            id: input
            width: parent.width
            height: 200
            placeholderText: '请输入反馈内容'
            focus: true
            background: Rectangle {
               color: '#eee'
            }
        }
        Button {
            text: '提交'
            onClicked: {
                var content = input.text.trim()
                if (!content) return
                enabled = false
                text = '提交中'
                Util.api('feedback', {
                         content: content
                         })
                .then(function() {
                    topIndicator.success('提交成功')
                    return Util.sleep(1000)
                })
                .then(function () {
                    window.visible = false
                })
                .catch(function(error) {
                    topIndicator.fail(error.message)
                    text = '提交'
                    enabled = true
                })
            }

        }
    }
}
