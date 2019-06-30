# 目录结构

## 三个模块

```
pan-ligth
├── demo/ // 在线演示系统模块
├── go.mod
├── LICENSE
├── pan-light.go // 开发时的辅助脚本, 把一些常用的命令集成一起方便使用
├── pan-light.pro // qt 工程文件, 编写qml代码时, 在Qt Creator ide中打开此文件.
├── pan-light.pro.user // 同上
├── pc/ // 客户端模块
├── qt/ // 用来将 qt 绑定到 go 的模块, 这个模块只在开发时运行, 程序跑起来不会运行
├── README.md
└── server/ // 服务端模块

```

### 客户端模块

```
pc
├── dep // 解决包互相依赖的问题, 以及初始化顺序的问题, 这个包不依赖其他包.
│   │          参照laravel的service provider处理方式, 把包分为注册和初始化2个时机
│   ├── dep.go 
│   ├── env-dev.go // 开发环境配置
│   ├── env.go // 全局配置变量, 服务端地址, 版本号等
│   └── env-prod.go // 正式环境配置
├── docker // 存放docker file, 这些容器用来在linux上打windows的包
│   ├── winodws_64_shared
│   │   └── Dockerfile
│   └── winodws_64_static
│       └── Dockerfile
├── downloader // 下载器包
│   ├── download-helper.go // 下载器帮助函数
│   ├── internal // protobuf 原型文件放在这里, 用来将下载状态序列化后持久化到磁盘.
│   │   └── types.proto
│   ├── manager.go // 下载管理器类, 管理所有任务的下载, 暂停, 重启软件后的恢复等
│   ├── segment.go // 分段下载用到的, 下载片段结构体, 记录某片段起始地址, 长度, 已下载完成的长度
│   ├── task.go // 下载任务结构体, 皴法下载链接, 文件保存地址等
│   └── worker.go // 下载工作协程一一对应的结构体, 一个并发对应一个worker, 一个任务中有多个work
├── functions // 这个包用来存放gui端, 也就是qml执行环境调用的函数. 这里一般不会直接实现相应功能, 而是调用其他包的函数实现, 这里只是做了一层封装.
│   ├── base.go // 软件重启, 获取配置等基本操作
│   ├── download.go // 下载模块的接口
│   ├── login.go // 几种登录方式的接口
│   ├── pan-api.go // 百度网盘的API, 获取文件链接, 列出目录下的文件等
│   ├── regitser.go // 用来把上面这些功能注册到qml
│   ├── regitser-plugin.go // plugin模式下的注册
│   └── testing.go
├── go.mod
├── go.sum
├── gui // 和qt打交道的模块, 其他包都不会和qt直接打交道. 
│   ├── bridge // 桥接, 连接 qml 环境和 go环境
│   │   ├── api_for_qml.go // 存放上面functions包的map
│   │   └── router.go // 路由, 把qml的调用请求从字符串解析到相应的go函数
│   ├── comp 
│   │   ├── BridgeComp.go // qml 原生组件, '继承'qt的c++类实现, 这个类联通了qml和go 2个执行环境
│   ├── gui.go // gui 界面初始化入口
│   ├── gui-plugine.go // plugin 模式下的gui界面初始化入口
│   ├── icon.ico // windows下的图标
│   ├── icon.rc
│   ├── icon_windows.syso
│   ├── qml // 界面描述文件, 也就是qml文件都存放在这里面, 详情在下方
│   │   └──...
│   └── qt-rpc // 存放几个全局变量, 解决包依赖用的
│       └── rpc.go
├── login // 登录功能实现
│   ├── baidu.go // 百度app扫码登录
│   ├── login-http-client.go // 登录用到的http客户端
│   ├── qq.go // qq扫码登录
│   └── wx.go // 微信扫码登录
├── pan-api // 把web版百度网盘提供的http接口转换成go包
│   ├── http-client.go // http客户端
│   ├── login-session.go // 登录会话结构体
│   └── pan-api.go // 调用接口具体实现
├── pan-download // 对于上面download的封装
│   ├── pan-download.go // 封装实现
│   └── video-agent.go // 一个简单的代理, 突破百度防盗链, 是的qt的视频播放器能播放网盘中的视频
├── pan-light-pc-dev.go // 开发模式, 即plugin模式程序入口
├── pan-light-pc.go // 程序main函数
├── server-api // 调用服务端模块的包
│   └── server-api.go
├── storage // 程序状态持久化包
│   ├── state.go // 全局状态结构体
│   ├── types.pb.go // protobuf生成的文件
│   └── types.proto // protobuf 原型文件
└── util // 工具包
    ├── file.go
    └── fn.go

```

