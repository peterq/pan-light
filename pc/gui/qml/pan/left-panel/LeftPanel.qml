import QtQuick 2.0
import QtGraphicalEffects 1.0
import QtQuick.Controls 2.1
import "../../comps"
import "../../widget"
import "../../js/util.js" as Util
import "../../js/app.js" as App

Item {

    property var loginSession: App.appState.loginSession

    User {
        id: user
    }
    // 网盘用量
   DiskUsage {
       id: usage
       anchors.top: user.bottom
       anchors.topMargin: 10
   }
}
