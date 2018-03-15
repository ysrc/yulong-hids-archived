package main

import (
	"context"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"time"

	"yulong-hids/server/action"
	"yulong-hids/server/models"
	"yulong-hids/server/safecheck"

	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
)

const authToken string = "67080fc75bb8ee4a168026e5b21bf6fc"

type Watcher int

// GetInfo agent 提交主机信息获取配置信息
func (w *Watcher) GetInfo(ctx context.Context, info *action.ComputerInfo, result *action.ClientConfig) error {
	action.ComputerInfoSave(*info)
	config := action.GetAgentConfig(info.IP)
	log.Println("getconfig:", info.IP)
	*result = config
	return nil
}

// PutInfo 接收处理agent传输的信息
func (w *Watcher) PutInfo(ctx context.Context, datainfo *models.DataInfo, result *int) error {
	//保证数据正常
	if len(datainfo.Data) == 0 {
		return nil
	}
	datainfo.Uptime = time.Now()
	log.Println("putinfo:", datainfo.IP, datainfo.Type)
	err := action.ResultSave(*datainfo)
	if err != nil {
		log.Println(err)
	}
	err = action.ResultStat(*datainfo)
	if err != nil {
		log.Println(err)
	}
	safecheck.ScanChan <- *datainfo
	*result = 1
	return nil
}

func auth(ctx context.Context, req *protocol.Message, token string) error {
	if token == authToken {
		return nil
	}
	return errors.New("invalid token")
}

func init() {
	log.Println(models.Config)
	// 从数据库获取证书和RSA私钥
	ioutil.WriteFile("cert.pem", []byte(models.Config.Cert), 0666)
	ioutil.WriteFile("private.pem", []byte(models.Config.Private), 0666)
	// 启动心跳线程
	go models.Heartbeat()
	// 启动推送任务线程
	go action.TaskThread()
	// 启动安全检测线程
	go safecheck.ScanMonitorThread()
	// 启动客户端健康检测线程
	go safecheck.HealthCheckThread()
	// ES异步写入线程
	go models.InsertThread()
}
func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "private.pem")
	if err != nil {
		log.Println("cert error!")
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	s := server.NewServer(server.WithTLSConfig(config))
	s.AuthFunc = auth
	s.RegisterName("Watcher", new(Watcher), "")
	log.Println("RPC Server started")
	err = s.Serve("tcp", ":33433")
	if err != nil {
		log.Println(err.Error())
	}
}
