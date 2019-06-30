import QtQuick 2.0
import QtQuick.Controls 1.4 as Controls
import "../js/app.js" as App
Controls.Menu {
    id: root
    Component {
        id: menuComp
        Controls.MenuItem {
            property var m
            text: m.name
            onTriggered: {
                 m.cb(m.name)
            }
        }
    }
    Component.onCompleted: {
        App.appState.globalRightMenu = root
    }
    function show(menus) {
        root.clear()
        menus.forEach(function(item) {
            var ins = menuComp.createObject(root, {m: item})
            root.insertItem(root.items.lenght, ins)
        })
        root.popup()
    }
}
