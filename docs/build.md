# 编译指南
----------
## 环境要求

Golang环境安装可以[参照](https://github.com/astaxie/build-web-application-with-golang/blob/master/zh/01.1.md)，各系统会与下面的示例稍有不同，Golang版本不可以低于1.9。

## 依赖

- Go依赖包都集成在相应工程的vendor目录中
- 编译Agent需要先安装libpcap-devel

**Windows 下编译 Agent 需要 [winpcap](https://www.winpcap.org/install/default.htm) 支持。且受到 [google/gopacket](https://github.com/google/gopacket) 影响可能会出现一些问题，具体请看 [Q&A#Q1](../qa.md#Q1)**
**Windows 下编译依赖gcc，可以通过mingw-w64 [32位](https://jaist.dl.sourceforge.net/project/mingw-w64/Toolchains%20targetting%20Win32/Personal%20Builds/mingw-builds/7.3.0/threads-posix/dwarf/i686-7.3.0-release-posix-dwarf-rt_v5-rev0.7z) [64位](https://jaist.dl.sourceforge.net/project/mingw-w64/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/7.3.0/threads-posix/seh/x86_64-7.3.0-release-posix-seh-rt_v5-rev0.7z)安装**


## 编译
### 客户端（Agent，Daemon、依赖）

- 下载安装对应安装包 https://golang.google.cn/dl/ 安装 golang 
- windows 下编译需注意：32位版本须在 32 位版本 win7/2003 下编译。Win版本 agent 代码做了兼容，如只编译32位亦可在64位系统中正常工作
- `go get github.com/ysrc/yulong-hids` 之后在 `{gopath}/src/github.com/ysrc/yulong-hids` 下的项目编译后是会出现预期外异常的，请根据接下来的指导文档进行编译。

```
cd C:\Go\src
git clone https://github.com/ysrc/yulong-hids/

// 编译agent
go build -o yulong-hids\bin\win-64\agent.exe --ldflags="-w -s" yulong-hids\agent\agent.go

// 编译daemon
go build -o yulong-hids\bin\win-64\daemon.exe --ldflags="-w -s" yulong-hids\daemon\daemon.go
```

```
# linux 64
// 下载并解压
wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz && tar -zxvf go1.10.linux-amd64.tar.gz -C /usr/local/

sudo vi /etc/profile 
// 并添加下面的内容：

export GOROOT=/usr/local/go
export GOBIN=$GOROOT/bin
export PATH=$PATH:$GOBIN
export GOPATH=$HOME/gopath (可选设置)

// 重新加载 profile 文件
source /etc/profile

cd /usr/local/go/src
git clone https://github.com/ysrc/yulong-hids/

// 编译agent
go build -o yulong-hids/bin/linux-64/agent --ldflags="-w -s" yulong-hids/agent/agent.go

// 编译daemon
go build -o yulong-hids/bin/linux-64/daemon --ldflags="-w -s" yulong-hids/daemon/daemon.go
```

> 编译后需压缩为不同系统的zip文件（如agent.exe、daemon.exe、data.zip），在向导过程中上传。  
为了方便大家，我们提供了一个编译打包脚本，在build目录下运行 python3 build.py 就能编译并打包。  
但是由于data.zip中提供的依赖文件不一定适配，可能需要你自行编译，所以这个文件需要手动打包进最终的部署包。

### 服务端（Server、Web）
```
go build -o yulong-hids/bin/server --ldflags="-w -s" yulong-hids/server/server.go
go build -o yulong-hids/web/web --ldflags="-w -s" yulong-hids/web/main.go
```


### 内核、驱动

win下的驭龙驱动文件pro.sys我们已经编译好了现成的放在了data.zip中，如果自行编译的话需要购买代码签名证书进行签名，否则无法正常加载使用。

下载地址：[wdk 7600](http://download.microsoft.com/download/4/A/2/4A25C7D5-EFBE-4182-B6A9-AE6850409A78/GRMWDK_EN_7600_1.ISO)

```
// win驱动
准备2008或win7 x64系统，安装 GRMWDK 7600。
从开始菜单中选择进入对应系统的控制台，如
Windows Driver Kits->WDK 7600.16385.1->Build Environments->Windows Vista and Windows Server 2008->x64 Checked
Build Environment
具体根据你的系统选择Build Environments，要编译32位的就选x86。

cd C:\Go\src\yulong-hids\driver\
build

// Linux内核
我们在bin\linux-64\data.zip中已经提供了一些编译好的对应内核版本的ko文件，直接用的话需要确认内核版本完全一致。
实际部署过程中需要 uname -r 统计下需要部署的机器linux内核版本，然后需要找到完全匹配对应版本的kernel-devel包
下下来并安装，yum安装的不一定完全匹配。虽然软链接到不一致的版本也能编译出来ko并加载，但是不知道会不会有影响，
稳妥起见还是参照上面的做法。
rpm -ivh kernel-devel-3.10.0-327.el7.x86_64.rpm

debian系的也类似 用 uname -r 的结果替换 4.9.0-3-amd64 
apt-get install linux-headers-4.9.0-3-amd64

openSUSE
zypper install kernel-devel

安装好对应kernel-devel包后cd进目录就可以编译了
cd /usr/local/go/src/yulong-hids/syscall_hook
make

```
编译好的内核和驱动需要替换进对应系统目录下的data.zip内。

驱动签名的我们给两个示例，具体命令会因证书商不一样而稍有不同。

signtool sign /v /ac "MSCV-ThawteClass3.cer" /a /s MY /n "Tongcheng Network Technology Co., Ltd" /fd sha256 /tr http://sha256timestamp.ws.symantec.com/sha256/timestamp/ "e:\pro.sys"

signtool sign /v /ac "DigiCert High Assurance EV Root CA.crt" /tr http://timestamp.digicert.com /td sha256 /fd sha256 /f "e:\sha256new.p12" /p 证书密码 "e:\pro.sys"
