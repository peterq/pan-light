import QtQuick 2.0
import PanLight 1.0

BridgeComp {
    id: root

    width: 250
    height: 450
    visible: false
    someString: "ItemTemplateString"

    Component.onCompleted: {
        root.logMsg('bridge ready')
        root.test(function () {})
    }
}
