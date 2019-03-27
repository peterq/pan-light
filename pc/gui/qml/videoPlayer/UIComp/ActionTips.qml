import QtQuick 2.0
import QtQuick.Controls 2.3
import "../../js/global.js" as G
import "../../js/app.js" as App

Text {
    id: playPauseButton
    property var player: App.appState.player
    visible: false
    property real hideCheck: 0
    text: 'pan-light'
    color: 'yellow'

    function show(msg) {
        text = msg
        visible = true
        hideCheck = Math.random()
        var t = hideCheck
        G.setTimeout(function() {
            if (t === hideCheck) {
                visible = false
            }
        }, 2e3)
    }

    Connections {
        target: player
        onCustomerEvent: {
            if (event != 'action.tips')
                return
            show(data)
        }
    }

}
