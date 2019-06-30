import QtQuick 2.0
import Qt.labs.platform 1.0

FileDialog {
    id: fileDialog
    fileMode: FileDialog.SaveFile
    title: '选择保存路径'
    nameFilters: ['全部文件 (*.*)']
    folder: defaultFolder || StandardPaths.standardLocations(
                StandardPaths.MoviesLocation)[0]
    property string defaultFolder: ''
    property string fileName: ''
    property string ext: fileName.split('.').pop()
    property var resolve: function (v) {
        console.log(v)
    }
    property var reject: function (v) {
        console.log(v)
    }
    onAccepted: {
        console.log(file.toString())
        var savePath = file.toString().replace('file://'
                                               + (Qt.platform.os == "windows" ? '/' : '')
                                               , '')
        if (Qt.platform.os == "windows") {
            savePath = savePath.split('/').join('\\')
        }
        resolve(savePath)
    }
    onRejected: {
        reject('用户取消选择保存路径')
    }
    Component.onCompleted: {
        currentFile = folder +  (Qt.platform.os == "windows" ? '\\' : '/') + fileName
    }
}
