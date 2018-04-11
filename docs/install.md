# 部署文档
----------
## 服务器资源需求参照表

所有程序均为跨平台解决方案，可 Windows 或 Linux 混合使用。

| 要求        |    版本要求   |   数量      |
| :------:    | :----:  | :----:  |
| MongoDB  | 3.x  |1|
| Elasticsearch 集群 |  5.x  |  3 + 主机数/1000  |
| server 集群       | |主机数/1000    |
|web||1|

## 部署流程

部署 MongoDB (3.x，驭龙不兼容2.x版本)；  
部署 Elasticsearch (5.x，驭龙暂不兼容6.x版本)；  
启动 MongoDB 、Elasticsearch；  
修改 web 的配置，启动 web ，在引导界面根据提示初始化数据库、规则等；  
启动 server（服务端）；  
部署 daemon （守护进程），启动 agent（客户端）。

### 部署 MongoDB (3.x，驭龙不兼容2.x版本)；

下载对应安装包 https://www.mongodb.com/download-center?jmp=nav#community

以 centos 上部署为例，如果报错，请去掉 --fork 查看原因。

```
mkdir /var/lib/mongodb/ && mkdir /var/log/mongodb && wget https://sec.ly.com/mirror/mongodb-linux-x86_64-3.6.3.tgz && tar -xvzf mongodb-linux-x86_64-3.6.3.tgz && mongodb-linux-x86_64-3.6.3/bin/mongod --dbpath /var/lib/mongodb/ --logpath /var/log/mongodb.log --fork --bind_ip 10.0.0.134
```

从 MongoDB 3.6版本开始，出于安全考虑，如果不指定实例绑定ip，默认是 bind 到 localhost  (127.0.0.1)的。

**但是驭龙 MongoDB 这边必需 bind_ip 指定非 localhost (127.0.0.1)的本机ip，否则后面会出错。**

> MongoDB 服务器需配置防火墙策略只允许 Server 集群和 Web 服务器的连接。

### 部署 Elasticsearch (5.x，驭龙暂不兼容6.x版本)

