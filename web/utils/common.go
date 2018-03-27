package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"yulong-hids/web/settings"

	"github.com/astaxie/beego"

	"gopkg.in/mgo.v2/bson"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// Md5String return MD5 hex string from origin string
func Md5String(str string) string {
	md5obj := md5.New()
	md5obj.Write([]byte(str))
	hex := fmt.Sprintf("%x", md5obj.Sum(nil))
	return hex
}

//GetFileMD5Hash return MD5 hex string from file
func GetFileMD5Hash(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		beego.Error("Open file error:", err)
		return ""
	}
	defer file.Close()
	md5h := md5.New()
	io.Copy(md5h, file)
	return fmt.Sprintf("%x", md5h.Sum([]byte("")))
}

// KeyEncode base64, why it need?
func KeyEncode(key string) string {
	if key == "" {
		return ""
	}
	bs := []byte(key)
	encodeStr := base64.StdEncoding.EncodeToString(bs)
	return fmt.Sprintf("%s", encodeStr)
}

// KeyDecode why?
func KeyDecode(encodeStr string) string {
	if encodeStr == "" {
		return ""
	}
	return encodeStr
}

//RStrip emmm..
func RStrip(s string, suffixlist []string) string {
	for _, element := range suffixlist {
		s = strings.TrimRight(s, element)
	}
	return s
}

// StringInSlice check string is a element of slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//GetCwd return current path
func GetCwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//AnyHasSuffix change HasSuffix function in any element
func AnyHasSuffix(s string, suffixlist []string) bool {
	for _, element := range suffixlist {
		if strings.HasSuffix(s, element) {
			return true
		}
	}
	return false
}

// RandStringBytesMaskImprSrc return random string
func RandStringBytesMaskImprSrc(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// SplitStrToMap Split string to map[string]string
func SplitStrToMap(str string, sep1 string, sep2 string) map[string]string {
	sublist := strings.Split(str, sep1)
	var resultmap = map[string]string{}
	for _, pair := range sublist {
		z := strings.SplitN(pair, sep2, 2)
		if len(z) == 2 {
			key := strings.TrimSpace(z[0])
			value := strings.TrimSpace(z[1])
			key = strings.ToLower(key)
			resultmap[key] = value
		}
	}
	return resultmap
}

func ValueInListMap(str string, listmap map[string][]string) string {
	key := ""
	for key, typelist := range listmap {
		typeArray := typelist
		if StringInSlice(str, typeArray) {
			return key
		}
	}
	return key
}

// KeyType change string to bool or int
func KeyType(key string, value string) interface{} {
	var vresult interface{}
	typemap := settings.ConfigTypeMap

	typestr := ValueInListMap(key, typemap)
	if typestr == "int" {
		vresult, _ = strconv.Atoi(value)
	} else if typestr == "bool" {
		vresult, _ = strconv.ParseBool(value)
	} else {
		vresult = value
	}

	return vresult
}

// FindSub which one is my sun?
// input (["a", "b"],"ac")
// return "a"
func FindSub(strlist []string, str string) string {
	for _, key := range strlist {
		if strings.Contains(str, key) {
			return key
		}
	}
	return ""
}

// AllKey return all key in map and sub map
func AllKey(skmap map[string]interface{}) []string {
	var result []string
	for key := range skmap {
		if reflect.ValueOf(skmap[key]).Kind() == reflect.String {
			result = append(result, key)
		} else {
			res := AllKey(skmap[key].(map[string]interface{}))
			result = append(result, res...)
		}
	}
	return result
}

// GetValue get value in map and sub map by key
func GetValue(skmap map[string]interface{}, key string) interface{} {
	result := false
	for k := range skmap {
		if reflect.ValueOf(skmap[k]).Kind() == reflect.Map {
			res := GetValue(skmap[k].(map[string]interface{}), key)
			if reflect.ValueOf(res).Kind() == reflect.String {
				return res
			}
		} else if k == key {
			return skmap[k]
		}
	}
	return result
}

// AllStructKey get all key in struct
func AllStructKey(sta interface{}) []string {

	var keylist []string
	re := regexp.MustCompile(",[\\w,]+$")

	val := reflect.ValueOf(sta)
	for i := 0; i < val.Type().NumField(); i++ {
		fieldnames := val.Type().Field(i).Tag.Get("json")
		if startswith := strings.HasPrefix(fieldnames, "_"); !startswith && fieldnames != "" {
			fieldnames = re.ReplaceAllString(fieldnames, "")
			keylist = append(keylist, fieldnames)
		}
	}

	return keylist
}

// Last7DateStr get formated date string list in last seven day
func Last7DateStr(timeformat string) []string {
	if timeformat == "" {
		timeformat = "01-02"
	}
	var last7DateList []string
	for day := 0; day < 7; day++ {
		date := TodayRounded().AddDate(0, 0, -day)
		datestr := date.Format(timeformat)
		last7DateList = append(last7DateList, datestr)
	}
	return last7DateList
}

// TodayRounded Get last midnight time
func TodayRounded() time.Time {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return today
}

// PPrintMap Petty Print Golang Map Struct
func PPrintMap(mlist ...map[string]interface{}) {
	for _, m := range mlist {
		b, _ := json.Marshal(m)
		fmt.Println(string(b))
	}
}

//ReverseStrList Reverse a strings slice
func ReverseStrList(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// MapSearch Return a map by filter in a map list
func MapSearch(ml []bson.M, key string, value interface{}) map[string]interface{} {
	for _, submap := range ml {
		if submap[key] == value {
			return submap
		}
	}
	return nil
}

// Round float64
func Round(f float64, n int) float64 {
	truncarg := math.Pow10(n)
	return math.Trunc((f+0.5/truncarg)*truncarg) / truncarg
}

// MapUpdate append a map to another
func MapUpdate(ori map[string]interface{}, sub map[string]interface{}) map[string]interface{} {
	for k, v := range sub {
		ori[k] = v
	}
	return ori
}

// ParseBsonM 把bson.M对象转化为人类可读的字符串
func ParseBsonM(obj interface{}) string {
	byteobj, _ := json.Marshal(obj)
	return string(byteobj)
}

// PPrintBsonM 把bson.M对象转化为人类可读的字符串
func PPrintBsonM(obj interface{}) {
	beego.Debug(ParseBsonM(obj))
}

// InterfaceSlice2BsonM 对interface{}的切片进行类型转换
func InterfaceSlice2BsonM(slice []interface{}) []bson.M {
	var reslice []bson.M
	for _, itf := range slice {
		reslice = append(reslice, itf.(bson.M))
	}
	return reslice
}

// ToBsonMSlice interface{} to []bson.M
func ToBsonMSlice(itf interface{}) []bson.M {
	return InterfaceSlice2BsonM(itf.([]interface{}))
}

//DeleteElementInSlient as name
func DeleteElementInSlient(sl []string, s string) []string {
	for i, v := range sl {
		if v == s {
			sl = append(sl[:i], sl[i+1:]...)
			break
		}
	}
	return sl
}
