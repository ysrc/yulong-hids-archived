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



安装过程中会尽可能使用国内镜像源以加快安装速度。

### 部署 MongoDB (3.x，驭龙不兼容2.x版本)

#### Debian/Ubuntu 用户

* 首先信任 MongoDB 的公钥：

```shell
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv EA312927
```

* 然后根据自己的 Debian/Ubuntu 版本，将相应的内容写入 `/etc/apt/sources.list.d/mongodb.list` 中，如果此文件不存在，创建之：

```
# Ubuntu 14.04 LTS
deb https://mirrors.tuna.tsinghua.edu.cn/mongodb/apt/ubuntu trusty/mongodb-org/stable multiverse
```

```
# Ubuntu 16.04 LTS
deb https://mirrors.tuna.tsinghua.edu.cn/mongodb/apt/ubuntu xenial/mongodb-org/stable multiverse
```

```
# Ubuntu 18.04 LTS
deb https://mirrors.tuna.tsinghua.edu.cn/mongodb/apt/ubuntu bionic/mongodb-org/stable multiverse
```

```shell
# Debian 7
deb http://mirrors.tuna.tsinghua.edu.cn/mongodb/apt/debian wheezy/mongodb-org/stable main
```

* 然后就可以安装相应版本的 mongodb 了：

```shell
$ sudo apt update
$ sudo apt install mongodb-org
```

默认安装的版本为 `3.2.21`。

#### RHEL/CentOS 用户

* 首先新建文件 `/etc/yum.repos.d/mongodb.repo`，并写入以下内容：

```shell
[mongodb-org]
name=MongoDB Repository
baseurl=https://mirrors.tuna.tsinghua.edu.cn/mongodb/yum/el$releasever/
gpgcheck=0
enabled=1
```

* 然后刷新包缓存并安装即可：

```shell
$ sudo yum makecache
$ sudo yum install mongodb-org
```

默认安装的版本为 `3.2.21`。



然后修改配置文件 `/etc/mongo.conf` 中的 `bindIp` 字段，将 `MongoDb` 绑定在本机的局域网 IP 上。

**驭龙 MongoDB 这边必需 bindIp 指定非 localhost （127.0.0.1）的本机 IP，否则后面会出错。**

> MongoDB 服务器需配置防火墙策略只允许 Server 集群和 Web 服务器的连接。

最后启动 MongoDB。在支持 `systemd` 的操作系统上（比如 CentOS 7、Ubuntu 16.04 及以上）使用 `systemctl` 启动，在不支持 `systemd` 的操作系统上使用 `service` 命令启动：

```shell
# CentOS7, Ubuntu 16.04, Ubuntu 18.04
$ sudo systemctl start mongod.service
# Ubuntu 14.04
$ sudo service mongod start
```



### 部署 Elasticsearch（5.x，驭龙暂不兼容 6.x 版本）

* 因为 Elasticsearch 使用 Java 开发，所以需要安装 Java 8 或者更高。通过包管理器可直接安装。CentOS 7、Ubuntu 16.04 及以上的官方仓库中 openjdk 的版本已经为 Java 8，所以可以直接安装，对应的包名为 `java-1.8.0-openjdk-headless`（CentOS 7）和 `openjdk-8-jre-headless`。Ubuntu 14.04 的官方仓库源中为 Java 7，需要添加 PPA 源进行安装：

```shell
# For Ubuntu 14.04
$ sudo add-apt-repository ppa:openjdk-r/ppa
$ sudo apt update
$ sudo apt install openjdk-8-jre-headless
```

* 导入 Elasticsearch 的签名公钥：

```shell
# RHEL/CentOS 用户
rpm --import https://artifacts.elastic.co/GPG-KEY-elasticsearch
# Debian/Ubuntu 用户
wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo apt-key add -
```

