import QtQuick 2.0

QtObject {
    id: dataObj
    property string $key: 'pan-light'

    Component.onCompleted: {
        for (var k in dataObj) {
            if (k.length > 7 && k.slice(k.length -7, k.length) === 'Changed')
                continue
            dataObj[k + 'Changed'].connect(function() {
                console.log('change', k, dataObj[k])
            })
        }
    }
}
