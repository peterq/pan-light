import QtQuick 2.11
import QtGraphicalEffects 1.0
import QtQuick.Controls 2.1
import "../js/util.js" as Util
import "../js/app.js" as App

Item {
    id: root

    property var meta
    property int idx
    property var listComp
    width: parent.width
    height: content.height + 20

    Rectangle {
        id: content
        width: parent.width
        height: 200
        color: 'white'
        Column {
            padding: 10
            spacing: 10
            Row {
                spacing: 10
                Rectangle {
                    width: 50
                    height: width
                    radius: width / 2
                    anchors.verticalCenter: parent.verticalCenter
                    color: 'orange'
                    Image {
                        id: avatar
                        smooth: true
                        visible: false
                        anchors.fill: parent
                        source: root.meta.user.avatar
                        antialiasing: true
                    }
                    Rectangle {
                        id: mask
                        anchors.fill: parent
                        radius: width / 2
                        visible: false
                        antialiasing: true
                        smooth: true
                    }
                    OpacityMask {
                        anchors.fill: avatar
                        source: avatar
                        maskSource: mask
                        antialiasing: true
                        visible: avatar.status === Image.Ready
                    }
                }

                Column {
                    Text {
                        text: root.meta.user.mark_username
                        font.pointSize: 12
                    }
                    Text {
                        text: Util.unixTime(root.meta.share_at)
                        font.pointSize: 10
                        color: '#aaa'
                    }
                }
            }

            Text {
                width: content.width - 30
                wrapMode: Text.WrapAnywhere
                text: root.meta.title
            }

            Row {
                spacing: 10
                Tag {
                    text.text: Util.humanSize(root.meta.file_size)
                }
                Tag {
                    border.color: 'red'
                    text.text: '热度指数: ' + root.meta.hot_index
                }
                Tag {
                    visible: root.meta.official
                    border.color: '#409EFF'
                    text.text: '官方认证'
                }
            }

            Row {
                spacing: 15
                Button {
                    id: btnDown
                    text: '下载'
                    onClicked: {
                        btnDown.text = '请稍后...'
                        btnDown.enabled = false
                        root.clickDownload().then(function () {
                            btnDown.text = '已添加至下载队列'
                        }).catch(function (err) {
                            btnDown.text = '下载'
                            btnDown.enabled = true
                            throw err
                        })
                    }
                }
                Button {
                    id: btnSave
                    text: '转存'
                    onClicked: {
                        btnSave.text = '请稍后...'
                        btnSave.enabled = false
                        root.clickSave().then(function () {
                            btnSave.text = '已转存'
                        }).catch(function (err) {
                            btnSave.text = '转存'
                            btnSave.enabled = true
                            throw err
                        })
                    }
                }
                Button {
                    id: btnPlay
                    visible: Util.isVideo(root.meta.title)
                    text: '播放'
                    onClicked: {
                        btnPlay.text = '请稍后...'
                        btnPlay.enabled = false
                        root.clickPlay().finally(function () {
                            btnPlay.text = '播放'
                            btnPlay.enabled = true
                        })
                    }
                }
            }
        }
    }

    function hit() {
        Util.api('share-hit', {id: meta._id})
    }

    function clickDownload() {
        hit()
        var fid = ['share', meta.md5, meta.slice_md5, meta.file_size].join('.')
        return Util.Promise.resolve().then(function () {
            return Util.callGoAsync('pan.link', {
                                        "fid": fid
                                    })
        }).then(function () {
            return App.appState.transferComp.addDownloadShare(meta.md5,
                                                              meta.slice_md5,
                                                              meta.file_size,
                                                              meta.title)
        }).catch(function (err) {
            topIndicator.fail(err)
            throw err
        })
    }

    function clickPlay() {
        hit()
        var fid = ['share', meta.md5, meta.slice_md5, meta.file_size].join('.')
        return Util.Promise.resolve().then(function () {
            return Util.callGoAsync('pan.link', {
                                        "fid": fid
                                    })
        }).then(function () {
            return Util.playVideoByLink(meta.title, Util.videoAgentLink(fid))
        }).catch(function (err) {
            topIndicator.fail(err)
            throw err
        })
    }

    function clickSave() {
        hit()
        return Util.Promise.resolve().then(function () {
            return Util.callGoAsync('pan.save.md5', {
                                        "md5": meta.md5,
                                        "sliceMd5": meta.slice_md5,
                                        "fileSize": meta.file_size,
                                        "path": '/pan-light-share/' + meta.title.replace('/', '_')
                                    })
        }).then(function (serverPath) {
            topIndicator.success('已保存到: ' + serverPath)
        }).catch(function (err) {
            topIndicator.fail(err)
            throw err
        })
    }
}
