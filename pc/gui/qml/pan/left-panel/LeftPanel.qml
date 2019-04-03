import QtQuick 2.0
import QtQuick.Controls 2.1

Item {
    Column {
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
    }
}
