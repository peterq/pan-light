import QtQuick 2.0
import QtQuick.Controls 2.2
Button {
    property string tooltip: ''
    property string iconImage: ''
    Text {
        id: name
        text: tooltip
    }
}