#### 客户端界面文件结构
```
pc/gui/qml
├── assets 静态文件目录
│   └── images // 图片
│       ├── icons // icon图标
│       │   ├── baidu-cloud.svg
│       │   └ ...
│       └── pan-light-1.png
├── comps // 公共组件, js动态创建的组件
│   ├── Alert.qml // 警告弹窗
│   ├── bridge.qml // 和go交互的组件
│   ├── confirm-window.qml // 确认组件
│   ├── DataSaver.qml // 将界面信息, xy坐标等状态自动存盘的组件
│   ├── desktop-widget.qml // 桌面悬浮窗组件
│   ├── dialog.qml // 对话框
│   ├── IconButton.qml // 图标按钮
│   ├── IconFont.qml // 已弃用, 打算使用阿里的Iconfont字体图标
│   ├── loginWebView.qml // 一起用, 打算通过webview来实现网页登录
│   ├── Modal.qml // 忽略, 测试用的mask窗口
│   ├── PromiseDialog.qml // promise话的dialog
│   ├── prompt-window.qml // 用户输入弹框
│   ├── rightClickMenu.qml // 右键按钮通用组件
│   ├── select-save-path.qml // 文件保存路径选择组件
│   ├── timer.qml // 计时器组件, qml未提供setTimeout函数, 通过这个组件实现polyfill
│   └── tool-tip-window.qml // 鼠标选题显示提示的窗口
├── explore // 探索页面
│   ├── explore.qml // 探索页面布局
│   ├── ShareItem.qml // 分享内容展示组件
│   ├── ShareList.qml // 分享列表 list view
│   ├── TagForm.ui.qml  // 忽略
│   └── Tag.qml // 分享徽标
├── js // js 文件
│   ├── app.js // 组件全局状态
│   ├── appState.qml // 全局状态组件
│   ├── global.js // 全局变量
│   ├── promise.js // 带进度的promise实现
│   ├── transfer.js // 弃用 
│   └── util.js // qml组件调用的工具
├── layout // 界面布局
│   ├── Header.qml // 头部
│   └── Layout.qml // 布局
├── login // 登录方式的界面
│   ├── Baidu.qml
│   ├── FixedWindow.qml // 登录窗口公用组件
│   ├── Login.qml
│   ├── QQ.qml
│   └── Wx.qml
├── main.qml // 界面入口
├── pages // 一些小窗口
│   ├── about-window.qml // 关于窗口
│   ├── feedback-window.qml // 反馈窗口
│   ├── setting-window.qml // 设置窗口
│   └── share-window.qml // 分享交互窗口
├── pan // 我的网盘界面
│   ├── FileIcon.qml // 文件图标组件
│   ├── FileList.qml // 文件列表
│   ├── FileNode.qml // 文件item
│   ├── left-panel // 侧边栏
│   │   ├── DiskUsage.qml
│   │   ├── LeftPanel.qml
│   │   ├── PathCollectionItem.qml
│   │   ├── PathCollection.qml
│   │   └── User.qml
│   ├── LoadDirError.qml
│   ├── pan.qml
│   └── PathNav.qml
├── qml.qrc // xml文件, 组织qml用的
├── transfer
│   ├── DownloadingList.qml
│   ├── DownloadItem.qml
│   ├── DownloadList.qml
│   ├── HeaderBar.qml
│   └── transfer.qml
├── videoPlayer // 视频播放器
│   ├── ButtonImage.qml
│   ├── ControlArea.qml
│   ├── icons
│   │   ├── backward.svg
│   │   ├── forward.svg
│   │   ├── fullscreen.svg
│   │   ├── open-file.svg
│   │   ├── pause.svg
│   │   ├── play.svg
│   │   ├── rotate.svg
│   │   ├── volume-down.svg
│   │   ├── volume-mute.svg
│   │   └── volume-up.svg
│   ├── MPlayer.qml
│   ├── screen.png
│   ├── Tips.qml
│   └── UIComp
│       ├── ActionTips.qml
│       ├── AniIcon.qml
│       ├── ButtonPlay.qml
│       ├── DataSaver.qml
│       ├── ForwardBackward.qml
│       ├── FullScreenButton.qml
│       ├── LoadingTips.qml
│       ├── OpenFileButton.qml
│       ├── PlayIcon.qml
│       ├── PlayRateButton.qml
│       ├── RotateButton.qml
│       ├── TimeSlider.qml
│       ├── TimeText.qml
│       ├── VideoTitle.qml
│       ├── VolumeButton.qml
│       ├── VolumeIcon.qml
│       └── VolumeSlider.qml
└── widget // 其他组件
    ├── MoveWindow.qml // 鼠标移动窗口支持
    ├── Resize.qml // 鼠标改变窗口尺寸
    ├── RightMenu.qml // 右键菜单v2
    ├── ToolTip.qml
    ├── TopIndicator.qml // 顶部信息提示
    └── VirtualFrame.qml // 虚拟边框
```

### 服务端模块
```
server
├── artisan // 一些不暂时放这的功能 ;)
│   ├── cache // 缓存包
│   │   └── redis.go // redis实现缓存
│   ├── errors.go // 接口错误处理
│   ├── utils.go // 工具函数
│   └── web.go // web 应用特有
├── conf // 配置包
│   ├── conf.go // 配置
│   └── consts.go // 常量
├── dao // 数据库, mongobd
│   ├── feedback.go // 反馈表
│   ├── file-share.go // 文件分享表
│   ├── user.go // 用户表
│   ├── vip.go // 开通vip账号的表
│   └── vip-save-file.go // vip转存文件表
├── go.mod
├── go.sum
├── pan-light-server.go // 服务端入口
├── pan-light-server.yaml // 配置文件
├── pan-viper // vip账号包
│   ├── http-client.go // http 客户端
│   └── vip.go // vip结构体, 提供文件转存功能
├── pc-api // 给客户端调用的api
│   ├── handlers.go // api实现函数
│   ├── middleware // 中间件
│   │   └── pc-jwt.go // jwt auth
│   └── pc-api-router.go // 路由
├── realtime // 实时通信包, 基于websocket
│   ├── room.go // 房间结构体
│   ├── server.go // 实时通信服务器
│   └── session.go // 会话
└── timewheel // 时间轮算法实现
    └── timewheel.go
```

### qt 绑定模块

> 这个模块是基于 [therecipe/qt](https://github.com/therecipe/qt) 修改的, 可以访问原项目获取详情, 会用就行


 