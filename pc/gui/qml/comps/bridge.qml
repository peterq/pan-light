import QtQuick 2.0				//needed for js
import PanLight 1.0	//ItemTemplate

BridgeComp {
        id: root

        width: 250
        height: 450
        visible: false
        someString: "ItemTemplateString"

        Component.onCompleted: {
            root.logMsg('bridge ready')
            root.test(function(){})
        }
    }
