package safecheck

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"yulong-hids/server/models"
)

func isLan(ip string) bool {
	if strings.Count(ip, ".") != 3 {
		return true
	}
	ipInt := ipToInt(ip)
	if (ipInt >= 167772160 && ipInt <= 184549375) || (ipInt >= 2886729728 && ipInt <= 2887778303) || (ipInt >= 3232235520 && ipInt <= 3232301055) {
		return true
	}
	return false
}
func ipToInt(ip string) int64 {
	bits := strings.Split(ip, ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])
	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}
func inArray(list []string, value string, regex bool) bool {
	for _, v := range list {
		if regex {
			if v == "" {
				continue
			}
			if ok, err := regexp.MatchString(v, value); ok {
				return true
			} else if err != nil {
				log.Println(err.Error())
			}
		} else {
			if value == v {
				return true
			}
		}
	}
	return false
}
func sendNotice(level int, info string) error {
	log.Println(info)
	if models.Config.Notice.Switch {
		if models.Config.Notice.OnlyHigh {
			if level == 0 {
				_, err := http.Get(strings.Replace(models.Config.Notice.API, "{$info}", url.QueryEscape(info), 1))
				if err != nil {
					return err
				}
			}
			return nil
		}
		_, err := http.Get(strings.Replace(models.Config.Notice.API, "{$info}", url.QueryEscape(info), 1))
		if err != nil {
			return err
		}
	}
	return nil
}
