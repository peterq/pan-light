import QtQuick 2.4
import "../pan"
import "../transfer"
import "../explore"
import "../js/util.js" as Util
Item {
    id: mainLayout
    property var tabsMap: {'我的网盘': pan, '传输列表': transfer, '探索': explore}
    property var tabs: ['我的网盘', '传输列表', '探索']
    property var colors: ['blue', 'red', 'green']
    property string activeTab: '我的网盘'
    anchors.fill: parent
    Header {
        id: header
        tabs: mainLayout.tabs
        activeTab: mainLayout.activeTab
        onActiveChange: {
            tabsCon.doSlideAnimation(tab)
            mainLayout.activeTab = tab
        }
    }
    function notiActive() {
        tabsMap[mainLayout.activeTab].active()
    }
    onActiveTabChanged: {
        notiActive()
    }

    Rectangle {
        id: tabsViewport
        width: parent.width
        height: parent.height - header.height
        y: header.height
        clip: true
        color: '#eaeaea'
        Text {
            id: tabBgAppName
            text: "pan-light"
            font.pointSize: 40
            anchors.centerIn: parent
        }
        Text {
            font.pointSize: 20
            color: 'red'
            text: "by PeterQ"
            anchors.left: tabBgAppName.right
            anchors.bottom: tabBgAppName.bottom
            anchors.leftMargin: 10
        }
        Row {
            id: tabsCon
            width: parent.width * tabs.length
            height: parent.height
            spacing: 0.2 * tabsViewport.width
            x: {
                if (slideChangeAnimation.doing) return x
                var idx = mainLayout.tabs.findIndex(function (tab) {
                    return tab === mainLayout.activeTab
                })
                return -(tabsViewport.width + tabsCon.spacing) * idx
            }
            function doSlideAnimation(newTab) {
                var idx = mainLayout.tabs.findIndex(function (tab) {
                    return tab === newTab
                })
                var newX = -(tabsViewport.width + tabsCon.spacing) * idx
                slideChangeAnimation.stop()
                tabSlideAnimation.from = x
                tabSlideAnimation.to = newX
                slideChangeAnimation.start()
            }

            SequentialAnimation {
                id: slideChangeAnimation
                property real speed: 1
                property bool doing: false
                onStarted: {
                    doing = true
                }
                onStopped: {
                    if (tabsConScale.scale === 1) doing = false
                }
                PropertyAnimation {
                   id: tabScalAnimation
                   target: tabsConScale
                   property: 'scale'
                   from: tabsConScale.scale
                   to: 0.5
                   duration: Math.abs(from - to) / 0.5 * 200 * slideChangeAnimation.speed
                }
                PauseAnimation {
                    duration: 150 * slideChangeAnimation.speed
                }
                PropertyAnimation {
                   id: tabSlideAnimation
                   target: tabsCon
                   property: 'x'
                   from: 0
                   to: 0
                   duration: Math.abs(from - to) / (tabsViewport.width + tabsCon.spacing) * 200 * slideChangeAnimation.speed || 1
                }
                PauseAnimation {
                    duration: 150 * slideChangeAnimation.speed
                }
                PropertyAnimation {
                   id: tabScalAnimation2
                   target: tabsConScale
                   property: 'scale'
                   from: 0.5
                   to: 1
                   duration: 200 * slideChangeAnimation.speed
                }
            }

            transform: Scale {
                id: tabsConScale
                property real scale: 1
                xScale: scale
                yScale: scale
                origin.x: -tabsCon.x + tabsViewport.width / 2
                origin.y: parent.height /2
            }

            Pan {
                id: pan
                width: tabsViewport.width
                height: tabsViewport.height
            }

            Transfer {
                id: transfer
                width: tabsViewport.width
                height: tabsViewport.height
            }

            Explore {
                id: explore
                width: tabsViewport.width
                height: tabsViewport.height
            }


        }
    }

    Component.onCompleted: {
        notiActive()
        console.log(Util.callGoSync('time'))
        Util.openDesktopWidget()
        Util.callGoAsync('wait', {time: 1})
            .then(function(data) {
                console.log('reslove', data)
            }, function(data) {
                console.log('reject', data)
            }, function(data) {
                console.log('progress', data)
            })
    }
}