* 然后创建相应的仓库信息文件：

  * 对于 RHEL/CentOS 用户，创建文件 `/etc/yum.repos.d/elasticsearch.repo` 并加入以下内容：

  ```shell
  [elasticsearch-5.x]
  name=Elasticsearch repository for 5.x packages
  baseurl=https://mirrors.tuna.tsinghua.edu.cn/elasticstack/5.x/yum
  gpgcheck=1
  gpgkey=https://artifacts.elastic.co/GPG-KEY-elasticsearch
  enabled=1
  autorefresh=1
  type=rpm-md
  ```

  * 对于 Debian/Ubuntu 用户，创建文件 `/etc/apt/sources.list.d/elastic-5.x.list` 并加入以下内容：

  ```shell
  deb https://mirrors.tuna.tsinghua.edu.cn/elasticstack/5.x/apt stable main
  ```

* 然后更新软件仓库，安装：

```shell
# RHEL/CentOS 用户
$ sudo yum makecache
$ sudo yum install elasticsearch

# Debian/Ubuntu 用户
$ sudo apt update
$ sudo apt install elasticsearch
```

默认安装的版本为 `5.6.13`。配置文件目录路径为 `/etc/elasticsearch`。

* 修改配置文件 `/etc/elasticsearch/elasticsearch.yml` 加入以下内容：

```
bootstrap.system_call_filter: false
```

* Elasticsearch 的 JVM 配置文件默认指定的虚拟机内存大小为 2G，所以如果你的机器内存不够大，需要调整 JVM 参数。修改配置文件 `/etc/elasticsearch/jvm.options` 的以下两项以调整：

```
-Xms512m
-Xmx512m
```

* 最后启动 elasticsearch：

```shell
# systemd 用户
sudo systemctl start elasticsearch.service
# SysV init 用户
sudo service elasticsearch start
```

通过命令 `curl -XGET -s "http://localhost:9200/_cluster/health?pretty"` 可判断是否成功启动。

#### ES 集群部署

如果只有几百台机器，单节点ES就足够；部署实例较多的话，ES就需要集群部署了。

以5000台服务器为例子，需准备好8台服务器左右。如果单ES实例配置较高，数量可以减少，可根据服务器资源情况调整节点角色和配置。

- client node 1台（10.100.100.100）
- master node 2台（10.100.100.101-102）
- data node 其余的5台（10.100.100.103-107）

分别修改配置文件 `/etc/elasticsearch/elasticsearch.yml`：

* client node

```
cluster.name: yulonghids
node.name: client-node-1
node.data: false

# 监听的IP地址
network.host: 10.100.100.100
```

* master node

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

* data node

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

 - 修改 web 的配置：必须改名为 app.conf：

   * 首先将 `app-config-sample.conf` 重命名为 `app.conf`，接着修改 `app.conf`：

   ```shell
   $ mv yulong-hids/web/conf/app-config-sample.conf yulong-hids/web/conf/app.conf
   ```

   * 管理密码 passwordhex 是密码的 32 位 MD5 值，可以 `echo -n password | md5sum` 或者去 cmd5 生成一个替换掉；

   * TwoFactorAuthKey 是开启二次验证后，敏感操作都需要 Google Authenticator 生成的动态口令做二次验证，请确保服务器跟手机的时间都正确；

   * 将 mongodb 部分的 ip 和 port 修改为 MongoDB 配置文件中设置的 ip 和 port；
   * ES 修改为 ES 实例对应的 ip:9200，ip 不对会导致 web 面板报错；

   * 如果需要 web 运行在其他端口，还需要修改对应的 HTTPPort 和 HTTPSPort。

#### 启动 web

可以直接用 YSRC 编译好的[版本](https://github.com/ysrc/yulong-hids/releases),也可以参照[编译指南](./build.md)自行编译。

* 进入 web 目录然后执行 `./web` 即可启动。注意启动前要保证 `web` 可执行程序具有执行权限。

如果能正常访问，可以放到后台去运行：

```shell
$ nohup ./web &
```

* win 版本控制台运行 web.exe 后通过浏览器访问进入向导过程，根据向导提示进行即可。

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
