# Q&A

一些常见的问题和解决方法，在提 issue 之前，请先查看一下本篇里面的内容，或许可以找到满意的答案。

## Q1

**Question:** 编译 agent 时报错: agent\vendor\github.com\akrennmair\gopcap\pcap.go:12:18: fatal error: pcap.h: No such file or directory.

**Answer:** 根据 [编译文档](https://github.com/ysrc/yulong-hids/blob/master/docs/build.md) 对 agent 进行编译，并报 `agent\vendor\github.com\akrennmair\gopcap\pcap.go:12:18: fatal error: pcap.h: No such file or directory.` 的情况, 有几种原因及解决方法:

1. 没有安装 libpcap 或者 winpcap 依赖
2. 由于 [gopacket的代码](https://github.com/google/gopacket/blob/master/pcap/pcap.go#L17) 写死了WpdPack/Include 包的位置，所以安装 winpcap 时不能修改该位置，请卸载并重新安装 libpcap/winpcap 到默认位置。
3. 部分 Windows 系统的安装后的默认位置并非 C:/WpdPack/Include, 下载 [WinPcap Developer's Pack](https://www.winpcap.org/devel.htm)， 并解压到C盘根目录，保证有该文件路径及路劲下的依赖文件即可。
