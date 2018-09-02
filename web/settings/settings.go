package settings

import "gopkg.in/mgo.v2/bson"

var (
	// Version 版本号
	Version = "v0.4.4 BETA"

	// TFAPassHistorys 双因子验证密码历史, 永远只接受6个历史
	TFAPassHistorys = []uint32{0, 0, 0, 0, 0, 0}

	// SystemArray System参数白名单
	SystemArray = []string{"linux", "windows"}

	// PlatformArray Platform参数白名单
	PlatformArray = []string{"32", "64"}

	// TypeArray Platform参数白名单
	TypeArray = []string{"agent", "data", "daemon", "agent.exe", "data.zip", "daemon.exe"}

	// ClientHealthTag 主机健康搜索标签
	ClientHealthTag = bson.M{
		"online":       0,
		"offline":      1,
		"can-not-push": 2,
	}

	// FileName2Type 文件名和类型的对应关系
	FileName2Type = bson.M{
		"agent":      "agent",
		"data":       "data",
		"daemon":     "daemon",
		"agent.exe":  "agent",
		"data.zip":   "data",
		"daemon.exe": "daemon",
	}

	// ProjectPath 项目目录
	ProjectPath = ""

	// FilePath 文件路径
	FilePath = ""

	// PageLimit 默认翻页器的Limit
	PageLimit = 10

	// NotDataPrefixLst 无需添加data.前缀的列表
	NotDataPrefixLst = []string{"count"}

	// ValidNoticeQ 有效的告警查询条件，非观察模式
	ValidNoticeQ = bson.M{"status": bson.M{"$lt": 3}}

	// LearnNoticeQ 观察模式下生成的告警的查询条件
	LearnNoticeQ = bson.M{"status": bson.M{"$gt": 2}}

	// Learn2WriteList 观察模式的告警用于添加到白名单的几个类型
	Learn2WriteList = []string{"process", "loginlog", "connection"}

	// StatisticsDBKeys 当搜索条件出现以下key时转换搜索表
	StatisticsDBKeys = []string{"count", "exist_count"}

	// ElasticSearchTypeList 当搜索条件出现以下key时应该调用ElasticSearch接口
	ElasticSearchTypeList = []string{
		"connection", "process", "loginlog", "file", "count",
	}

	// MongoComparisonOperator mongodb比较符对比表
	MongoComparisonOperator = bson.M{
		">":  "$gt",
		"<":  "$lt",
		"=":  "$eq",
		">=": "$gte",
		"<=": "$lte",
	}

	// ConfigTypeMap 根据type判断配置类别
	ConfigTypeMap = map[string][]string{
		"bool": []string{"udp", "lan", "learn", "switch", "onlyhigh", "offlinecheck"},
		"int":  []string{"cycle"},
	}

	// TimeFormat 时间模板
	TimeFormat = "2006-01-02"

	// AnalyzeTypeDict 分析模块提示的关联配置
	AnalyzeTypeDict = []byte(`{
        "ip": "IP地址",
        "type": "数据类型",
        "system": "操作系统名称",
        "data": {
            "name": "数据内容",
            "command": "命令",
            "user": "用户名"
        },
        "crontab": {
            "arg": "启动参数",
            "rule": "",
            "description": "描述"
        },
        "listening": {
            "proto": "类型",
            "address": "监听地址",
            "pid": "监听程序pid"
        },
        "service": {
            "pathname": "启动命令，同command",
            "started": "当前启动状态",
            "startmode": "开机启动模式",
            "startname": "启动用户",
            "caption": "描述"
        },
        "startup": {
            "location": "来源"
        },
        "userlist": {
            "description": "描述 ",
            "status": "状态"
        },
        "file": {
            "path": "文件或者目录路径 file",
            "action": "行为类型",
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
            "pid": "进程pid",
            "ppid": "父进程pid",
            "parentname": "父进程名"
        },
        "connection": {
            "dir": "方向 ",
            "protocol": "类型（TCP、UDP） type",
            "local": "本机进行通讯ip:port",
            "remote": "远程进行通讯的ip:port",
            "pid ": "进程pid"
        }
    }`)

	// StatisticsPipeProjectQ 分析页面Aggregate $project 参数
	StatisticsPipeProjectQ = []byte(`
    {
        "$project": {
            "d" : {"$dateToString" : {"format": "%Y-%m-%d", "date": "$time" }}
        }
    }
    `)

	// StatisticsPipeGroupQ 分析页面Aggregate $group 参数
	StatisticsPipeGroupQ = []byte(`
    {
        "$group": {
            "_id": "$d",
            "count": {
                "$sum": 1
            }
        }
    }
    `)

	// StatisticsTimeProjectQ 分析页面Aggregate $project 参数
	StatisticsTimeProjectQ = []byte(`
    {
        "$project": {
            "d" : {"$dateToString" : {"format": "%m-%d:%H", "date": "$time" }}
        }
    }
    `)

	// LevelString 告警等级和其中文解释的对应关系 {1:"危险", 2:"可疑", 3:"风险"}
	LevelString = []string{"危险", "可疑", "风险"}

	// AuthURILst 需要二次验证的url
	AuthURILst = []string{
		"/file",
		"/config",
		"/tasks",
		"/rules",
	}

	// HTTPURLLst 允许HTTP的url
	HTTPURLLst = []string{
		"/download",
	}

	// DefualtConfig 默认配置
	DefualtConfig = []byte(`
    [
        {
            "type": "client",
            "dic": {
                "cycle": 2,
                "udp": false,
                "lan": false,
                "monitorPath": [
                    "%windows%",
                    "%system32%",
                    "%web%",
                    "/etc/",
                    "/bin/",
                    "/sbin/",
                    "/usr/bin/",
                    "/usr/sbin/"
                ]
            }
        },
        {
            "type": "server",
            "dic": {
                "learn": true,
                "offlinecheck": false,
                "publickey": "",
                "privatekey": "",
                "cert": ""
            }
        },
        {
            "type": "intelligence",
            "dic": {
                "switch": false,
                "ipapi": "http://127.0.0.1/api/?ip={$ip}",
                "fileapi": "http://127.0.0.1/api/?hash={$hash}",
                "regex": "black"
            }
        },
        {
            "type" : "notice",
            "dic" : {
                "switch" : false,
                "onlyhigh" : true,
                "api" : "http://127.0.0.1/test/?text={$info}"
            }
        },
        {
            "type": "whitelist",
            "dic": {
                "file": [],
                "ip": [],
                "process": [],
                "other" : []
            }
        },
        {
            "type": "blacklist",
            "dic": {
                "file": [],
                "ip": [],
                "process": [
                    "mssecsvc\\.exe",
                    "tasksche\\.exe"
                ],
                "other": []
            }
        },
        {
            "type": "filter",
            "dic": {
                "file": ["^c:\\\\windows\\\\temp$","\\.(png|js|css|jpg|gif|wolff|svg)$"],
                "ip": [],
                "process": ["c:\\\\windows\\\\system32\\\\wbem\\\\wmiprvse.exe"]
            }
        },
        {
            "type" : "web",
            "dic" : {
                "tfakey": ""
            }
        }
    ]`)

	// InstallStep 安装步骤次数
	InstallStep = 4

	// SessionGCMaxLifetime session 超时时长
	SessionGCMaxLifetime = int64(60 * 60 * 24 * 15)

	//SecretKeyLst 不必展示的密文
	SecretKeyLst = []string{"cert", "publickey", "privatekey"}

	// key files' name
	PublicKeyName  = "public.pem"
	PrivateKeyName = "private.pem"
	CertKeyName    = "cert.pem"

	// msg list
	AddTaskSucceed = "添加任务成功"
	AddTaskFailure = "添加任务失败"
	EditCfgSucceed = "修改配置成功"
	EditCfgFailure = "修改配置失败"
	Failure        = "您的操作失败，请检查输入是否非法"
	Succeed        = "操作成功"
)
