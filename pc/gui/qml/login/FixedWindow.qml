import QtQuick 2.0
import QtQuick.Window 2.2
Window {
    width: 300
    height: 340
    minimumHeight: height
    maximumHeight: height
    minimumWidth: width
    maximumWidth: width
    flags: Qt.Dialog | Qt.WindowModal |Qt.WindowCloseButtonHint
    visible: false
    modality: Qt.WindowModal
}
