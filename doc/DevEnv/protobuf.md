# protobuf

> protobuf 是google出品的一个跨语言数据序列化库.
 它由2个部分组成:一个是protoc程序, 用来将数据格式描述文件(proto文件)转换成需要的编程语言代码.
 另一部分就是编程语言运行的支持库, 比如golang中的 github.com/golang/protobuf/proto 包.
 这里我们先下载安装protoc程序
 
## 安装protoc

```bash
# 获取源码包
wget https://github.com/google/protobuf/archive/v3.5.0.tar.gz

# 解压缩并进入源码目录
tar -zxvf v3.5.0.tar.gz
cd protobuf-3.5.0

# 生成configure文件
./autogen.sh

# 编译安装
./configure
make
make check
make install
```

## 相关资料
- [golang-protobuf快速上手指南](https://studygolang.com/articles/14337)

