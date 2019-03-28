import QtQuick 2.0
import "../js/global.js" as G
import "../js/util.js" as Util

QtObject {
    id: dataObj
    property string $key: '$pan-light'

    Component.onCompleted: {
        // 必须指定数据存储前缀
        if ($key === '$pan-light')
            throw new Error('$key must be specified')
        // 检测是前缀否唯一
        if (G.dataSaverKeys[$key])
            throw new Error('duplicated $key: ' + $key)
        G.dataSaverKeys[$key] = true
        // 遍历属性, 从磁盘恢复
        for (var k in dataObj) {
            if (!dataObj[k + 'Changed'])
                continue
            var storageKey = ['comp-data', $key, k].join('.')
            var $default = dataObj[k]
            var v = Util.storageGet(storageKey, $default)
            if (v !== $default)
                dataObj[k] = v
            watch(k, storageKey)
        }
    }

    // 监听属性
    function watch(key, storageKey) {
        dataObj[key + 'Changed'].connect(function() {
            Util.storageSet(storageKey, dataObj[key])
        })
    }

    // 前缀释放
    Component.onDestruction: {
        delete G.dataSaverKeys[$key]
    }
}
