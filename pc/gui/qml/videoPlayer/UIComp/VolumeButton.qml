import QtQuick 2.0
import QtQuick.Controls 2.3
import "../../js/global.js" as G
import "../../js/app.js" as App

Button {
    id: playPauseButton
    property var player: App.appState.player
    icon.source: {
        if (player.muted)
            return '../icons/volume-mute.svg'
        if (player.volume < 0.5)
            return '../icons/volume-down.svg'
        return '../icons/volume-up.svg'
    }
    hoverEnabled: true
    icon.color: hovered ? 'white' : '#ddd'
    display: AbstractButton.IconOnly
    onClicked: {
        player.toogleMute()
        focus = false
    }
    background: Item {
    }
}
