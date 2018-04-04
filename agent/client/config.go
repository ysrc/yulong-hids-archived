package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/smallnest/rpcx/client"
)

const (
	AUTH_TOKEN           string          = "67080fc75bb8ee4a168026e5b21bf6fc"
	CONFIGR_REF_INTERVAL int             = 60
	FAILMODE             client.FailMode = client.Failtry
	SERVER_API           string          = "/json/serverlist"
	TESTMODE             bool            = false
	CONF_FILE            string          = "./broker.conf"
)

// Config 配置文件
type Config struct {
	KafkaBroker string `json:"kafka"`
}

// 读取配置文件
func readConfig() *Config {
	data, err := ioutil.ReadFile(CONF_FILE)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return &c
}
