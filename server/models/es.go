package models

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/olivere/elastic"
)

var processMapping = `
{
	"properties": {
		"data": {
			"properties": {
				"command": {
					"type": "text",
					"fields": {
						"keyword": {
							"ignore_above": 256,
							"type": "keyword"
						}
					}
				},
				"name": {
					"type": "text",
					"fields": {
						"keyword": {
							"ignore_above": 128,
							"type": "keyword"
						}
					}
				},
				"parentname": {
					"type": "text",
					"fields": {
						"keyword": {
							"ignore_above": 128,
							"type": "keyword"
						}
					}
				},
				"pid": {
					"type": "keyword"
				},
				"ppid": {
					"type": "keyword"
				}
			}
		},
		"ip": {
			"type": "ip"
		},
		"time": {
			"type": "date"
		}
	}
}`
var fileMapping = `
{
	"properties": {
		"data": {
			"properties": {
				"action": {
					"type": "keyword"
				},
				"path": {
					"type": "text",
					"fields": {
						"keyword": {
							"ignore_above": 256,
							"type": "keyword"
						}
					}
				},
				"hash":{
					"type" :"keyword"
				},
				"user":{
					"type" :"string",
					"fields": {
						"keyword": {
							"ignore_above": 40,
							"type": "keyword"
						}
					}
				}
			}
		},
		"ip": {
			"type": "ip"
		},
		"time": {
			"type": "date"
		}
	}
}`

var loginlogMapping = `
{
	"properties": {
		"data": {
			"properties": {
				"username": {
					"type": "text",
					"fields": {
						"keyword": {
							"ignore_above": 40,
							"type": "keyword"
						}
					}
				},
				"hostname": {
					"type": "text",
					"fields": {
						"keyword": {
							"ignore_above": 40,
							"type": "keyword"
						}
					}
				},
				"remote": {
					"type": "text",
					"fields": {
						"keyword": {
							"ignore_above": 25,
							"type": "keyword"
						}
					}
				},
				"status": {
					"type": "keyword"
				}
			}
		},
		"ip": {
			"type": "ip"
		},
		"time": {
			"type": "date"
		}
	}
}`

var connectionMapping = `
{
	"properties": {
		"data": {
			"properties": {
				"dir": {
					"type": "keyword"
				},
				"remote": {
					"type": "text",
					"fields": {
						"keyword": {
							"ignore_above": 25,
							"type": "keyword"
						}
					}
				},
				"local": {
					"type": "text"
				},
				"name":{
					"type":"text",
					"fields": {
						"keyword": {
							"ignore_above": 40,
							"type": "keyword"
						}
					}
				},
				"pid":{
					"type":"keyword"
				},
				"protocol": {
					"type": "keyword"
				}
			}
		},
		"ip": {
			"type": "ip"
		},
		"time": {
			"type": "date"
		}
	}
}`

// ESSave 插入es记录结构
type ESSave struct {
	IP   string            `json:"ip"`
	Data map[string]string `json:"data"`
	Time time.Time         `json:"time"`
}

type esData struct {
	dataType string
	data     ESSave
}

// Client http请求client
var Client *elastic.Client

var esChan chan esData
var nowindicesName string

func init() {
	nowDate := time.Now().Local().Format("2006_01")
	nowindicesName = "monitor" + nowDate
	var err error
	Client, err = elastic.NewClient(elastic.SetURL("http://" + *es))
	if err != nil {
		log.Println("Elastic NewClient error:", err.Error())
		panic(1)
	}
	indexNameList, err := Client.IndexNames()
	if err != nil {
		log.Println("Client IndexNames error:", err.Error())
		return
	}
	if !inArray(indexNameList, nowindicesName, false) {
		newIndex(nowindicesName)
	}
	esChan = make(chan esData, 2048)
}

//InsertThread ES异步写入线程
func InsertThread() {
	var data esData
	p, err := Client.BulkProcessor().
		Name("YulongWorker-1").
		Workers(2).
		BulkActions(100).                // commit if # requests >= 100
		BulkSize(2 << 20).               // commit if size of requests >= 2 MB
		FlushInterval(30 * time.Second). // commit every 30s
		Do(context.Background())
	if err != nil {
		log.Println("start BulkProcessor: ", err)
	}
	for {
		data = <-esChan
		p.Add(elastic.NewBulkIndexRequest().Index(nowindicesName).Type(data.dataType).Doc(data.data))
	}
}

// InsertEs 将数据插入es
func InsertEs(dataType string, data ESSave) {
	esChan <- esData{dataType, data}
}

func esCheckThread() {
	ticker := time.NewTicker(time.Second * 3600)
	for _ = range ticker.C {
		nowDate := time.Now().Local().Format("2006_01")
		nowindicesName = "monitor" + nowDate
		indexNameList, err := Client.IndexNames()
		if err != nil {
			continue
		}
		if inArray(indexNameList, nowindicesName, false) {
			if time.Now().Local().Day() >= 28 {
				nextData := time.Now().Local().AddDate(0, 1, 0).Format("2006_01")
				indicesName := "monitor" + nextData
				if !inArray(indexNameList, indicesName, false) {
					newIndex(indicesName)
				}
			}
		} else {
			newIndex(nowindicesName)
		}
	}
}

func newIndex(name string) {
	log.Println("init indice", name)
	Client.CreateIndex(name).Do(context.Background())
	Client.PutMapping().Index(name).Type("process").BodyString(processMapping).Do(context.Background())
	Client.PutMapping().Index(name).Type("connection").BodyString(connectionMapping).Do(context.Background())
	Client.PutMapping().Index(name).Type("loginlog").BodyString(loginlogMapping).Do(context.Background())
	Client.PutMapping().Index(name).Type("file").BodyString(fileMapping).Do(context.Background())
}

// QueryLogLastTime 查询ip最后一条登录日志的时间
func QueryLogLastTime(ip string) (string, error) {
	termQuery := elastic.NewTermQuery("ip", ip)
	searchResult, err := Client.Search("monitor*").Type("loginlog").Query(termQuery).Sort("time", false).Size(1).Do(context.Background())
	if err != nil {
		return "", err
	}
	if searchResult.Hits.TotalHits != 0 {
		var res map[string]interface{}
		result, err := searchResult.Hits.Hits[0].Source.MarshalJSON()
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(result, &res)
		if err != nil {
			return "", err
		}
		return res["time"].(string), nil
	}
	return "all", nil
}

func inArray(list []string, value string, like bool) bool {
	for _, v := range list {
		if like {
			if strings.Contains(value, v) {
				return true
			}
		} else {
			if value == v {
				return true
			}
		}
	}
	return false
}
