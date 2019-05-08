import QtQuick 2.0
import QtQuick.Controls 2.1
import '../widget'

Button {
    id: iconBtn
    width: 25
    height: width
    anchors.verticalCenter: parent.verticalCenter
    property alias iconType: icon.type
    property alias title: toolTip.text
    property color color: 'black'
    property real lighter: 1.5
    ToolTip {
        id: toolTip
        show: iconBtn.hovered
        text: ''
    }
    IconFont {
        id: icon
        type: ''
        width: parent.width
        color: iconBtn.hovered ? iconBtn.color : Qt.lighter(iconBtn.color, lighter)
    }
    display: AbstractButton.IconOnly
    background: Item {
    }
}
