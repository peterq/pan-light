import QtQuick 2.0

Timer {
    id: timer
    property var cb
    onTriggered: {
        cb()
        timer.destroy()
    }
    Component.onCompleted: {
        timer.start()
    }
}
