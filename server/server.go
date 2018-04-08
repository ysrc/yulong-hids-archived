package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"yulong-hids/server/action"
	"yulong-hids/server/models"
	"yulong-hids/server/safecheck"

	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"gopkg.in/Shopify/sarama"
)

const authToken string = "67080fc75bb8ee4a168026e5b21bf6fc"

type Watcher int

const topic = "metrics"

// Kafka 客户端
type Kafka struct {
	consumer sarama.Consumer
}

func newKakfaClient(bs []string) *Kafka {
	config := sarama.NewConfig()
	conn, err := sarama.NewConsumer(bs, config)
	if err != nil {
		log.Fatalf("Connect to kafka error :%s\n", err.Error())
		return nil
	}
	return &Kafka{
		consumer: conn,
	}
}

// GetInfo agent 提交主机信息获取配置信息
func (w *Watcher) GetInfo(ctx context.Context, info *action.ComputerInfo, result *action.ClientConfig) error {
	action.ComputerInfoSave(*info)
	config := action.GetAgentConfig(info.IP)
	log.Println("getconfig:", info.IP)
	*result = config
	return nil
}

// PutInfo 接收处理agent传输的信息
func (k *Kafka) PutInfo() {
	cp, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln("Kafka consume error: ", err.Error())
	}

	// 开始消费数据
	for {
		select {
		case msg := <-cp.Messages():
			if msg != nil {
				var datainfo models.DataInfo
				err := json.Unmarshal(msg.Value, &datainfo)
				if err != nil {
					log.Println("Json unmarshal error: ", err.Error())
					continue
				}
				//保证数据正常
				if len(datainfo.Data) == 0 {
					continue
				}
				log.Println("putinfo:", datainfo.IP, datainfo.Type)
				datainfo.Uptime = time.Now()
				err = action.ResultSave(datainfo)
				if err != nil {
					log.Println(err)
				}
				err = action.ResultStat(datainfo)
				if err != nil {
					log.Println(err)
				}
				safecheck.ScanChan <- datainfo
			}
		}
	}
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

var brokers = flag.String("b", "127.0.0.1:9092", "Kafka brokers")

func main() {
	flag.Parse()
	if *brokers == "" {
		log.Fatal("Kafka Brokers can not be empty")
		return
	}

	if kafkaClient := newKakfaClient(strings.Split(*brokers, ",")); kafkaClient != nil {
		go kafkaClient.PutInfo()
	}

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
