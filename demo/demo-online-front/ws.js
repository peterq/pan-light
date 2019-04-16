var ws = require("nodejs-websocket")

ws.createServer(function(conn){
    conn.on("text", function (str) {
        console.log("收到的信息为:"+str)
    })
    conn.on("close", function (code, reason) {
        console.log("关闭连接")
    })
    conn.on("error", function (code, reason) {
        console.log("异常关闭")
    })
    conn.on('binary', function(inStream) {
        // 创建空的buffer对象，收集二进制数据
        var data = new Buffer(0)
        // 读取二进制数据的内容并且添加到buffer中
        inStream.on('readable', function() {
            var newData = inStream.read()
            if (newData)
                data = Buffer.concat([data, newData], data.length + newData.length)
        })
        inStream.on('end', function() {
            console.log(data)
        })
    })
}).listen(8001)
console.log("WebSocket建立完毕")