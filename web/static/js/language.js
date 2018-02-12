global_data_info = {
    "crontab": {
        "name": "计划任务名",
        "command": "要执行的程序或命令以及参数",
        "arg": "启动参数",
        "user": "启动用户",
        "rule": "",
        "description": "描述"
    },
    "listening": {
        "proto": "类型",
        "address": "监听地址",
        "name": "监听程序名",
        "pid": "监听程序pid"
    },
    "service": {
        "name": "服务名",
        "pathname": "启动命令，同command",
        "started": "当前启动状态",
        "startmode": "开机启动模式",
        "startname": "启动用户",
        "caption": "描述"
    },
    "startup": {
        "name": "名称",
        "command": "启动程序或命令",
        "location": "来源",
        "user": "启动用户"
    },
    "userlist": {
        "name": "用户名",
        "description": "描述",
        "status": "状态"
    },
    "file": {
        "path": "文件或者目录路径 file",
        "action": "行为类型",
        "user": "操作用户",
        "hash": "文件md5 hash"
    },
    "loginlog": {
        "username": "用户名",
        "hostname": "远程主机名",
        "remote": "远程IP",
        "status": "认证结果",
        "time": "时间"
    },
    "process": {
        "name": "进程名",
        "command": "程序或命令以及参数",
        "pid": "进程pid",
        "ppid": "父进程pid",
        "parentname": "父进程名"
    },
    "connection": {
        "dir": "方向 ",
        "protocol": "类型（TCP、UDP） type",
        "local": "本机进行通讯ip:port ",
        "remote": "远程进行通讯的ip:port",
        "name": "进程名",
        "pid ": "进程pid"
    }
}

// 加载部分css配置信息
HostWData.update("style", {
    "notice": {
        "level": {
            2: "badge badge-info",
            1: "badge badge-warning",
            0: "badge badge-danger",
        },
        "status": {
            1: "badge badge-success",
            0: "badge badge-danger",
            2: "badge badge-info",
        }
    }
});

// config chinese description
HostWData.update("zh_cn", {
    "config" : {
        "client": {
            "type_description": "客户端 Agent配置",
            "cycle": "间隔 收集型信息回传间隔",
            "lan": "内网连接 是否记录内网网络连接信息",
            "mode": "模式 （规划中）",
            "monitorPath": "监控目录 文件操作监控目录，%web%为自动识别的web目录，*结尾为迭代监控（例如/tmp/*）",
            "udp": "记录UDP 是否记录UDP连接信息"
        },
        "server": {
            "type_description": "服务端 Server配置",
            "cert": "证书",
            "learn": "观察模式 开启观察模式后所有的警报都只做记录统计方便判断是否为误报，在关闭时可进行汇总处理。（部署后默认为观察模式）",
            "privatekey": "私钥 任务指令加密传输RSA私钥（Server用）",
            "offlinecheck": "离线告警 周期循环检测主机是否离线",
            "publickey": "公钥 任务指令加密传输RSA公钥（Daemon用）"
        },
        "intelligence": {
            "type_description": "威胁情报",
            "ipapi": "IP检测接口 IP威胁情报接口，格式为：http://x.x.x.x/api/check_ip/?ip={$ip}，{$ip}为IP的占位符",
            "fileapi": "文件检测接口 文件威胁情报接口，格式为：http://x.x.x.x/api/check_file/?hash={$hash}，{$hash}为文件md5值的占位符",
            "regex": "正则匹配 接口返回结果判断正则，如果匹配会产生警报信息",
            "switch": "开关"
        },
        "blacklist": {
            "type_description": "黑名单",
            "file": "文件 文件行为，可文件md5或文件路径的正则(自动识别)",
            "ip": "IP IP地址，不包含端口",
            "process": "进程 进程名称或参数的正则",
            "other": "其他 其他类型信息的正则（自动识别对应的关键字段）"
        },
        "whitelist": {
            "type_description": "白名单",
            "file": "文件 文件行为 可文件md5或文件路径的正则(自动识别)",
            "ip": "IP IP地址 不包含端口",
            "process": "进程 进程名称或参数的正则",
            "other": "其他 其他类型信息的正则（自动识别对应的关键字段）"
        },
        "filter": {
            "type_description": "过滤 （不传回Server记录，直接抛弃）",
            "file": "文件 文件行为 可文件md5或文件路径的正则(自动识别)",
            "ip": "IP IP地址 不包含端口",
            "process": "进程 进程名称或参数的正则"
        },
        "notice": {
            "type_description": "通知",
            "api": "通知接口 （例如短信、微信、邮件），格式为：http://x.x.x.x/sendmsg/?text={$info}，{$info}为消息通知占位符",
            "onlyhigh": "仅危险警告 仅对危险等级的告警进行通知",
            "switch": "开关"
        },
        "update": {
            "type_description": "更新Agent"
        }
    }
})

