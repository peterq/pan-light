import QtQuick 2.0
import QtQuick.Controls 2.1

Rectangle {
    clip: true
    Column {
        id: col
        spacing: 10
        anchors.fill: parent
        User {
            id: user
        }
        DiskUsage {
            id: diskUsage
        }
        PathCollection {
            height: Math.min(
                        500,
                        parent.height - user.height - diskUsage.height - 10 * 2)
        }
        Item {
            width: 1
            height: 30
        }
    }
    Text {
        visible: col.implicitHeight <= parent.height
        text: '★ <a href="https://github.com/peterq/pan-light">给作者点个star</a>'
        textFormat: Text.RichText
        height: 30
        x: 30
        anchors.bottom: parent.bottom
        wrapMode: Text.Wrap
        onLinkActivated: {
            Qt.openUrlExternally(link)
        }
    }
}
