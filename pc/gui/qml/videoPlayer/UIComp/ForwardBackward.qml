import QtQuick 2.8
import QtQuick.Controls 2.3
import "../../js/app.js" as App
import "../../js/global.js" as G

AniIcon {
    id: forwardBackward
    property var player: App.appState.player
    property bool isBack: false
    icon.source: isBack ? '../icons/backward.svg' : '../icons/forward.svg'
    property int offset: 0
    property real moveCheck: 0

    // 防止过快切换时间, 导致大量网络请求
    function moveLater(direction) {
        offset += direction * 10e3
        isBack = direction === -1
        if (player.position + offset < 0)
            offset = -player.position
        else if (player.position + offset > player.duration)
            offset = player.duration - player.position
        ani()
        player.customerEvent('action.tips', '定位到 ' + timeFormat( player.position + offset))
        moveCheck = Math.random()
        var t = moveCheck
        G.setTimeout(function() {
            if (t == moveCheck) {
                player.seekAbs(player.position + offset)
                offset = 0
            }
        }, 500)
    }

    function timeFormat(time) {
        var sec = Math.floor(time / 1000)
        var hours = Math.floor(sec / 3600)
        var minutes = Math.floor((sec - hours * 3600) / 60)
        var seconds = sec - hours * 3600 - minutes * 60
        var hh, mm, ss
        if (hours.toString().length < 2)
            hh = "0" + hours.toString()
        else
            hh = hours.toString()
        if (minutes.toString().length < 2)
            mm = "0" + minutes.toString()
        else
            mm = minutes.toString()
        if (seconds.toString().length < 2)
            ss = "0" + seconds.toString()
        else
            ss = seconds.toString()
        return hh + ":" + mm + ":" + ss
    }
    Action {
        onTriggered: {
            moveLater(-1)
        }
        shortcut: 'Left'
    }
    Action {
        onTriggered: {
            moveLater(1)
        }
        shortcut: 'Right'
    }

}