// 加载中文语言信息
HostWData.update("zh_cn", {
    "type": {
        "loginlog": "登录记录",
        "listening": "监听端口",
        "service": "服务",
        "userlist": "用户",
        "startup": "开机启动",
        "connection": "连接请求",
        "crontab": "计划任务",
        "process": "进程",
        "abnormal": "异常",
        "file" : "文件操作"
    }
})

HostWData.update("zh_cn", {
    "notice": {
        "key": {
            "info": "告警信息",
            "description": "描述",
            "status": "状态",
            "time": "时间",
            "type": "类型",
            "ip": "发出告警的主机IP",
            "source": "告警原因",
            "level": "告警等级",
            "raw": "原始数据"
        },
        "data": {
            "type": HostWData.zh_cn.type,
            "status": {
                0: "未处理",
                1: "已处理",
                2: "已忽略",
            },
            "level": {
                2: "提示",
                1: "可疑",
                0: "危险",
            }
        }
    },
    "file": {
        "fileupload": "上传文件",
        "upload": "上传"
    }
})

// format
format_connection_msg = function(source) {
    msg = "<span class='title'>网络连接信息:</span>";
    if (source.data.dir == "out") {
        msg += "出方向|";
    } else {
        msg += "入方向|";
    }
    if(source.data.remote)
        msg += "<span class='key'>远程地址:</span>" + source.data.remote + "  ";
    if(source.data.local)
        msg += "<span class='key'>本地地址:</span>" + source.data.local + "  ";
    if(source.data.name)
        msg += "<span class='key'>进程名:</span>" + source.data.name + "  ";
    if(source.data.pid)
        msg += "<span class='key'>PID:</span>" + source.data.pid + "  ";
    if(source.data.protocol)
        msg += "<span class='key'>网络类型:</span>" + source.data.protocol + "  ";
    msg += "<span class='key'>时间:</span>" + timeformat(source.time);
    return msg
}

format_file_msg = function(source) {
    msg = "<span class='title'>文件信息:</span>";
    if(source.data.action)
        msg += "<span class='key'>操作类型:</span>" + source.data.action + "  ";
    if(source.data.path)
        msg += "<span class='key'>文件路径:</span>" + source.data.path + "  ";
    if(source.data.hash)
        msg += "<span class='key'>文件hash:</span>" + source.data.hash + "  ";
    if(source.data.user)
        msg += "<span class='key'>用户:</span>" + source.data.user + "  ";
    msg += "<span class='key'>时间:</span>" + timeformat(source.time);
    return msg
}

format_process_msg = function(source) {
    msg = "<span class='title'>进程信息:</span>";
    if(source.data.command)
        msg += "<span class='key'>命令:</span>" + source.data.command + "  ";
    if(source.data.name)
        msg += "<span class='key'>进程信息:</span>" + source.data.name + ":" + source.data.pid + "  ";
    if(source.data.parentname)
        msg += "<span class='key'>父进程信息:</span>" + source.data.parentname + ":" + source.data.ppid + "  ";
    msg += "<span class='key'>时间:</span>" + timeformat(source.time);
    return msg
}

format_loginlog_msg = function(source) {
    msg = "<span class='title'>登录信息:</span>";
    if(source.data.hostname)
        msg += "<span class='key'>主机名:</span>" + source.data.hostname + "  ";
    if(source.data.remote)
        msg += "<span class='key'>remote:</span>" + source.data.remote + "  ";
    if(source.data.username)
        msg += "<span class='key'>用户名:</span>" + source.data.username + "  ";
    if(source.data.status)
        msg += "<span class='key'>登录状态:</span>" + source.data.status + "  ";
    msg += "<span class='key'>时间:</span>" + timeformat(source.time);
    return msg
}

format_msg = function(data) {
    source_type = data["_type"];
    if (source_type) {
        return window["format_"+source_type+"_msg"](data._source)
    }
}
