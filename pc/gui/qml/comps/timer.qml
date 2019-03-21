import QtQuick 2.0

Timer {
    id: timer
    property alias interval: timer.interval
    property var cb

    running: true
    onTriggered: {
        cb()
        timer.destroy()
    }
}
