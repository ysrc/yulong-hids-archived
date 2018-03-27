package controllers

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"yulong-hids/web/models"
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2/bson"
)

var (
	// TypeConf contain setting.AnalyzeTypeDict
	TypeConf map[string]interface{}
)

// AnalyzeController /analyze
type AnalyzeController struct {
	BaseController
}

// Post method POST
func (c *AnalyzeController) Post() {

	var result interface{}
	keyword := c.GetString("keyword")
	result = kwParser(keyword)
	c.Data["json"] = result
	c.ServeJSON()
	return

}

// Get method GET
func (c *AnalyzeController) Get() {

	query := c.GetString("q")
	isEsType := false

	p := c.InitPaginator()
	start, limit := p.ToParameter()

	queryobj, err := queryString2Map(query)
	if err != nil {
		c.Data["json"] = models.NewErrorInfo("type-not-in-query")
		c.ServeJSON()
		return
	}

	if isCountQuery(queryobj) {
		isEsType = false
	} else {
		isEsType = utils.StringInSlice(queryobj["type"], settings.ElasticSearchTypeList)
	}

	if isEsType {
		timequery := c.GetString("tq")
		c.Data["json"] = searchInElasticSearch(queryobj, timequery, start, limit)
	} else {
		c.Data["json"] = searchInMongo(queryobj, start, limit)
	}

	c.ServeJSON()
	return

}

// isCountQuery check is StatisticsDBKeys in query or not
func isCountQuery(query map[string]string) bool {
	for _, ktype := range settings.StatisticsDBKeys {
		_, exist := query[ktype]
		if exist {
			return true
		}
	}
	return false
}

// searchInElasticSearch as name
func searchInElasticSearch(queryobj map[string]string, timequery string, start int, limit int) interface{} {

	var indexs = []string{"monitor"}
	var mustlst []bson.M
	var keyNotPrefixData = []string{"ip", "time"}

	for key, value := range queryobj {
		if key == "type" {
			indexs = append(indexs, value)
		} else if key == "all" {
			// {"match": {"_all": ""}}
			allquery := bson.M{"match": bson.M{"_all": value}}
			mustlst = append(mustlst, allquery)
		} else {
			var eskey string
			if utils.StringInSlice(key, keyNotPrefixData) {
				eskey = key
			} else {
				eskey = "data." + key
			}
			// {"match": {"": ""}}
			termquery := bson.M{"match": bson.M{eskey: value}}
			mustlst = append(mustlst, termquery)
		}
	}

	// time query string to bson.M
	re := regexp.MustCompile("^[\\d]+\\-[\\dnow]+$")
	if re.MatchString(timequery) {
		splitlst := strings.Split(timequery, "-")
		timeMatch := bson.M{
			"range": bson.M{
				"time": bson.M{
					"gte": splitlst[0],
					"lte": splitlst[1],
				},
			},
		}
		mustlst = append(mustlst, timeMatch)
	}

	// {"query":{"bool":{"must":[]}}}
	var query = bson.M{"query": bson.M{"bool": bson.M{"must": mustlst}}}
	query["from"] = start
	query["size"] = limit
	query["sort"] = []bson.M{bson.M{"time": bson.M{"order": "desc"}}}

	// request es web api
	var result interface{}
	es := utils.NewSession()
	result = es.SearchByJSON(indexs, query)

	result = result.(bson.M)["hits"]
	result = result.(map[string]interface{})["hits"]

	itemlist := result.([]interface{})
	var datalist []interface{}
	for _, item := range itemlist {
		var data interface{}
		data = item.(map[string]interface{})["_source"]
		datalist = append(datalist, data)
	}

	return datalist
}

// searchInMongo as name
func searchInMongo(queryobj map[string]string, start int, limit int) []bson.M {
	var result []bson.M
	matchQuery, projectQuery, dbstring := searchParser(queryobj)

	// two collection for search : info and statistics
	if dbstring == "Info" {
		cli := models.NewInfo()
		result = cli.Aggregate(
			bson.M{"$match": matchQuery},
			bson.M{"$skip": start},
			bson.M{"$limit": limit},
			bson.M{"$project": projectQuery},
			bson.M{"$sort": bson.M{"uptime": -1}},
		)
	} else {
		cli := models.NewStatistics()
		result = cli.Query(matchQuery, start, limit)
	}

	return result
}

