# pan-light

- [环境搭建](DevEnv/introduction.md)
    - [golang](DevEnv/golang.md)
    - [qt5](DevEnv/qt5.md)
    - [protobuf](DevEnv/protobuf.md)

- 开始
    - [设计思路-客户端部分](Start/design-pc.md)   
    - [设计思路-在线演示系统](Start/design-demo.md)   
    - [目录结构](Start/directory.md)   
    - [项目初始化](Start/init.md)
    
- [Qt & Go](LangBind/introduction.md)
    - [Go中使用Qt](LangBind/qt_in_go.md)
    - qml & go通信
        - [方法互调](LangBind/qml/call.md)
        - [qml异步调用go](LangBind/qml/async.md)
    - [用plugin加速编译](LangBind/plugin.md)
    
- [下载器实现](Downloader/introduction.md)
    - [独立包设计](Downloader/package.md)
    - [并发下载](Downloader/parallel.md)
    - [速度统计](Downloader/speed_count.md)
    - [断点续传](Downloader/point_continue.md)

- [服务端相关](Server/introduction.md)
    - [频率限制器](Server/throttle.md)
    - [时间轮算法](Server/time_wheel.md)
    - [用户认证](Server/auth.md)
    - [参数处理](Server/param.md)
    - [错误处理](Server/error.md)
    
- [在线演示系统](Demo/introduction.md)
    - rpc实现
    - docker & vnc
    - web rtc
    - 即时通信
    - 排队系统
    - 匿名聊天-武侠角色