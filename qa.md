# Q&A

一些常见的问题和解决方法，在提 issue 之前，请先查看一下本篇里面的内容，或许可以找到满意的答案。

## Q1

**Question:** 编译 agent 时报错: agent\vendor\github.com\akrennmair\gopcap\pcap.go:12:18: fatal error: pcap.h: No such file or directory.

**Answer:** 根据 [编译文档](https://github.com/ysrc/yulong-hids/blob/master/docs/build.md) 对 agent 进行编译，并报 `agent\vendor\github.com\akrennmair\gopcap\pcap.go:12:18: fatal error: pcap.h: No such file or directory.` 的情况, 有几种原因及解决方法:

1. 没有安装 libpcap 或者 winpcap 依赖
2. 由于 [gopacket的代码](https://github.com/google/gopacket/blob/master/pcap/pcap.go#L17) 写死了WpdPack/Include 包的位置，所以安装 winpcap 时不能修改该位置，请卸载并重新安装 libpcap/winpcap 到默认位置。
3. 部分 Windows 系统的安装后的默认位置并非 C:/WpdPack/Include, 下载 [WinPcap Developer's Pack](https://www.winpcap.org/devel.htm)， 并解压到C盘根目录，保证有该文件路径及路径下的依赖文件即可。

## Q2

**Question:** 运行web程序报: "panic: prefix should has path" 错误。

**Answer:** 这个错误来自 beego, 当beego找不到配置文件的时候会产生该错误， 配置路径为在 'web/conf/app.conf', 该路径下必须存在配置文件。 配置文件样本及示例可参见: [web/conf/app-config-sample.conf](https://github.com/ysrc/yulong-hids/blob/master/web/conf/app-config-sample.conf)。

## Q3

**Question:** Web前端显示错误，查看log看到以下信息： Collections pipe(pipe.All) aggregate all exception: invalid operator '$dateToString'。

**Answer:** 几乎所有的Web问题都可以通过看log解决，Web的log非常全，只要保证配置文件的 loglevel 为 6 以上，就可以看到很多日志信息。这个问题里的报错来自 Mongodb，当前的 Mongodb 不支持 $dateToString 操作符，请升级 Mongodb 到文档要求的版本。
