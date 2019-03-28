import QtQuick 2.0
import QtQuick.Controls 2.3
import '../../js/global.js' as G
import '../../js/app.js' as App
Button {
    id: btn
    property var player: App.appState.player
    icon.source: '../icons/rotate.svg'
    hoverEnabled: true
    icon.color: hovered ? 'white' : '#ddd'
    display: AbstractButton.IconOnly
    onClicked: {
        player.rotateVideo()
        focus = false
    }
    background: Item {}
}
