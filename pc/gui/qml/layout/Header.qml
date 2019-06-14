import QtQuick 2.0
import "../comps"
import "../widget"
Rectangle {
    id: root
    color: "#EEEEF6"
    width: parent.width
    height: 80
    property alias tabs: tabBtns.btns
    property alias activeTab: tabBtns.activeTab
    signal activeChange(string tab)
    MoveWindow{
        anchors.fill: parent
    }
    // 下边界
    Rectangle {
        color: "#D7D7DE"
        height: 1
        width: parent.width
        anchors.bottom: parent.bottom
    }
    // logo
    Image {
        id: logo
        source: "../assets/images/pan-light-1.png"
        fillMode: Image.Stretch
        height: parent.height * 0.6
        y: (parent.height - height) * 0.5 * 0.6
        anchors.left: parent.left
        anchors.leftMargin: 20
        width: height
    }
    // app name
    Text {
        id: appName
        text: "pan-light"
        anchors.verticalCenter: logo.verticalCenter
        anchors.left: logo.right
        anchors.leftMargin: 10
        font.pointSize: 18
        font.bold: true
    }
    // tab 按钮
    Row {
        id: tabBtns
        property int btnWidth: 100
        property color activeColor: '#6441FF'
        property color inactiveColor: 'black'
        property string activeTab: '我的网盘'
        property var btns: ['我的网盘', '传输列表', '探索']
        signal activeChange(string tab)
        onActiveChange: {
            activeTab = tab
            root.activeChange(tab)
        }
        anchors.verticalCenter: parent.verticalCenter
        anchors.left: appName.right
        anchors.leftMargin: 40
        height: parent.height * 0.8
        Repeater {
            id: tabBtnRepeater
            model: parent.btns
            Item {
                id: tabItem
                property bool active: modelData === tabBtns.activeTab
                height: tabBtns.height
                width: tabBtns.btnWidth
                Text {
                    text: modelData
                    color: {
                       if (parent.active)
                           return tabBtns.activeColor
                       if (ma.containsMouse)
                           return Qt.lighter(tabBtns.activeColor, 1.2)
                       return tabBtns.inactiveColor
                    }
                    anchors.centerIn: parent
                    font.pointSize: 12
                    MouseArea {
                        id: ma
                        hoverEnabled: true
                        anchors.fill: parent
                        onClicked: !tabItem.active && tabBtns.activeChange(modelData)
                        cursorShape: tabItem.active ? Qt.ArrowCursor
                                                    : Qt.PointingHandCursor
                    }
                }
            }
        }
    }
    // 激活tab的下划线
    Rectangle {
        color: tabBtns.activeColor
        height: 3
        width: tabBtns.btnWidth
        anchors.bottom: tabBtns.bottom
        anchors.left: tabBtns.left
        anchors.leftMargin: {
            var idx = tabBtns.btns.findIndex(function (tab) {
                return tab === tabBtns.activeTab
            })
            return idx * tabBtns.btnWidth
        }

        Behavior on anchors.leftMargin {
            NumberAnimation {duration: 200; easing.type: Easing.OutQuad}
        }
    }

}
