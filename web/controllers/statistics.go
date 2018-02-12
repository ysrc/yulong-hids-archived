package controllers

import (
	"time"
	"yulong-hids/web/models"
	"yulong-hids/web/settings"
	"yulong-hids/web/utils"

	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2/bson"
)

// StatisticsController /statistics
type StatisticsController struct {
	BaseController
}

// Get method
func (c *StatisticsController) Get() {
	var result bson.M
	switchType := c.GetString("type")
	switch switchType {
	case "line":
		// notice number datetime base line chart
		result = getLineData()
	case "pie":
		// type notice pie chart
		result = getPieData()
	case "time":
		// time base
		result = getTimeData()
	case "topmsg":
		// last notice list
		result = getTopMsgData()
	case "total":
		// nav data for each type count
		result = getTotalData()
	}
	c.Data["json"] = result
	c.ServeJSON()
	return
}

func initQueryMap() ([]string, bson.M) {
	orderedkeys := []string{"onlineClient", "allClient", "doneNotice", "allNotice", "ignoreNotice", "doneTask", "allTask", "failTask"}
	querymap := bson.M{
		"onlineClient": bson.M{"health": bson.M{"$in": []int{0, 2}}},
		"allClient":    bson.M{},
		"doneNotice":   bson.M{"status": 1},
		"allNotice":    settings.ValidNoticeQ,
		"ignoreNotice": bson.M{"status": 2},
		"doneTask":     bson.M{"status": "true"},
		"allTask":      bson.M{},
		"failTask":     bson.M{"status": "false"},
	}
	return orderedkeys, querymap
}

func getTotalData() bson.M {
	var result []int
	orderedkeys, querymap := initQueryMap()
	keylist := []string{"Client", "Notice", "Task"}
	for _, name := range orderedkeys {
		var count int
		args := querymap[name]
		key := utils.FindSub(keylist, name)
		switch key {
		case "Client":
			cli := models.NewClient()
			count = cli.Count(args.(bson.M))
		case "Notice":
			cli := models.NewNotice()
			count = cli.Count(args.(bson.M))
		case "Task":
			cli := models.NewTaskResult()
			count = cli.Count(args.(bson.M))
		}
		result = append(result, count)
	}
	hostdata := result[:2]
	alarmdata := result[2:5]
	taskdata := result[5:8]

	return bson.M{
		"hostdata":    hostdata,
		"alarmdata":   alarmdata,
		"taskdata":    taskdata,
		"servicedata": getServiceData(),
		"totaldata":   getInfoCount(),
	}
}

func getServiceData() []float64 {
	serverlist := GetAliveServerList()

	servercount := len(serverlist)

	client := models.NewClient()
	clientcount := client.Count(nil)
	perload, _ := beego.AppConfig.Int("perloadcount")

	var loadPercentage float64
	if servercount > 0 {
		loadPercentage = float64(clientcount) / float64(servercount*perload)
		loadPercentage = loadPercentage * 100
	}

	loadPercentage = utils.Round(loadPercentage, 1)

	return []float64{float64(servercount), loadPercentage}
}

func getInfoCount() []int {
	cli := models.NewInfo()
	infoCount := cli.CountSubList(bson.M{}, "data")

	es := utils.NewSession()
	monitorCount := es.CountAllMonitor()

	return []int{infoCount, int(monitorCount)}
}

func getTopMsgData() bson.M {
	xdata := []string{"标题", "创建时间"}
	var listdata [][]interface{}
	cli := models.NewNotice()
	reslist := cli.GetSortedTop(settings.ValidNoticeQ, 0, 8, "status", "level", "-time")
	for _, res := range reslist {
		timestr := res["time"].(time.Time).Format("2006-01-02 15:04")
		item := []interface{}{res["info"].(string), res["source"].(string), res["level"].(int), timestr}
		listdata = append(listdata, item)
	}
	return newStatisticsJson(xdata, listdata)
}

func getTimeData() bson.M {
	var result [][]int
	cli := models.NewNotice()
	today := utils.TodayRounded()
	var flag = today
	daylist := []int{0, 1}
	for _, day := range daylist {
		match := bson.M{
			"time": bson.M{"$gte": flag.AddDate(0, 0, -1), "$lt": flag},
		}
		match["status"] = settings.ValidNoticeQ["status"]
		flag = flag.AddDate(0, 0, 1)
		reslist := cli.CountPerHour(match)
		for _, res := range reslist {
			pos := res["time"].(time.Time).Hour()
			item := []int{day, pos, res["count"].(int)}
			result = append(result, item)
		}
	}
	xdata := []string{"昨天", "今天"}
	return newStatisticsJson(xdata, result)
}

func getPieData() bson.M {
	var cli = models.NewNotice()
	var listdata []bson.M
	var xdata []string
	var listdataInner []bson.M

	// json["listdata"]
	match := bson.M{"$match": settings.ValidNoticeQ}
	SourceCounts := cli.CountPerByKey(match, "source")
	if len(SourceCounts) > 6 {
		othercount := 0.0
		for _, res := range SourceCounts[:6] {
			listdata = append(listdata, bson.M{"value": res["count"], "name": res["_id"]})
		}
		for _, res := range SourceCounts[6:] {
			othercount += res["count"].(float64)
		}
		listdata = append(listdata, bson.M{"value": othercount, "name": "其他"})
	} else {
		for _, res := range SourceCounts {
			listdata = append(listdata, bson.M{"value": res["count"], "name": res["_id"]})
		}
	}

	// json["xdata"]
	for _, data := range listdata {
		xdata = append(xdata, data["name"].(string))
	}

	// json["listdataInner"]
	LevelCounts := cli.CountPerByKey(match, "level")
	for _, res := range LevelCounts {
		index := res["_id"].(int)
		listdataInner = append(listdataInner, bson.M{"value": res["count"], "name": settings.LevelString[index]})
	}

	result := newStatisticsJson(xdata, listdata)
	result["listdataInner"] = listdataInner
	return result
}

func getLineData() bson.M {
	var result bson.M
	cli := models.NewNotice()
	var listdata []bson.M
	levelstring := settings.LevelString

	// get date
	last7date := utils.Last7DateStr(settings.TimeFormat)
	tomorrow := utils.TodayRounded().AddDate(0, 0, 1)
	sevenDayAgo, _ := time.Parse(settings.TimeFormat, last7date[6])

	// get count per day for each level in last 7 day
	last7date = utils.ReverseStrList(last7date)
	matchDate := bson.M{"$gte": sevenDayAgo, "$lt": tomorrow}
	for level := 0; level < 3; level++ {
		var data [7]float64
		matchQuery := bson.M{"time": matchDate, "level": level, "status": settings.ValidNoticeQ["status"]}
		mgoResult := cli.CountPerDay(matchQuery)
		for index, date := range last7date {
			var count float64
			res := utils.MapSearch(mgoResult, "_id", date)
			if res != nil {
				count = res["count"].(float64)
				data[index] = count
			}
		}
		format := bson.M{
			"name": levelstring[level],
			"type": "line",
			"data": data,
		}
		listdata = append(listdata, format)
	}

	result = newStatisticsJson(last7date, listdata)
	return result
}

func newStatisticsJson(xdata []string, listdata interface{}) bson.M {
	var result bson.M
	result = bson.M{"xdata": xdata, "listdata": listdata}
	return result
}
