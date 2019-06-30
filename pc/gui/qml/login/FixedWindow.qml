import QtQuick 2.0
import QtQuick.Window 2.2
Window {
    width: 300
    height: 340
    minimumHeight: height
    maximumHeight: height
    minimumWidth: width
    maximumWidth: width
    flags: Qt.MSWindowsFixedSizeDialogHint | Qt.WindowTitleHint | Qt.WindowCloseButtonHint
           | Qt.WindowModal | Qt.Dialog
    visible: false
    modality: Qt.WindowModal
}
