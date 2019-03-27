import QtQuick 2.0
import QtQuick.Controls 1.4 as Controls

Controls.Menu {
    id: root
    property var menus: []
    Component {
        id: menuComp
        Controls.MenuItem {
            property var m: null
            text: m.name
            onTriggered: {
                 m.cb(m.name)
            }
        }
    }
    Component.onCompleted: {
        menus.forEach(function(item) {
            var ins = menuComp.createObject(root, {m: item})
            items.push(ins)
        })
    }
}
