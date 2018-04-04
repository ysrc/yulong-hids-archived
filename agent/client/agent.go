package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
	"yulong-hids/agent/collect"
	"yulong-hids/agent/common"
	"yulong-hids/agent/monitor"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"gopkg.in/Shopify/sarama"
)

var err error

type dataInfo struct {
	IP     string              `json:"ip"`     // 客户端的IP地址
	Type   string              `json:"type"`   // 传输的数据类型
	System string              `json:"system"` // 操作系统
	Data   []map[string]string `json:"data"`   // 数据内容
}

const sendTopic = "metrics"

// KafkaClient 客户端
type KafkaClient struct {
	// 生产者
	producer sarama.SyncProducer
}

// Agent agent客户端结构
type Agent struct {
	ServerNetLoc  string         // 服务端地址 IP:PORT
	Client        client.XClient // RPC 客户端
	ServerList    []string       // 存活服务端集群列表
	PutData       dataInfo       // 要传输的数据
	Reply         int            // RPC Server 响应结果
	Mutex         *sync.Mutex    // 安全操作锁
	IsDebug       bool           // 是否开启debug模式，debug模式打印传输内容和报错信息
	KafkaProducer KafkaClient    // kafka生产者
	ctx           context.Context
}

var httpClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
}

func (a *Agent) init() {
	a.ServerList, err = a.getServerList()
	if err != nil {
		a.log("GetServerList error:", err)
		panic(1)
	}
	a.ctx = context.WithValue(context.Background(), share.ReqMetaDataKey, make(map[string]string))
	a.log("Available server node:", a.ServerList)
	if len(a.ServerList) == 0 {
		time.Sleep(time.Second * 30)
		a.log("No server node available")
		panic(1)
	}
	a.newClient()
	if common.LocalIP == "" {
		a.log("Can not get local address")
		panic(1)
	}
	a.Mutex = new(sync.Mutex)
	err := a.Client.Call(a.ctx, "GetInfo", &common.ServerInfo, &common.Config)
	if err != nil {
		a.log("RPC Client Call Error:", err.Error())
		panic(1)
	}

	conf := readConfig()
	if conf != nil {
		if brokers := strings.Split(conf.KafkaBroker, ","); brokers != nil {
			config := sarama.NewConfig()
			config.Producer.Return.Successes = true
			a.KafkaProducer, err = sarama.NewSyncProducer(brokers)
			if err != nil {
				a.log(err.Error())
			}
		}
	}
	a.log("Common Client Config:", common.Config)
}

// Run 启动agent
func (a *Agent) Run() {

	// agent 初始化
	// 请求Web API，获取Server地址，初始化RPC客户端，获取客户端IP等
	a.init()

	// 每隔一段时间更新初始化配置
	a.configRefresh()

	// 开启各个监控流程 文件监控，网络监控，进程监控
	a.monitor()

	// 每隔一段时间获取系统信息
	// 监听端口，服务信息，用户信息，开机启动项，计划任务，登录信息，进程列表等
	a.getInfo()
}

func (a *Agent) newClient() {
	var servers []*client.KVPair
	for _, server := range a.ServerList {
		common.ServerIPList = append(common.ServerIPList, strings.Split(server, ":")[0])
		s := client.KVPair{Key: server}
		servers = append(servers, &s)
		if common.LocalIP == "" {
			a.setLocalIP(server)
			common.ServerInfo = collect.GetComInfo()
			a.log("Host Information:", common.ServerInfo)
		}
	}
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	option := client.DefaultOption
	option.TLSConfig = conf
	serverd := client.NewMultipleServersDiscovery(servers)
	a.Client = client.NewXClient("Watcher", FAILMODE, client.RandomSelect, serverd, option)
	a.Client.Auth(AUTH_TOKEN)
}

func (a Agent) getServerList() ([]string, error) {
	var serlist []string
	var url string
	if TESTMODE {
		url = "http://" + a.ServerNetLoc + SERVER_API
	} else {
		url = "https://" + a.ServerNetLoc + SERVER_API
	}
	a.log("Web API:", url)
	request, _ := http.NewRequest("GET", url, nil)
	request.Close = true
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(result), &serlist)
	if err != nil {
		return nil, err
	}
	return serlist, nil
}

