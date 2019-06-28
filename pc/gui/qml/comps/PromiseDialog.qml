import QtQuick.Dialogs 1.1
import QtQuick 2.11
import QtQuick.Window 2.2
import QtQuick.Controls 2.2
import QtQuick.Layouts 1.3
import '../js/util.js' as Util
import "../js/global.js" as G

Item {
    id: root
    property string confirmLabel: '确定'
    property string cancelLabel: '取消'
    property string title: '对话框'
    property var options
    property Item header
    property Item contentItem
    property var result
    property bool confirmBtnEnabled: true
    property var onClickConfirm: function(){return false}
    property int w
    property int h

    QtObject {
        id: promise
        property bool done: true
        property var promise
        property var resovle
        property var reject
    }

    Component {
        id: windowComp
        Item {
            parent: G.root
            Window {
                id: window
                flags:  Qt.MSWindowsFixedSizeDialogHint | Qt.WindowTitleHint | Qt.WindowCloseButtonHint
                        | Qt.WindowModal | Qt.Dialog
                modality: Qt.ApplicationModal
                title: root.title
                minimumHeight: 100
                minimumWidth: 100
                width: root.w
                height: root.h
                Dialog {
                    id: dialog
                    width: window.width
                    height: window.height
                    focus: true
                    header: root.header
                    contentItem: root.contentItem
                    footer: Item {
                        width: parent.width
                        height: 50
                        Row {
                            spacing: 10
                            height: parent.height
                            anchors.centerIn: parent
                            Button {
                                width: 80
                                height: 30
                                text: root.cancelLabel
                                onClicked: {
                                    promise.reject('cancle')
                                }
                            }
                            Button {
                                width: 80
                                height: 30
                                text: root.confirmLabel
                                visible: !!root.confirmLabel
                                enabled: root.confirmBtnEnabled
                                onClicked: {
                                    if (root.onClickConfirm()) {
                                        promise.resovle(root.result)
                                    }
                                }
                            }
                        }
                    }
                }

                onClosing: {
                    promise.reject('closed')
                }

                Component.onCompleted: {
                    visible = true
                    requestActivate()
                    dialog.open()
                }
            }

        }
    }

    Loader {
        id: windowLoader
        focus: true
    }

    function open(opt) {
        if (!promise.done)
            throw new Error('last dialog is not closed')
        promise.done = false
        promise.promise = new Util.Promise(function(res, rej) {
            options = opt
            promise.resovle = res
            promise.reject = rej
            windowLoader.sourceComponent = windowComp
        }).finally(function() {
            console.log(windowLoader.item)
            windowLoader.sourceComponent = null
            promise.done = true
        })
        return promise.promise
    }

    function forceClose() {
        if (!promise.done)
            promise.reject('force closed')
    }
}
