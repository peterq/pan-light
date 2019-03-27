import QtQuick 2.0
import QtQuick.Controls 2.3
import Qt.labs.platform 1.0
import "../../js/global.js" as G
import "../../js/app.js" as App

Button {
    id: btn
    property var player: App.appState.player
    icon.source: '../icons/open-file.svg'
    hoverEnabled: true
    icon.color: hovered ? 'white' : '#ddd'
    display: AbstractButton.IconOnly
    onClicked: {
        fileDialog.open()
    }
    background: Item {
    }

    DataSaver {
        $key: 'video.open-file'
        property alias lastFolder: fileDialog.lastFolder
    }

    FileDialog {
        id: fileDialog
        fileMode: FileDialog.OpenFile
        title: '播放本地视频'
        nameFilters: ["视频文件 (*.ts *.mp4 *.avi *.flv *.mkv *.3gp)", "全部文件 (*.*)"]
        options: FileDialog.ReadOnly
        folder: lastFolder || StandardPaths.standardLocations(StandardPaths.MoviesLocation)[0]
        property string lastFolder: ''
        onAccepted: {
            var sep = Qt.platform.os == "windows" ? '\\' : '/'
            var t = String.prototype.split.call(file, sep)
            var filename = t.pop()
            lastFolder = t.join(sep)
            player.playVideo(filename, file)
        }
    }
}