// queryString2Map return dict type from parse query
func queryString2Map(query string) (map[string]string, error) {
	queryObj := utils.SplitStrToMap(query, "|", ":")

	if _, flg := queryObj["type"]; !flg {
		return queryObj, errors.New("'type:' not in query string")
	}

	return queryObj, nil
}

// searchParser : parse the query string to aggregate query list [match, project]
// more doc string in doc/dev.md#searchParser and analyze search design and doc
func searchParser(queryObj map[string]string) (bson.M, bson.M, string) {
	var dbstring string
	infocli := models.NewInfo()
	matchQuery := bson.M{}
	projectQuery := bson.M{"_id": 1, "ip": 1, "type ": 1, "system": 1, "uptime": 1}
	condAndQuery := []bson.M{}

	for skey, svalue := range queryObj {
		// Judge the collections to use Statistics or Info
		if utils.StringInSlice(skey, settings.StatisticsDBKeys) {
			dbstring = "Statistics"
		} else {
			dbstring = "Info"
		}

		notDataPrefixLst := append(utils.AllStructKey(infocli), settings.NotDataPrefixLst...)

		if skey == "exist_count" {
			// parse to {"server_list": {"$exists":true}, "$where":"this.tag.length>3"}
			if match, _ := regexp.MatchString("^[\\>\\<\\=]{1,2}\\d+$", svalue); match {
				skey = "server_list"
				matchQuery[skey] = bson.M{"$exists": true}
				matchQuery["$where"] = "this.server_list.length" + svalue
			}

		} else if utils.StringInSlice(skey, notDataPrefixLst) {
			// parse to {"skey":"value"}
			var value interface{}
			if skey == "count" {
				kvalue, operator := parseIntQuery(svalue)
				value = bson.M{operator: kvalue}
			} else {
				value = svalue
			}
			matchQuery[skey] = value

		} else {
			// parse to {"data.skey":"value"}
			matchQuery["data."+skey] = svalue
			condAndQuery = append(condAndQuery, bson.M{"$eq": []string{"$$d." + skey, svalue}})
		}
	}

	dataFilter := bson.M{"$filter": bson.M{"input": "$data", "as": "d", "cond": bson.M{"$and": condAndQuery}}}
	projectQuery["data"] = dataFilter
	beego.Debug("Show matchQuery, projectQuery, dbstring: ", matchQuery, projectQuery, dbstring)
	return matchQuery, projectQuery, dbstring
}

// kwParser return the search item input tips
func kwParser(keyword string) interface{} {

	infocli := models.NewInfo()
	result := map[string]interface{}{
		"msg":        "",
		"samplelist": make([]string, 0),
	}
	json.Unmarshal(settings.AnalyzeTypeDict, &TypeConf)

	if keyword == "" {
		keylist := utils.AllKey(TypeConf)
		result["samplelist"] = append(
			result["samplelist"].([]string), keylist...,
		)
	} else if strings.HasSuffix(keyword, ":") {
		key := strings.TrimRight(keyword, ":")
		value := utils.GetValue(TypeConf, key)

		if !utils.StringInSlice(key, utils.AllStructKey(infocli)) {
			key = "data." + key
		}
		if key == "type" {
			statistic := models.NewStatistics()
			result["samplelist"] = statistic.AllValue(key, 20)
		} else {
			result["samplelist"] = infocli.AllValue(key, 20)
		}

		if reflect.ValueOf(value).Kind() == reflect.String {
			result["msg"] = value
		}
	}

	return result
}

// parseIntQuery parse query string to int and comparison operators
func parseIntQuery(qvalue string) (int64, string) {

	var operator = "$eq"
	var value int64
	var err error

	value, err = strconv.ParseInt(qvalue, 0, 64)
	if err != nil {
		for common, oper := range settings.MongoComparisonOperator {
			afterTrim := strings.TrimLeft(qvalue, common)
			if value, err = strconv.ParseInt(afterTrim, 0, 64); err == nil {
				operator = oper.(string)
				break
			}
		}
	}

	return value, operator
}

func toTimeMatch(timequery string) bson.M {
	return nil
}
