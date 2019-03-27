import QtQuick 2.0
import '../js/util.js' as Util
import '../js/app.js' as App
Item {

    property var loginSession: App.appState.loginSession
    Text {
        text: App.appState.loginSession.username
    }
}
