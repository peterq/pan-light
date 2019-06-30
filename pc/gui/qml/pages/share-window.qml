import QtQuick 2.0
import QtQuick.Window 2.2
import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3
import "../js/app.js" as App
import "../js/util.js" as Util
import "../widget"

Window {
    id: window
    flags: Qt.MSWindowsFixedSizeDialogHint | Qt.WindowTitleHint | Qt.WindowCloseButtonHint
           | Qt.WindowModal | Qt.Dialog
    modality: Qt.ApplicationModal
    title: '分享到资源广场'
    minimumHeight: height
    minimumWidth: width
    maximumHeight: height
    maximumWidth: width
    visible: true
    width: 550
    height: 300

    property bool submitting: false
    property var sliceMd5Promise
    property var meta
    property var timeMap: {
        "永久": 365,
        "7天": 7,
        "一个月": 30
    }
    property string fileExtention

    Component.onCompleted: {
        getSliceMd5()
    }

    function getSliceMd5() {
        sliceMd5Promise = Util.callGoAsync('pan.rapid.md5', {
                                               "fid": meta.fs_id
                                           }).then(function (sliceMd5) {
                                               console.log('slice md5',
                                                           sliceMd5)
                                               return sliceMd5
                                           })
        sliceMd5Promise.catch(function () {
            sliceMd5Promise = null
        })
    }

    TopIndicator {
        id: topIndicator
        z: 2
    }

    GridLayout {
        columns: 2
        width: parent.width - 20
        y: 20
        rowSpacing: 10
        anchors.horizontalCenter: parent.horizontalCenter
        Label {
            text: '分享文件名'
            width: parent.width * 30
            Layout.alignment: Qt.AlignRight
        }
        TextField {
            id: titleInput
            enabled: !window.submitting
            width: parent.width * 60
            placeholderText: "文件名"
            Component.onCompleted: {
                var t = meta.server_filename.split('.')
                window.fileExtention = t.pop()
                text = t.join('.')
            }
        }
        Label {
            text: '有效期'
            width: parent.width * 30
            Layout.alignment: Qt.AlignRight
        }
        ComboBox {
            id: selectDuraion
            enabled: !window.submitting
            model: Object.keys(window.timeMap)
        }

        Label {
        }
        Button {
            id: btnShare
            text: window.submitting ? '请稍后' : '分享'
            enabled: !submitting
            onClicked: {
                if (!titleInput.text)
                    return
                if (!sliceMd5Promise) {
                    getSliceMd5()
                }
                window.submitting = true
                sliceMd5Promise.then(function (sliceMd5) {
                    return Util.api('share', {
                                        "md5": window.meta.md5,
                                        "sliceMd5": sliceMd5,
                                        "title": titleInput.text + '.' + window.fileExtention,
                                        "duration": window.timeMap[selectDuraion.currentText],
                                        "fileSize": window.meta.size
                                    })
                })
                .then(function() {
                    topIndicator.success('分享成功')
                    return Util.sleep(1000)
                })
                .then(function(){
                    window.visible = false
                })
                .catch(function (err) {
                    topIndicator.fail(err.message)
                    window.submitting = false
                })
            }
        }
    }
}
