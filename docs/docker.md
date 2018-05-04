## Docker 体验版安装

### 说明

**Docker 版只为快速体验使用，请不要在生产环境下使用!**

该版本中只包括了服务端：`Server`、`Web`、`ElasticSearch`和`MongoDB`, Client 端的 `Agent` 和 `Daemon` 请在相应的机器中运行即可.

### 依赖

- docker-ce >= 18
- docker-compose >= 1.20

> 安装 Docker 和 Docker-Compose 步骤请参考 Docker 官方文档和搜索引擎(非常简单)

### 使用步骤

#### Step1. 下载源码

```
$ git clone https://github.com/ysrc/yulong-hids.git
```

#### Step2. 初次编译并运行

```
$ cd yulong-hids/
$ docker-compose up
```

> 由于需要映射 Web 80/443 端口到宿主机80/443端口，请保证有权限，如果提示 Permission denied, 请执行 `sudo docker-compose up`


> 第一次启动时由于 Server 需要配置文件不存在会导致启动失败，不要慌，只要保证 web、mongo、es 正常启动即可

> 如果提示获取 docker image 失败，请检查网络并重试，或者给 docker pull 挂上 http/https 代理

#### Step3. 通过 Web 界面初始化

假定宿主机（物理机）的 IP 地址是: 192.168.1.101

打开浏览器访问 `http://192.168.1.101` 如果启动正常，就可以看到驭龙的Login界面了，输入下面的登录名和密码进入后台。

登录名 | 密码 | 二次验证秘钥
:-: | :-: | :-:
`yulong` | `All_life_is_a_game_of_luck.` | <img src="./docker_totp_default.png" width="128px"> `IVFHGS2OGYTXIVDGEIZWCNC2MVMHYWDRK44GOQALPNJHGRS6FE2QUCT4`

值得一提的是，初始化的第3步所需要上传的文件需前往[Release发布页](https://github.com/ysrc/yulong-hids/releases) 下载发行版 zip 包，并解压，然后上传对应的 `win-32.zip`,`win-64.zip`,`linux-64.zip` 即可。

#### Step4. 重新启动

通过 Web 初始化完毕后，切回 `docker-compose up` 的终端，按下 `Ctrl + C` 组合键，结束进程，然后执行：

```
$ docker-compose up -d
```

如果看到如下提示，证明启动完毕

```
$ docker-compose up -d
Creating yulonghids_ids_elasticsearch_1 ... done
Creating ids_mongodb                    ... done
Creating ids_web                        ... done
Creating ids_server                     ... done
```

浏览器打开 http://192.168.1.101 就可以看到正常功能的界面了.

#### Step5. 启动Server调校程序

**这一步是Docker版本特有，其它方式不需要**

> 由于 Docker 使用的 network 是 bridge 方式，server 的 getLocalIP 方法中只能获取到 container 内部的 IP, agent 中通过 getserverlist 拿到的 IP 地址是不对的，所以需要额外走这一步来调整 Server IP 地址

1.根据**宿主机(物理机)**的操作系统，下载对应的程序:

Adapter | 可执行文件 MD5
:-- | :--
[Windows_amd64](https://sec.ly.com/mirror/yulong_server_docker_adapter_win_x64.zip)| `6cd550b9443be8f6c19dc58b460ef877`
[Linux_amd64](https://sec.ly.com/mirror/yulong_server_docker_adapter_linux_amd64.zip) | `c2775b65c24b479c138e14338fc893d1`
[Mac_amd64](https://sec.ly.com/mirror/yulong_server_docker_adapter_darwin_amd64.zip) | `d7ca534fb0b553dedd6d8327d4e329bb`
 
 其它操作系统的自行下载源码编译，可执行文件源码地址：[yulong_server_docker_adapter.go](https://gist.github.com/Medicean/8ffd6babb42b8078064d8306069aee79)

 > Docker 不支持 32 位操作系统，所以没有32位的版本

2.在**宿主机(物理机)**上面执行该程序

以 Linux 为例:

```
$ ./yulong_server_docker_adapter_linux_amd64
```

> 如果宿主机无法联网，那么需要手动指定宿主机的IP地址，使用 `-ip` 参数即可.
> 
> 嫌文件名长的话自己改名 = =

 这个程序需要一直运行，所以可以使用 nohup 或者 screen 来执行，具体方法就不在赘述了.

### 其它

Agent 连接 Server 请直接参考真机布署文档即可

### 异常情况解决

**Q1**: 启动之后，Web 或者 Server 启动不成功？

**A1**: 由于这两者启动时需要 mongo 和 es 已经启动成功，不同的机器启动时的速度不一致导致会有启动不成功的情况，遇到这种情况，请手动启动 Web 和 Server 即可: `docker start ids_web ids_server`

**Q2**: 如何停止所有服务?

**A2**: 在 yulong-hids 目录下执行: `docker-compose stop`

**Q3**: 怎么修改 web 的配置文件?

**A3**: 需要进入容器中进行修改: `$ docker exec -it ids_web /bin/sh` 然后用 vi 去编辑


**Q4**: 编译出现 parsing reference 错误怎么办？具体如下

```
Building ids_web
Step 1/20 : FROM golang:1.10 as builder
ERROR: Service 'ids_web' failed to build: Error parsing reference: "golang:1.10 as builder" is not a valid repository/tag: invalid reference format
```
**A4**: 是因为你的 docker-ce 版本太低导致，升级至 docker-ce >= 18, 具体看本文档最上方 **依赖** 项
