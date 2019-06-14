import QtQuick 2.0
import "../js/util.js" as Util
import "../comps"
import QtQuick.Layouts 1.1
import QtQuick.Controls 2.2

Rectangle {
    anchors.fill: parent
    visible: false
    Component.onCompleted: {
        Util.event.on('init.not-login', function () {
            visible = true
        })
    }
    Loader {
        id: loginLoader
    }
    Component {
        id: wxLoginComp
        Wx {}
    }
    Component {
        id: baiduLoginComp
        Baidu {}
    }
    Component {
        id: qqLoginComp
        QQ {}
    }
    ColumnLayout {
        anchors.centerIn: parent
        spacing: 20
        Text {
            Layout.alignment: Qt.AlignCenter
            text: '请选择一种方式登录你的百度网盘'
        }


        RowLayout {
            Layout.alignment: Qt.AlignCenter
            spacing: 20
            Component {
                id: loginMethodComp
                IconFont {
                    width: 80
                    type: method
                    MouseArea {
                        anchors.fill: parent
                        onClicked: {
                            showLogin(method)
                        }
                    }
                }
            }
            Loader {
                property string method: 'wx'
                sourceComponent: loginMethodComp
            }
            Loader {
                property string method: 'baidu'
                sourceComponent: loginMethodComp
            }
            Loader {
                property string method: 'qq'
                sourceComponent: loginMethodComp
            }
        }
    }

    function showLogin(type) {
        if (loginLoader.sourceComponent) return
        var compMap = {
            'wx': wxLoginComp,
            'baidu': baiduLoginComp,
            'qq': qqLoginComp
        }
        loginLoader.sourceComponent = compMap[type]
        loginLoader.item.start()
        .then(function(){
            Util.event.fire('login.success', type)
            visible = false
        })
        .finally(function(){
            console.log('login finish')
            loginLoader.sourceComponent = null
        })
    }
}
