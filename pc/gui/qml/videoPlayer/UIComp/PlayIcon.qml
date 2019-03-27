import QtQuick 2.0
import QtQuick.Controls 2.3
import "../../js/app.js" as App

AniIcon {
    id: playIcon
    property var player: App.appState.player
    icon.source: !player.playing ? '../icons/pause.svg' : '../icons/play.svg'

    Action {
        onTriggered: {
            player.tooglePlay()
            playIcon.ani()
        }
        shortcut: ' '
    }

}
