# qt5

## 下载

> Qt 是由c++构建的夸平台gui框架,  本项目使用的Qt \>= 5.11 版本.
 本项目gui部分基于qt quick, 也就是qml语言实现. 如果你不懂, 没关系, 作者也是第一次接触, 
 很简单, 是JavaScript语法. 如果你会js, 入门只需半天

下载地址 http://download.qt.io/archive/qt/5.12/5.12.4/
也可以通过系统的包管理器安装

## 配置环境变量

将如下内容加入`~/.bashrc`

```bash
export QT_VERSION=5.11.3 # 你安装的qt的版本
export QT_DIR=/media/peterq/files/dev/env/qt # 你安装的qt的路径
```

这2个环境变量是 github.com/peterq/pan-light/qt 这个模块用来创建qt 的 c++ 到 golang 的中间代码的

## 相关资料
- [QML入门教程](https://blog.csdn.net/qq_40194498/article/details/79849807)
