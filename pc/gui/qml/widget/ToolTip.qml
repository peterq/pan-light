import QtQuick 2.0
import '../js/util.js' as Util
Item {
    id: root
    property string text: ''
    property int delay: 800
    property bool show: false
    property real showId: 0
    property var tipIns: Util.tooTip()

    onShowChanged: {
        if (!show) {
            tipIns.hide(showId)
            showDelay.stop()
        } else {
            if (text)
                showDelay.restart()
        }
    }

    Timer {
        id: showDelay
        running: false
        interval: root.delay
        onTriggered: {
            if (root.show)
                root.showId = tipIns.show(root.text)
        }
    }
}
