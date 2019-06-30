import QtQuick.Dialogs 1.1
import QtQuick 2.11
import QtQuick.Window 2.2
import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3

Window {
    id: root
    //提示框内容
    property alias tipText: msg.text
    //Dialog header、contentItem、footer之间的间隔默认是12
    // 提示框的最小宽度是 100
    property real minWidth: 300
    property real maxWidth: 800
    property var closeCb: function () {
        console.log('close confirm')
    }
    width: {
        if (msg.implicitWidth <= minWidth)
            return minWidth
        return Math.min(msg.implicitWidth, maxWidth) + 24
    }
    height: msg.implicitHeight + 24 + 150

    flags: Qt.MSWindowsFixedSizeDialogHint | Qt.WindowTitleHint | Qt.WindowCloseButtonHint
            | Qt.WindowModal | Qt.Dialog
    modality: Qt.WindowModal

    Dialog {
        id: dialog
        width: root.width
        height: root.height
        focus: true
        header: Rectangle {
            width: parent.height
            height: textContent.implicitHeight + 12
            Text {
                id: textContent
                width: parent.width - 12
                anchors.centerIn: parent
                color: "gray"
                text: tipText
                wrapMode: Text.Wrap
                verticalAlignment: Text.AlignVCenter
                horizontalAlignment: Text.AlignHCenter
            }
        }
        contentItem: Rectangle {
            width: dialog.width
            height: 100
            Text {
                id: errMsg
                height: 30
                text: typeof checkResult === 'string' ? checkResult : ''
                color: 'red'
            }
        }
        footer: Item {
            width: parent.width
            height: 50
            Row {
                spacing: 10
                height: parent.height
                anchors.centerIn: parent
                Button {
                    width: 80
                    height: 30
                    text: '取消'
                    onClicked: {
                        root.userClose(false)
                    }
                }
                Button {
                    width: 80
                    height: 30
                    text: '确定'
                    onClicked: {
                        root.userClose(true)
                    }
                }
            }
        }
    }

    //利用Text 的implicitWidth属性来调节提示框的大小
    //该Text的字体格式需要与contentItem中的字体一模一样
    Text {
        id: msg
        visible: false
        width: maxWidth
        wrapMode: Text.Wrap
        verticalAlignment: Text.AlignVCenter
        horizontalAlignment: Text.AlignHCenter
    }

    function userClose(result) {
        root.destroy()
        root.closeCb(result)
    }

    onClosing: {
        userClose(false)
    }

    Component.onCompleted: {
        visible = true
        requestActivate()
        dialog.open()
    }
}
