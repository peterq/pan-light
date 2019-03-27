import QtQuick 2.0
import "../../js/app.js" as App

Item {
    Text {
        property var player: App.appState.player
        anchors.verticalCenter: parent.verticalCenter
        text: parent.timeFormat(player.position) +
              " / " + parent.timeFormat(player.duration)
        color: "white"
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
}
