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
    property real minWidth: 200
    property real maxWidth: 800
    property bool copyButton: false
    property var closeCb: function () {
        console.log('close alert')
    }
    width: {
        if (msg.implicitWidth <= minWidth)
            return minWidth
        return Math.min(msg.implicitWidth, maxWidth) + 24
    }
    height: msg.implicitHeight + 24 + 100

    flags: Qt.MSWindowsFixedSizeDialogHint | Qt.WindowTitleHint | Qt.WindowCloseButtonHint
           | Qt.WindowModal | Qt.Dialog
    modality: Qt.WindowModal

    Dialog {
        id: dialog
        width: root.width
        height: root.height
        header: Rectangle {
            width: dialog.width
            height: 50
            IconFont {
                width: 50
                height: 50
                anchors.centerIn: parent
                type: 'error'
            }
        }
        contentItem: Rectangle {
            anchors.fill: parent
            TextEdit {
                id: textContent
                width: parent.width
                anchors.centerIn: parent
                color: "gray"
                text: tipText
                wrapMode: Text.Wrap
                verticalAlignment: Text.AlignVCenter
                horizontalAlignment: Text.AlignHCenter
                onLinkActivated: Qt.openUrlExternally(link)
                selectByMouse: true
                readOnly: true
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
                    text: '复制'
                    visible: copyButton
                    onClicked: {
                        textContent.selectAll()
                        textContent.copy()
                        text = '已复制'
                    }
                }
                Button {
                    width: 80
                    height: 30
                    text: '确定'
                    onClicked: {
                        root.userClose()
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

    function userClose() {
        root.destroy()
        root.closeCb()
    }

    onClosing: {
        userClose()
    }

    Component.onCompleted: {
        dialog.open()
        visible = true
        requestActivate()
    }
}
