import QtQuick 2.8
import QtQuick.Controls 2.3
import "../../js/app.js" as App
import "../../js/global.js" as G

AniIcon {
    id: volumIcon
    property var player: App.appState.player
    property bool isUp: false
    icon.source: isUp ? '../icons/volume-up.svg' : '../icons/volume-down.svg'

    function changeVolume(direction) {
        var v = player.volume + direction * 0.05
        if (v > 1) v = 1
        else if (v  < 0) v = 0
        player.volume = v
        player.customerEvent('action.tips', '音量: ' + Math.round(v * 100) + '%')
        isUp = direction === 1
        ani()
    }

    Action {
        onTriggered: {
            volumIcon.changeVolume(-1)
        }
        shortcut: 'Down'
    }
    Action {
        onTriggered: {
            volumIcon.changeVolume(1)
        }
        shortcut: 'Up'
    }

}
