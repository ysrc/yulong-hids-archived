package monitor

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"yulong-hids/agent/common"

	pcap "github.com/akrennmair/gopcap"
	"github.com/go-fsnotify/fsnotify"
)

func getFileMD5(path string) (string, error) {
	fileinfo, err := os.Stat(path)
	// log.Println(fileinfo.Size())
	if fileinfo.Size() >= fileSize {
		return "", errors.New("big file")
	}
	if err != nil {
		return "", err
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	md5Ctx := md5.New()
	if _, err = io.Copy(md5Ctx, file); err != nil {
		return "", err
	}
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr), nil
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

func isLan(ip string) bool {
	ipInt := ipToInt(ip)
	if (ipInt >= 167772160 && ipInt <= 184549375) || (ipInt >= 2886729728 && ipInt <= 2887778303) || (ipInt >= 3232235520 && ipInt <= 3232301055) {
		return true
	}
	return false
}
func isFilterPort(port int) bool {
	for _, v := range filter.Port {
		if v == port {
			return true
		}
	}
	return false
}

func getPcapHandle(ip string) (*pcap.Pcap, error) {
	devs, err := pcap.Findalldevs()
	if err != nil {
		return nil, err
	}
	var device string
	for _, dev := range devs {
		for _, v := range dev.Addresses {
			if v.IP.String() == ip {
				device = dev.Name
				break
			}
		}
	}
	if device == "" {
		return nil, errors.New("find device error")
	}
	h, err := pcap.Openlive(device, 65535, true, 0)
	if err != nil {
		return nil, err
	}
	log.Println("StartConnMonitor")
	err = h.Setfilter("tcp or udp and (not broadcast and not multicast)")
	if err != nil {
		return nil, err
	}
	return h, nil
}

func iterationWatcher(monList []string, watcher *fsnotify.Watcher, pathList []string) {
	for _, p := range monList {
		filepath.Walk(p, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				log.Println(err.Error())
				return err
			}
			if f.IsDir() {
				pathList = append(pathList, path)
				err = watcher.Add(strings.ToLower(path))
				if err != nil {
					log.Println(err)
				}
			}
			return nil
		})
	}
}

// isFileWhite param @resultdata key list: [source, action, path, hash, user]
func isFileWhite(resultdata map[string]string) bool {
	for _, v := range common.Config.Filter.File {
		if ok, _ := regexp.MatchString(`^[0-9a-zA-Z]{32}$`, v); ok {
			if strings.ToLower(v) == strings.ToLower(resultdata["hash"]) {
				return true
			}
		} else {
			if ok, _ := regexp.MatchString(v, strings.ToLower(resultdata["path"])); ok {
				return true
			}
		}
	}
	return false
}
