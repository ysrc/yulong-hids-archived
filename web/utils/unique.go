package utils

import (
	"path"
	"path/filepath"
	"regexp"

	"gopkg.in/mgo.v2/bson"

	"github.com/astaxie/beego"
)

// DloadFilePath relative dload path to abs path
func DloadFilePath(propath string) string {
	pathInConfig := beego.AppConfig.String("FilePath")
	if filepath.IsAbs(pathInConfig) {
		return pathInConfig
	}
	return path.Join(propath, pathInConfig)
}

// IsDevMode is debug mode or not
func IsDevMode() bool {
	return beego.AppConfig.String("runmode") == "dev"
}

// AllKeyRegexQuery 为结构体定义所有key添加搜索模糊搜索条件
func AllKeyRegexQuery(filter string, sta interface{}) bson.M {

	var query = make(map[string]interface{})

	if filter != "" {
		allkey := AllStructKey(sta)
		regex := bson.M{"$regex": regexp.QuoteMeta(filter), "$options": "$i"}
		var filterlst []bson.M
		for _, key := range allkey {
			var dict = bson.M{key: regex}
			filterlst = append(filterlst, dict)
		}
		query["$or"] = filterlst
		PPrintBsonM(query)
	} else {
		query = nil
	}

	return query
}