下载安装[jre](https://www.java.com/zh_CN/download/manual.jsp)依赖， 因为官网下载较慢，这边缓存了一份。

```
wget https://sec.ly.com/mirror/jre-8u161-linux-x64.rpm && yum -y localinstall jre-8u161-linux-x64.rpm
```

下载ES并解压

```
wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-5.6.8.tar.gz && tar -zxvf elasticsearch-5.6.8.tar.gz -C /opt
```

Elasticsearch 不建议以 root 权限运行，新建一个非 root 权限用户，-p 后跟自行设定的密码

```
groupadd elasticsearch && useradd elasticsearch -g elasticsearch -p ElasticSearch666
```

修改文件夹及内部文件的所属用户及组为 elasticsearch:elasticsearch

```
chown -R elasticsearch:elasticsearch /opt/elasticsearch-5.6.8
```

centos7 以下的系统需编辑 config/elasticsearch.yml 添加 

```
bootstrap.system_call_filter: false
```

启动es

```
su - elasticsearch -c '/opt/elasticsearch-5.6.8/bin/elasticsearch -d'
```

非单机测试部署可以修改 network.host: 后面的ip，监听对应ip。

curl请求下确认ES启动成功

```
curl -XGET -s "http://localhost:9200/_cluster/health?pretty"
```

#### ES集群部署

如果只有几百台机器，单节点ES就足够；部署实例较多的话，ES就需要集群部署了。

以5000台服务器为例子，需准备好8台服务器左右。如果单ES实例配置较高，数量可以减少，可根据服务器资源情况调整节点角色和配置。

 - client node 1台（10.100.100.100）
 - master node 2台（10.100.100.101-102）
 - data node 其余的5台（10.100.100.103-107）

**配置文件：**  

##### client node

```
cluster.name: yulonghids
node.name: client-node-1
node.data: false

# 监听的IP地址
network.host: 10.100.100.100
```
##### master node

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
##### data node

```
cluster.name: yulonghids

# 数字累加区分
node.name: data-node-1
node.data: true

# 监听的IP地址
network.host: 10.100.100.103
discovery.zen.ping.unicast.hosts: ["10.100.100.100"]
```

> ES服务器需配置防火墙策略只允许集群之间和Server以及Web服务器的连接（9200,9300端口）。

### web 配置

 - 将 web 目录拷贝到 WebServer 服务器上，单台机器测试也可以都在一台机子上。

 - 修改 web 的配置，必须改名为 app.conf

   `mv yulong-hids/web/conf/app-config-sample.conf yulong-hids/web/conf/app.conf`
   `vi yulong-hids/web/conf/app.conf`

   主要是改3个地方

   管理密码 passwordhex 是密码的32位MD5值，可以 echo -n password | md5sum 或者去cmd5生成一个替换掉；

   TwoFactorAuthKey 是开启二次验证后，敏感操作都需要Google Authenticator生成的动态口令做二次验证，请确保服务器跟手机的时间都正确；

   mongodb ip:port 修改为 MongoDB 之前 bind 的 ip:27017，ES修改为ES实例的 ip:9200，ip不对会导致web面板报错；

   如果需要 web 运行在其他端口，还需要修改对应的 HTTPPort 和 HTTPSPort。

#### 启动 web

可以直接用 YSRC 编译好的[版本](https://github.com/ysrc/yulong-hids/releases),也可以参照[编译指南](./build.md)自行编译。

cd 进 web 目录
`cd yulong-hids/web/`

`./web`  启动 web，如果是下的编译好的二进制需要赋予执行权限 `chmod +x web/web`

如果能正常访问，可以放到后台去运行。

`nohup ./web &`

win 版本控制台运行 web.exe 后通过浏览器访问进入向导过程，根据向导提示进行即可。

> TwoFactorAuthKey 使用 Google Authenticator，需在手机中安装并导入生成的base32编码的密钥（如图），具体生成方式见 app.conf 注释，开启后敏感操作均需通过生成的动态口令进行二次验证。

![](./auth.png)

## server 集群
将 server 二进制文件拷贝到各 server 集群服务器上，ip不对会导致报错。
`server -db 10.0.0.134:27017 -es 10.100.100.100:9200`

如果能正常访问，可以nohup放到后台去运行。db 跟MongoDB地址端口，es 跟 Elasticsearch 地址端口。

> server会在33433端口开放RPC服务，请保持此端口与所有Agent机器通信畅通。  

## agent 部署
完成以上所有步骤且没有任何报错后即可开始安装 agent，建议小规模部署测试**确定稳定后**再灰度部署。

```
# 在主机列表添加处可查看自动生成的安装命令
# 例 web 地址为为http://10.100.100.254，netloc 后跟的ip即为 web 的ip
# Windows 安装命令
cd %SystemDrive% & certutil -urlcache -split -f http://10.100.100.254/json/download?type=daemon^&system=windows^&platform=64^&action=download daemon.exe & daemon.exe -netloc 10.100.100.254:443 -install

# 手动卸载
net stop yulong-hids & C:\yulong-hids\daemon.exe -uninstall

# Linux 安装命令（依赖libpcap，未安装的需先安装libpcap）
wget -O /tmp/daemon http://10.100.100.254/json/download?type=daemon\&system=linux\&platform=64\&action=download;chmod +x /tmp/daemon;/tmp/daemon -install -netloc 10.100.100.254:443

# 手动卸载
service yulong-hids stop & /usr/yulong-hids/daemon -uninstall

#如果看不到agent上线，参见下面的命令调试，ip跟web的ip。
一般来说报错信息都比较明显，server/MongoDB/ES没起，MongoDB/ES连不上之类的。
agent 10.100.100.254 debug
```
> 目前驭龙系统的设计仅适合服务器场景，不适合部署在线下办公环境 ;
> daemon 会开放监听 tcp 65512 端口用于接收任务，请保持此端口不被禁止;  
> agent 会本地监听 udp 65530 端口用于接收进程创建信息。  

## 开始使用

部署完成后默认为观察模式(建议1-3天)，在此模式下所有警报都不会提示仅在后台进行统计。

由于驭龙系统规则定义中有首次出现的概念，在关闭观察模式后会有大量的结果需要人工审核，处置完结果后就可以开始正式使用了，具体功能介绍可查看[使用帮助](./help.md)。
