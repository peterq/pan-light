import QtQuick 2.0
import QtQuick.Controls 2.3
import '../../js/global.js' as G
import '../../js/app.js' as App
Button {
    id: playPauseButton
    property var player: App.appState.player
    icon.source: player.playing ? '../icons/pause.svg' : '../icons/play.svg'
    hoverEnabled: true
    icon.color: hovered ? 'white' : '#ddd'
    display: AbstractButton.IconOnly
    onClicked: {
        player.tooglePlay()
        focus = false
    }
    background: Item {
    }
}
