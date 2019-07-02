import QtQuick 2.0
import "../comps"

Item {
    id: root
    property var tasbs: Object.keys(tasbsMap)
    property var tasbsMap: {
        "newest": '最新',
        "hottest": '最热',
        "official": '官方'
    }
    property string currentTab: 'newest'

    signal active

    onCurrentTabChanged: {
        listRepeater.itemAt(tasbs.indexOf(currentTab)).active()
    }


    onActive: {
        listRepeater.itemAt(tasbs.indexOf(currentTab)).active()
    }

    DataSaver {
        $key: 'page-explorer'
        property alias showTip: tip.showTip
    }

    Column {
        anchors.fill: parent
        Rectangle {
            id: tip
            property string showTip: 'show'
            visible: showTip === 'show'
            color: '#409EFF'
            width: parent.width
            height: visible ? textTip.implicitHeight + 20 : 0
            Text {
                id: textTip
                color: 'white'
                width: parent.width - 40
                anchors.centerIn: parent
                wrapMode: Text.WrapAnywhere
                text: '开放型的资源广场, 突破版权文件分享限制(被和谐的违规文件除外).实验性功能, 暂只支持分享大于256k的文件.'
            }
            IconButton {
                iconType: 'close'
                color: 'red'
                anchors.right: parent.right
                anchors.rightMargin: 10
                anchors.verticalCenter: parent.verticalCenter
                onClicked: {
                    tip.showTip = ''
                }
            }
        }

        Rectangle {
            id: tabsBar
            width: parent.width
            height: 60
            Rectangle {
                width: parent.width
                height: 1
                color: '#ddd'
                anchors.bottom: parent.bottom
            }
            Item {
                id: typeTab
                property int padding: 4
                property int tabWidth: 80
                width: tabWidth * tabs.length + 2 * padding
                height: parent.height * 0.7
                anchors.verticalCenter: parent.verticalCenter
                x: 10
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
                    x: typeTab.padding + root.tasbs.indexOf(
                           root.currentTab) * typeTab.tabWidth
                    Behavior on x {
                        PropertyAnimation {
                            duration: 400
                            easing.type: Easing.OutCubic
                        }
                    }
                }

                Repeater {
                    model: root.tasbs
                    delegate: Item {
                        width: typeTab.tabWidth
                        height: typeTab.height - 2 * typeTab.padding
                        anchors.verticalCenter: parent.verticalCenter
                        x: typeTab.padding + index * typeTab.tabWidth
                        Text {
                            text: root.tasbsMap[modelData]
                            anchors.centerIn: parent
                        }
                        MouseArea {
                            anchors.fill: parent
                            onClicked: {
                               root.currentTab = modelData
                            }
                        }
                    }
                }
            }


        }

        Repeater {
            id: listRepeater
            model: root.tasbs
            ShareList {
                id: shareList
                type: modelData
                height: parent.height - tip.height - tabsBar.height
                width: parent.width
                visible: root.currentTab == type
            }
        }
    }
}
