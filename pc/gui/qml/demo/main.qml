import QtQuick 2.0
import QtQuick.Window 2.3
import '../js/util.js' as Util

Window {
    id: mainWindow
    width: Screen.desktopAvailableWidth * 0.23
    height: width * 0.8
    x: Screen.desktopAvailableWidth - width - 10
    y: 10
    flags: Qt.WA_TranslucentBackground | Qt.WA_TransparentForMouseEvents| Qt.FramelessWindowHint
                      | Qt.WindowSystemMenuHint
                       | Qt.WindowStaysOnTopHint | Qt.X11BypassWindowManagerHint
    color: 'transparent'
    modality: Qt.WindowModal
    visible: shuttingMsg == ''
    property int cursorIndex: -1
    property int itemHeight: 50
    property string nickname: 'nickname'
    property int endTime: +(new Date) / 1000
    property int currentTime: +(new Date) / 1000
    property string shuttingMsg: ''

    Timer {
        running: true
        repeat: true
        onTriggered: {
            currentTime = +(new Date) / 1000
        }
    }

    Component.onCompleted: {
        var conf = Util.callGoSync('conf')
        endTime = conf.endTime
        nickname = conf.nickname
        Util.callGoAsync('shutdownMsg')
            .then(function (msg) {
                shuttingMsg = msg
            })
    }

    Rectangle {
        anchors.fill: parent
        radius: 20
        color: Qt.rgba(0, 107/255, 1, 0.85)
        Column {
            anchors.fill: parent
            spacing: 5
            Text {
                text: 'pan-light 在线体验'
                anchors.horizontalCenter: parent.horizontalCenter
                font.pointSize: 25
                color: 'white'
            }
            Rectangle {
                height: 2
                width: parent.width
                color: '#0013ff'
            }
            Text {
                text: '本次体验时长剩余 :'
                anchors.horizontalCenter: parent.horizontalCenter
                font.pointSize: 16
                color: 'white'
            }
            Text {
                text: {
                    var du = endTime - currentTime
                    if (du < 0) du = 0
                    var seconds = du % 60
                    var minute = ~~((du - seconds) / 60)
                    if (seconds < 10) seconds = '0' + seconds
                    if (minute < 10) minute = '0' + minute
                    return minute + ':' + seconds
                }
                anchors.horizontalCenter: parent.horizontalCenter
                font.pointSize: 50
                color: 'white'
            }
            Text {
                text: '<font color="black"">*</a> 当前操作用户: ' + nickname
                    + '<br><font color="black"">*</a> 由于网络和远程pc的配置等原因, 部分功能无法使用'
                    + '<br><font color="black"">*</a> 为达到最佳体验效果建议下载安装体验'
                font.pointSize: 16
                x: 20
                color: 'white'
                width: parent.width - 40
                wrapMode: Text.WrapAnywhere
            }
        }
    }

    MouseArea {
        property point clickPos: "0,0"
        id: dragRegion
        anchors.fill: parent
        drag.minimumX: 0
        drag.maximumX: Screen.desktopAvailableWidth - mainWindow.width
        drag.minimumY: 0
        drag.maximumY: Screen.desktopAvailableHeight - mainWindow.heigh
        onPressed: {
            mainWindow.requestActivate()
            clickPos = Qt.point(mouseX, mouseY)
        }

        onPositionChanged: {
            var delta = Qt.point(mouse.x - clickPos.x, mouse.y - clickPos.y)
            mainWindow.x += delta.x
            mainWindow.y += delta.y
        }
    }

    Window {
        width: Screen.desktopAvailableWidth
        height: Screen.desktopAvailableHeight * 0.5
        y: Screen.desktopAvailableHeight * 0.25
        flags: Qt.WA_TranslucentBackground | Qt.WA_TransparentForMouseEvents| Qt.FramelessWindowHint
                          | Qt.WindowSystemMenuHint
                           | Qt.WindowStaysOnTopHint | Qt.X11BypassWindowManagerHint
        color: '#0013ff'
        modality: Qt.WindowModal
        visible: shuttingMsg != ''
        Text {
            text: shuttingMsg
            anchors.centerIn: parent
            font.pointSize: 50
            color: 'white'
        }
    }
}