func (a Agent) setLocalIP(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		a.log("Net.Dial:", ip)
		a.log("Error:", err)
		panic(1)
	}
	defer conn.Close()
	common.LocalIP = strings.Split(conn.LocalAddr().String(), ":")[0]
}

func (a *Agent) configRefresh() {
	ticker := time.NewTicker(time.Second * time.Duration(CONFIGR_REF_INTERVAL))
	go func() {
		for _ = range ticker.C {
			ch := make(chan struct{})
			go func() {
				err = a.Client.Call(a.ctx, "GetInfo", &common.ServerInfo, &common.Config)
				if err != nil {
					a.log("RPC Client Call:", err.Error())
					return
				}
				close(ch)
			}()
			// Server集群列表获取
			select {
			case <-ch:
				serverList, err := a.getServerList()
				if err != nil {
					a.log("RPC Client Call:", err.Error())
					break
				}
				if len(serverList) == 0 {
					a.log("No server node available")
					break
				}
				if len(serverList) == len(a.ServerList) {
					for i, server := range serverList {
						// TODO 可能会产生问题
						if server != a.ServerList[i] {
							a.ServerList = serverList
							// 防止正在传输重置client导致数据丢失
							a.Mutex.Lock()
							a.Client.Close()
							a.newClient()
							a.Mutex.Unlock()
							break
						}
					}
				} else {
					a.log("Server nodes from old to new:", a.ServerList, "->", serverList)
					a.ServerList = serverList
					a.Mutex.Lock()
					a.Client.Close()
					a.newClient()
					a.Mutex.Unlock()
				}
			case <-time.NewTicker(time.Second * 3).C:
				break
			}
		}
	}()
}

func (a *Agent) monitor() {
	resultChan := make(chan map[string]string, 16)
	go monitor.StartNetSniff(resultChan)
	go monitor.StartProcessMonitor(resultChan)
	go monitor.StartFileMonitor(resultChan)
	go func(result chan map[string]string) {
		var resultdata []map[string]string
		var data map[string]string
		for {
			data = <-result
			data["time"] = fmt.Sprintf("%d", time.Now().Unix())
			a.log("Monitor data: ", data)
			source := data["source"]
			delete(data, "source")
			a.Mutex.Lock()
			a.PutData = dataInfo{common.LocalIP, source, runtime.GOOS, append(resultdata, data)}
			a.put()
			a.Mutex.Unlock()
		}
	}(resultChan)
}

func (a *Agent) getInfo() {
	historyCache := make(map[string][]map[string]string)
	for {
		if len(common.Config.MonitorPath) == 0 {
			time.Sleep(time.Second)
			a.log("Failed to get the configuration information")
			continue
		}
		allData := collect.GetAllInfo()
		for k, v := range allData {
			if len(v) == 0 || a.mapComparison(v, historyCache[k]) {
				a.log("GetInfo Data:", k, "No change")
				continue
			} else {
				a.Mutex.Lock()
				a.PutData = dataInfo{common.LocalIP, k, runtime.GOOS, v}
				a.put()
				a.Mutex.Unlock()
				if k != "service" {
					a.log("Data details:", k, a.PutData)
				}
				historyCache[k] = v
			}
		}
		if common.Config.Cycle == 0 {
			common.Config.Cycle = 1
		}
		time.Sleep(time.Second * time.Duration(common.Config.Cycle) * 60)
	}
}

func (a Agent) put() {
	_, err := a.Client.Go(a.ctx, "PutInfo", &a.PutData, &a.Reply, nil)
	if err != nil {
		a.log("PutInfo error:", err.Error())
	}

	// 发送payload至消息队列
	go func() {
		if a.KafkaProducer != nil {
			payload, err := json.Marshal(&a.PutData)
			if err != nil {
				a.log(err.Error())
				return
			}
			msg := &sarama.ProducerMessage{Topic: sendTopic, Value: sarama.ByteEncoder(payload)}
			_, _, err = a.KafkaProducer.SendMessage(msg)
			if err != nil {
				a.log(err.Error())
			}
		}
	}()
}

func (a Agent) mapComparison(new []map[string]string, old []map[string]string) bool {
	if len(new) == len(old) {
		for i, v := range new {
			for k, value := range v {
				if value != old[i][k] {
					return false
				}
			}
		}
		return true
	}
	return false
}

func (a Agent) log(info ...interface{}) {
	if a.IsDebug {
		log.Println(info...)
	}
}
