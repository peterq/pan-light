import QtQuick 2.4

Rectangle {
    width: text.implicitWidth + 10
    height: 25
    radius: 6
    color: 'transparent'
    border.color: 'blue'
    property alias text: text
    Text {
        id: text
        anchors.centerIn: parent
        color: parent.border.color
        font.pointSize: 10
    }
}
