import QtQuick.Dialogs 1.1
import QtQuick 2.0
import QtQuick.Window 2.2
import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3
Window {
    id: root
    //提示框内容
    property alias tipText: msg.text
    //Dialog header、contentItem、footer之间的间隔默认是12
    // 提示框的最小宽度是 100
    property real maxWidth: 800
    property var closeCb: function(){console.log('close alert')}
    width: {
//        console.log('w', msg.implicitWidth)
        if(msg.implicitWidth < 100 || msg.implicitWidth == 100)
            return 100;
        else
            return msg.implicitWidth > maxWidth ? maxWidth + 24 : (msg.implicitWidth + 24);
    }
    height: msg.implicitHeight + 24 + 100

    flags: Qt.Dialog | Qt.WindowModal | Qt.WindowCloseButtonHint
    modality: Qt.WindowModal

    Dialog {
        id: dialog
        width: root.width
        height: root.height
        header: Rectangle {
            width: dialog.width
            height: 50
            radius: 5
            IconFont {
                width: 50
                height: 50
                anchors.centerIn: parent
                type: 'error'
            }
        }
        contentItem: Rectangle {
            Text {
                anchors.fill: parent
                anchors.centerIn: parent
                color: "gray"
                text: tipText
                wrapMode: Text.WordWrap
                verticalAlignment: Text.AlignVCenter
                horizontalAlignment: Text.AlignHCenter
                onLinkActivated: Qt.openUrlExternally(link)
            }
        }
        footer: Rectangle {
            width: msg.width
            height: 50
            radius: 5
            Button {
                anchors.centerIn: parent
                width: 80
                height: 30
                text: '确定'
                onClicked: {
                    root.userClose()
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
        wrapMode: Text.WordWrap
        verticalAlignment: Text.AlignVCenter
        horizontalAlignment: Text.AlignHCenter
    }

    function userClose(){
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

