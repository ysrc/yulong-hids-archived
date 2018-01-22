# 规则编写
----------
## 规则介绍

可通过此功能自定义警报规则，格式如下：

```
{
  "and": true, // rules的逻辑，true为全部符合，false则只需要符合一个
  "enabled": true, // 规则启用开关
  "meta": {
    "author": "wolf",
    "description": "Guest用户正常情况为禁止状态",
    "level": 0, // 警报等级，0-2，分别为危险、可疑、提示
    "name": "Guest用户异常"
  },
  "rules": {
    "name": { 
      "data": "Guest", // 判断值
      "type": "string" // 判断方式（string、regex、non-regex、count）
    }, // 用户名为Guest，key为字段
    "status": {
      "data": "OK",
      "type": "string"
    } // 状态为启用
  },
  "source": "userlist", // 数据来源类型
  "system": "windows" // 操作系统类型（windows,linux）
}
```
> 正则表达式(regex,non-regex)的相关字符串匹配需使用小写字母，字符串(string)则不区分大小写。  
> 部分内置规则在不同的环境下可能会存在误报和无效，需根据自身环境和业务特点进行改动。（例如`可疑动态脚本写入`规则，如果你的web服务是以管理员权限运行或者与代码发布所有者权限一致的话将无法发挥作用）

这里引用[职业欠钱](https://xianzhi.aliyun.com/forum/topic/1626/)关于入侵检测基本原则的描述，在定义规则的时候可以思考一下。

1. 不能把每一条告警都彻底跟进的模型，等同于无效模型 ——有入侵了再说之前有告警，只是太多了没跟过来/没查彻底，这是马后炮，等同于不具备发现能力；
2. 我们必须屏蔽一些重复发生的相似的误报告警，以集中精力对每一个告警都闭环掉 —— 这会产生白名单，也就是漏报，因此单个模型的漏报是不可避免的；
3. 由于任何单模型都会存在漏报，所以我们必须在多个纬度上做多个模型，形成纵深 —— 假设WebShell静态文本分析被黑客变形绕过了，在RASP（运行时环境）的恶意调用还可以监控到，这样可以选择接受单个模型的漏报，但在整体上仍然不漏；
4. 任何模型都有误报漏报，我们做什么，不做什么，需要考虑的是“性价比” —— 比如某些变形的WebShell可以写成跟业务代码非常相似，人的肉眼几乎无法识别，再追求一定要在文本分析上进行对抗，就是性价比很差的决策，通过RASP的检测方案，其性价比更高一些；
5. 我们不可能知道黑客所有的攻击手法，也不可能针对每一种手法都建设策略（不具备性价比），但是，针对重点业务，我们可以通过加固的方式，让黑客能攻击的路径极度收敛，仅在关键环节进行对抗（包括加固的有效性检测）可能会让100%的目标变得现实

> 基于上述几个原则，我们可以知道一个事实，或许，我们永远不可能在单点上做到100分，但是，我们可以通过一些组合方式，让攻击者很难绕过所有的点。

## 数据结构

**来源-字段**
- **crontab** // 计划任务
  - name // 计划任务名
  - command // 要执行的程序或命令以及参数
  - arg // 启动参数
  - user // 启动用户
  - rule
  - description // 描述 
- **listening** // 监听TCP端口
  - proto // 类型
  - address // 监听地址
  - name // 监听程序名
  - pid // 监听程序pid
- **service** // 服务
  - name // 服务名
  - pathname // 启动命令，同command
  - started // 当前启动状态
  - startmode // 开机启动模式
  - startname // 启动用户
  - caption // 描述
- **startup** // 开机启动项
  - name // 名称
  - command // 启动程序或命令
  - location // 来源
  - user // 启动用户
- **userlist** // 用户列表
  - name // 用户名
  - description // 描述 
  - status // 状态
- **file** // 文件操作行为
  - path // 文件或者目录路径
  - action 行为类型
  - user // 操作用户
  - hash // 文件md5 hash
- **loginlog** // 系统登录日志
  - username // 用户名
  - hostname // 远程主机名
  - remote // 远程IP
  - status // 认证结果
  - time // 时间
- **process** // 进程创建事件
  - name // 进程名
  - command // 程序或命令以及参数
  - pid // 进程pid
  - ppid // 父进程pid
  - parentname // 父进程名
  - info // 进程其他相关信息
- **connection** // 网络连接事件
  - dir // 方向
  - protocol // 类型（TCP、UDP）
  - local // 本机进行通讯ip:port
  - remote // 远程进行通讯的ip:port
  - name // 进程名
  - pid  // 进程pid

> 此结构windows、linux通用，但可能有一些细微区别，具体数据内容可在web控制台的数据分析功能查看



