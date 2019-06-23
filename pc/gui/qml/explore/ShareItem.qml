import QtQuick 2.11
import QtGraphicalEffects 1.0
import QtQuick.Controls 2.1
import "../js/util.js" as Util

Item {
    id: root

    property var meta
    property int idx
    property var listComp
    width: parent.width
    height: content.height + 20

    Rectangle {
        id: content
        width: parent.width
        height: 200
        color: 'white'
        Column {
            padding: 10
            spacing: 10
            Row {
                spacing: 10
                Rectangle {
                    width: 50
                    height: width
                    radius: width / 2
                    anchors.verticalCenter: parent.verticalCenter
                    color: 'orange'
                    Image {
                        id: avatar
                        smooth: true
                        visible: false
                        anchors.fill: parent
                        source: root.meta.user.avatar
                        antialiasing: true
                    }
                    Rectangle {
                        id: mask
                        anchors.fill: parent
                        radius: width / 2
                        visible: false
                        antialiasing: true
                        smooth: true
                    }
                    OpacityMask {
                        anchors.fill: avatar
                        source: avatar
                        maskSource: mask
                        antialiasing: true
                        visible: avatar.status === Image.Ready
                    }
                }

                Column {
                    Text {
                        text: root.meta.user.mark_username
                        font.pointSize: 12
                    }
                    Text {
                        text: Util.unixTime(root.meta.share_at)
                        font.pointSize: 10
                        color: '#aaa'
                    }
                }
            }

            Text {
                width: content.width - 30
                wrapMode: Text.WrapAnywhere
                text: root.meta.title
            }

            Row {
                spacing: 10
                Tag {
                    text.text: Util.humanSize(root.meta.file_size)
                }
                Tag {
                    border.color: 'red'
                    text.text: '热度指数: ' + root.meta.hot_index
                }
                Tag {
                    visible: root.meta.official
                    border.color: '#409EFF'
                    text.text: '官方认证'
                }
            }

            Row {
               spacing: 15
               Button {
                   text: '下载'
               }
               Button {
                   text: '转存'
               }
               Button {
                   visible: Util.isVideo(root.meta.title)
                   text: '播放'
               }
            }
        }
    }
}
