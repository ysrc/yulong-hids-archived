package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
)

// ElasticSearch struct
type ElasticSearch struct {
	baseurl *url.URL
}

// NewSession init ElasticSearch obj
func NewSession() ElasticSearch {
	var es ElasticSearch
	var err error
	baseurl := beego.AppConfig.String("elastic_search::baseurl")
	es.baseurl, err = url.Parse(baseurl)
	if err != nil {
		beego.Error("Es baseurl", baseurl)
		beego.Error("ElasticSearch base url parse error", err)
	}
	return es
}

// Search search by ElasticSearch web api
func (es ElasticSearch) Search(indexs []string, mode string, query []byte) bson.M {

	var url = *es.baseurl
	dateString := "*"

	for i, index := range indexs {
		if i == 0 {
			index = index + dateString
		}
		url.Path = path.Join(url.Path, index)
	}
	url.Path = path.Join(url.Path, mode)
	req, err := http.NewRequest("POST", url.String(), bytes.NewBuffer(query))
	beego.Debug("Query to ElasticSearch:", string(query))
	if err != nil {
		beego.Error("NewRequest Error", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		beego.Error("http client done requests", err)
	}
	defer resp.Body.Close()

	var res bson.M
	if resp.Status != "200 OK" {
		beego.Error("http requests error : ", url.String(), "'s response status is not 200 is", resp.Status)
		return res
	}

	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &res)

	return res
}

// SearchByJSON call search by bson.M
func (es ElasticSearch) SearchByJSON(indexs []string, query bson.M) bson.M {
	qbyte, _ := json.Marshal(query)
	return es.Search(indexs, "_search", qbyte)
}

// SearchInMonitor request /monitor/_search
func (es ElasticSearch) SearchInMonitor(query []byte) bson.M {
	var res bson.M
	indexs := []string{"monitor"}
	res = es.Search(indexs, "_search", query)
	return res
}

// Last3SecMonitorData request /monitor/_search "gte":"now-ns"
func (es ElasticSearch) LastSecMonitorData(ip string, second int) interface{} {
	var res bson.M
	indexs := []string{"monitor"}
	query := []byte(`{
		"query": {
			"bool": {
				"must": [
					{
						"term": {
							"ip": "` + ip + `"
						}
					},
					{
						"range": {
							"time": {
								"gte":"now-` + strconv.Itoa(second) + `s"
							}
						}
					}
				]
			}
		}
	}`)
	res = es.Search(indexs, "_search", query)
	hits := res["hits"].(map[string]interface{})
	return hits["hits"]
}

// Count count in es search
func (es ElasticSearch) Count(indexs []string, query []byte) float64 {
	res := es.Search(indexs, "_count", query)
	return res["count"].(float64)
}

// CountAllMonitor count all data in index monitor
func (es ElasticSearch) CountAllMonitor() float64 {
	indexs := []string{"monitor"}
	count := es.Count(indexs, []byte(`{}`))
	return count
}
