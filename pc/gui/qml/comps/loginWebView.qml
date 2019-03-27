import QtQuick 2.0
import QtQuick.Window 2.0
import QtWebView 1.1

Window {
    width: 1024
    height: 750
    visible: true
    WebView {
        id:webview
        anchors.fill: parent
        url: "http://www.baidu.com"
    }
}
