import QtQuick 2.0
import "../js/global.js" as G
import "../js/util.js" as Util

QtObject {
    id: dataObj
    property string $key: '$pan-light'

    Component.onCompleted: {
        if ($key === '$pan-light')
            throw new Error('$key must be specified')
        if (G.dataSaverKeys[$key])
            throw new Error('duplicated $key: ' + $key)
        G.dataSaverKeys[$key] = true
        for (var k in dataObj) {
            if (k.length > 7 && k.slice(k.length -7, k.length) === 'Changed')
                continue
            var key = k
            var storageKey = ['comp-data', $key, key].join('.')
            dataObj[key + 'Changed'].connect(function() {
                // console.log('change', key, dataObj[key])
                Util.storageSet(storageKey, dataObj[key])
            })
            var $default = dataObj[key]
            var v = Util.storageGet(storageKey, $default)
            if (v !== $default)
                dataObj[key] = v
        }
    }

    Component.onDestruction: {
        delete G.dataSaverKeys[$key]
    }
}
