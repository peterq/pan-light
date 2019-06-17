import QtQuick 2.0
import QtQuick.Controls 2.0
import "../js/app.js" as App

Item {
    width: parent.width - 20
    anchors.horizontalCenter: parent.horizontalCenter
    height: 50
    property alias currentTab: typeTab.currentTab
    Text {
        id: taskLeft
        text: "剩余任务: " + App.appState.downloadingList.length
        anchors.verticalCenter: parent.verticalCenter
    }

    Item {
        id: typeTab
        property int padding: 4
        property int tabWidth: 80
        property var tabs: ['下载中', '已完成']
        property string currentTab: '下载中'
        width: tabWidth * tabs.length + 2 * padding
        height: parent.height * 0.8
        anchors.verticalCenter: parent.verticalCenter
        anchors.left: taskLeft.right
        anchors.leftMargin: 30
        Rectangle {
            anchors.fill: parent
            radius: 3
            color: '#eee'
        }

        Rectangle {
            radius: 3
            width: typeTab.tabWidth
            height: typeTab.height - 2 * typeTab.padding
            color: '#fff'
            anchors.verticalCenter: parent.verticalCenter
            x: typeTab.padding + parent.tabs.indexOf(
                   parent.currentTab) * typeTab.tabWidth
            Behavior on x {
                PropertyAnimation {
                    duration: 400
                    easing.type: Easing.OutCubic
                }
            }
        }

        Repeater {
            model: parent.tabs
            delegate: Item {
                width: typeTab.tabWidth
                height: typeTab.height - 2 * typeTab.padding
                anchors.verticalCenter: parent.verticalCenter
                x: typeTab.padding + index * typeTab.tabWidth
                Text {
                    text: modelData
                    anchors.centerIn: parent
                }
                MouseArea {
                    anchors.fill: parent
                    onClicked: {
                        typeTab.currentTab = modelData
                    }
                }
            }
        }
    }

    Row {
        anchors.right: parent.right
        anchors.verticalCenter: parent.verticalCenter
        spacing: 10
        Button {
            text: '全部开始'
            onClicked: {
                App.appState.downloadingListComp.startAll()
            }
        }
        Button {
            text: '全部暂停'
            onClicked: {
                App.appState.downloadingListComp.pauseAll()
            }
        }
        Button {
            text: '全部取消'
            visible: false
        }
    }
}
