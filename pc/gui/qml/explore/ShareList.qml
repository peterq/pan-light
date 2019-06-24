import QtQuick 2.0
import QtQuick.Controls 2.1
import "../js/util.js" as Util
import "../widget"

Rectangle {
    id: root
    clip: true

    color: '#eee'

    property int pageSize: 10
    property int offset: 0
    property var loadPromise
    property string type: 'newest'
    property bool noMore: false
    signal active

    onActive: {
        if (offset === 0)
            loadMore()
    }

    TopIndicator {
        id: topIndicator
        z: 2
    }

    ListModel {
        id: listModel
    }

    Item {
        id: refreshIndicator
        width: parent.width
        height: 50
        visible: !listView.flicking
        y: -listView.contentY - refreshIndicator.height + listView.y
        Row {
            spacing: 5
            anchors.centerIn: parent
            BusyIndicator {
                visible: running
                running: refreshText.text === '刷新中'
                contentItem.implicitHeight: 25
                contentItem.implicitWidth: 25
            }
            Text {
                id: refreshText
                text: '下拉刷新'
                anchors.verticalCenter: parent.verticalCenter
            }
        }
        Timer {
            id: hideSuccessTipTimer
            property bool hide: true
            interval: 1000
            onTriggered: {
                hide = true
            }
        }
        states: [
            State {
                id: refreshDone
                name: "refreshDone"
                when: !hideSuccessTipTimer.hide
                PropertyChanges {
                    target: refreshText
                    text: '刷新成功'
                }
            },
            State {
                id: refreshing
                name: "refreshing"
                when: !!root.loadPromise && root.loadPromise.refresh
                StateChangeScript {
                    script: loadPromise.then(function () {
                        hideSuccessTipTimer.hide = false
                        hideSuccessTipTimer.restart()
                    })
                }
                PropertyChanges {
                    target: refreshText
                    text: '刷新中'
                }
            },
            State {
                id: refreshHold
                name: "refreshHold"
                when: -listView.contentY - refreshIndicator.height >= 0
                PropertyChanges {
                    target: refreshText
                    text: '释放立即刷新'
                }
            }

        ]
    }

    Component {
        id: listViewFooter
        Item {
            width: parent.width
            height: 50
            Text {
                text: {
                    if (root.type === 'hottest')
                        return '热度榜每分钟更新一次~'
                    if (root.loadPromise) {
                        return '加载中'
                    }
                    if (root.noMore)
                        return '没有更多了'
                    return '加载更多'
                }
                anchors.centerIn: parent
                MouseArea {
                    anchors.fill: parent
                    onClicked: {
                        if (parent.text === '加载更多')
                            root.loadMore()
                    }
                }
            }
        }
    }

    ListView {
        id: listView
        footer: listViewFooter
        y: ['刷新中', '刷新成功'].indexOf(refreshText.text) >= 0 ? refreshIndicator.height : 0
        width: parent.width
        height: parent.height
        model: listModel
        delegate: ShareItem {
            id: shareItem
            meta: listModel.get(index)
            idx: index
            listComp: root
        }
        ScrollBar.vertical: ScrollBar {
        }
        onDraggingChanged: {
            if (!dragging && refreshIndicator.state === 'refreshHold') {
                loadMore(true)
            }
        }
        Behavior on y {
            PropertyAnimation {
                duration: 150
                easing.type: Easing.Linear
            }
        }
    }

    function loadMore(refresh) {
        if (loadPromise)
            return
        if (!refresh && noMore) return
        loadPromise = Util.api('share-list', {
                                   "offset": refresh ? 0 : offset,
                                   "pageSize": pageSize,
                                   "type": type
                               })

        loadPromise = loadPromise.then(function (list) {
            if (refresh) {
                listModel.clear()
            }
            handleMore(list)
            return Util.sleep(100)
        }).catch(function (err) {
            topIndicator.fail(err.message)
            throw err.message
        }).finally(function () {
            loadPromise = null
        })
        loadPromise.refresh = refresh
    }

    function handleMore(list) {
        if (list.length === 0) {
            noMore = true
            return
        }
        noMore = false
        list.forEach(function (item) {
            listModel.append(item)
        })
        offset = listModel.get(listModel.count - 1).share_at
    }
}
