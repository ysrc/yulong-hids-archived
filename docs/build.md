# 编译指南
----------
## 环境要求

Golang 1.9  
Windows版本只需编译为32位，代码已做了兼容，可在64位系统中正常工作。

```
// 依赖
Go依赖包都集成在相应工程的vendor目录中
编译Agent需要先安装libpcap-devel
```

## 编译
### 客户端（Agent，Daemon、依赖）
```
# windows 32/64
cd %GOPATH%\src\
git clone github.com/ysrc/yulong-hids/

// 编译agent
go build -o bin/win-32/agent.exe --ldflags="-w -s" agent/agent.go
copy bin/win-32/agent.exe bin/win-64/agent.exe

// 编译daemon
go build -o bin/win-32/daemon.exe --ldflags="-w -s" daemon/daemon.go
copy bin/win-32/daemon.exe bin/win-64/daemon.exe
```

```
# linux 64
cd $GOPATH/src
git clone github.com/ysrc/yulong-hids/

// 编译agent
go build -o bin/linux-64/agent --ldflags="-w -s" agent/agent.go

// 编译daemon
go build -o bin/linux-64/daemon --ldflags="-w -s" daemon/daemon.go
```

> 编译后需压缩为不同系统的zip文件（agent.exe、daemon.exe、data.zip），在向导过程中上传。

### 服务端（Server、Web）
```
go build -o bin/server --ldflags="-w -s" server/server.go
go build -o bin/web --ldflags="-w -s" web/main.go
```


### 内核、驱动

下载地址：[wdk 7600](http://download.microsoft.com/download/4/A/2/4A25C7D5-EFBE-4182-B6A9-AE6850409A78/GRMWDK_EN_7600_1.ISO)
```
// win驱动
准备2008或win7 x64系统，安装 GRMWDK 7600。
进入对应系统的控制台，例:Windows Vista and Windows Server 2008\x64 Checked Build Environment。
cd %GOPATH%\src\yulong-hids\driver\
build

// Linux内核
安装对应内核版本的kernel-devel
cd $GOPATH\src\yulong-hids\syscall_hook\
make

```
> 内核和驱动替换到对应系统压缩文件的data.zip内。
> Windows 2008的驱动需自己购买证书进行签名，否则无法正常安装使用。