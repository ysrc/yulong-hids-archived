# 部署文档
----------
## 服务器准备

所有服务和后端均为跨平台解决方案，可为Windows或Linux或混合使用。

| 要求        |    版本     |   数量      |
| --------    | -----:  | -----:  |
| Mongodb    | 3.x  |1|
| Es集群     |  5.x  |   3+主机数/1000  |
| Server集群         | |主机数/1000    |
|WebServer||1|

## Es集群配置
以5000台服务器为例子，需准备好8台服务器并安装好jar [下载地址][1] 和elasticsearch [下载地址][2]。

 - client node 1台（10.100.100.100）
 - master node 2台（10.100.100.101-102）
 - data node 其余的5台（10.100.100.103-107）

> 根据实际情况调整节点角色和配置，如果只有几百台机器，单节点ES就足够。

**配置文件：**  
`client node`
```
cluster.name: yulonghids
node.name: client-node-1
node.data: false

# 监听的IP地址
network.host: 10.100.100.100
```
`master node`
```
cluster.name: yulonghids

# 数字累加区分
node.name: master-node-1
node.master: true
node.data: true

# 监听的IP地址
network.host: 10.100.100.101
discovery.zen.ping.unicast.hosts: ["10.100.100.100"]
```
`data node`
```
cluster.name: yulonghids

# 数字累加区分
node.name: data-node-1
node.data: true

# 监听的IP地址
network.host: 10.100.100.103
discovery.zen.ping.unicast.hosts: ["10.100.100.100"]
```

> Es服务器需**设置防火墙**只允许集群之间和Server以及Web服务器的连接（9200,9300端口）。


## MongoDB配置
安装运行实例即可无需其他操作。[下载地址][3]
> MongoDB服务器需**设置防火墙**只允许Server集群和Web服务器的连接。

## WebServer安装向导

 - 将web目录拷贝到WebServer服务器上
 - 修改配置文件里的管理密码（MD5）、TwoFactorAuthKey、mongodb链接、es client node地址
 - 控制台启动web.exe后通过浏览器访问进入向导过程，根据向导提示进行即可

> TwoFactorAuthKey使用Google Authenticator，需在手机中安装此APP并自行定义一个base32编码的密钥（如图），开启后敏感操作均需通过此APP的动态密码进行二次验证。

![](./auth.png)

## Server集群
将server.exe拷贝到各Server集群服务器上。  
运行参数：server.exe -db 10.100.100.200:27017 -es 10.100.100.100:9200

> Server会在33433端口开放RPC服务，请保持此端口与所有Agent机器通信畅通。  

## Agent安装
完成以上所有步骤且没有任何报错后即可开始安装Agent，建议小规模部署测试**确定稳定后**再批量覆盖。

```
# 在主机列表添加处可查看自动生成的安装命令
# 例WebServer为http://10.100.100.254:8080
# Windows安装命令
cd %SystemDrive% & certutil -urlcache -split -f http://10.100.100.254/json/download?type=daemon^&system=windows^&platform=64^&action=download daemon.exe & daemon.exe -netloc 10.100.100.254:443 -install
# 手动卸载
net stop yulong-hids & C:\yulong-hids\daemon.exe -uninstall

# Linux安装命令（依赖libpcap，未安装的需先编译安装libpcap）
wget -O /tmp/daemon http://10.100.100.254/json/download?type=daemon\&system=linux\&platform=64\&action=download;chmod +x /tmp/daemon;/tmp/daemon -install -netloc 10.100.100.254:443
# 手动卸载
service yulong-hids stop & /usr/yulong-hids/daemon -uninstall

```
> 此系统规则仅适合服务器，请不要安装在个人办公机上。  
> Daemon会开放监听tcp 65512端口用于接收任务，请保持此端口不被禁止。  
> Agent会本地监听udp 65530端口用于接收进程创建信息。  

## 开始使用

部署完成后默认为观察模式(建议1-3天)，在此模式下所有警报都不会提示仅在后台进行统计，由于此系统规则定义中有首次出现的概念，在关闭观察模式后会有大量的结果需要人工审核，处置完结果后就可以开始正式使用了，具体功能介绍可查看[使用帮助](./help.md)。





  [1]: https://sec.ly.com/mirror/jre-8u131-windows-x64.exe
  [2]: https://www.elastic.co/downloads/past-releases/elasticsearch-5-6-5
  [3]: https://sec.ly.com/mirror/mongodb-win32-x86_64-2008plus-ssl-3.4.0-signed.msi
