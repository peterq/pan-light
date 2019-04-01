import QtQuick 2.0
import QtQuick.Controls 2.0

Button {
    property string type
    property var images: ["txt","doc","ppt","pptx","pdf",
        "htm","html","docx","xls","srt","xlsx","mp3","wma",
        "wav","ogg","ape","midi","flac","aac","wmv","asf","asx",
        "rm","rmvb","3gp","mpg","mpe","mpeg","mp4","m4v","mov","avi",
        "dat","vob","mkv","flv","fla","swf","jpeg","jpg","gif","bmp",
        "tif","png","psd","svg","pcx","wmf","emf","dxf","eps","tga","cdr",
        "zip","rar","7z","cab","iso","tar","unknown","dir"]

    icon.source: '../assets/images/icons/file/' +
            (images.indexOf(type) > -1 ? type : 'unknown') +'.svg'
    smooth: true
    icon.color: Qt.rgba(0 / 250, 140 / 255, 238 / 255, 1)
    icon.width: width
    icon.height: height
    display: AbstractButton.IconOnly
    enabled: false
    background: Item{}
}
