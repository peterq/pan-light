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
    title: '关于'
    minimumHeight: height
    minimumWidth: width
    maximumHeight: height
    maximumWidth: width
    visible: true
    width: 600
    height: 400

    property string version: 'v0.0.1-preview'
    property string userAgreementLink: 'https://pan-light.peterq.cn/user-agreement'
    property string gitRepoUrl: 'https://github.com/peterq/pan-light'
    property string email: 'me@peterq.cn'

    Component.onCompleted: {
        visible = true
        requestActivate()
        version = Util.callGoSync('env.version')
    }

    onVisibleChanged: {
        if (!visible) {
            window.destroy()
        }
    }

    Item {
        anchors.fill: parent
        Rectangle {
            id: logoCon
            width: parent.width
            height: parent.height * 0.7
            Image {
                source: '../assets/images/pan-light-1.png'
                height: parent.height * 0.8
                width: height
                anchors.centerIn: parent
            }
        }
        Rectangle {
            width: parent.width
            height: parent.height - logoCon.height
            anchors.top: logoCon.bottom
            color: '#eee'
            Text {
                text: ['<b>pan-light&nbsp;&nbsp;&nbsp;' + window.version + '</b>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<a href="'
                    + window.userAgreementLink + '" >查看用户协议</a>',
                    '作&nbsp;&nbsp;&nbsp;&nbsp;者&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;PeterQ&lt;<a href="mailto:' + window.email
                    + '" >' + window.email
                    + '</a>&gt;', 'Git Repo&nbsp;&nbsp;' + '<a href="' + window.gitRepoUrl
                    + '" >' + window.gitRepoUrl + '</a>'].join(
                    '<br>')
                textFormat: Text.RichText
                x: 30
                anchors.verticalCenter: parent.verticalCenter
                wrapMode: Text.Wrap
                onLinkActivated: {
                    Qt.openUrlExternally(link)
                }
            }
        }
    }
}
