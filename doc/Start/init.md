# 项目初始化

## 克隆项目到本地

> 本项目使用 go mod, 不需要把项目放go path里

```bash
git clone git@github.com:peterq/pan-light.git

cd pan-light
```

## protobuf 生成文件

```bash
go generate

cd pc
```

## 用go封装c++类

```bash
go run ../qt/cmd/qtsetup/main.go
```

> 这里会把qt的模块生成对应的go包, 在开发时就会有相应的代码提示了.
实际上本项目只用到了少数几个模块. 这一步会耗时很久,但是这一步以后通常不需要在执行

## moc 原生组件

> gui/comp下有个go写的qml原生组件, 需要用c++封装才能被qt使用. 执行如下命令即可

```bash
go run ../qt/cmd/qtmoc/main.go desktop gui/comp
```

## 打包qml资源文件

```bash
go run ../qt/cmd/qtrcc/main.go desktop gui/qml
```

> 这一步会把qml资源文件打包成c++代码, 从而嵌入到程序里边

## 运行程序

```bash
go run pan-light-pc.go
```

