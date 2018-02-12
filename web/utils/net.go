package utils

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

// TCPAlive check tcp alive
func TCPAlive(server string) bool {
	timeOut := time.Duration(3) * time.Second
	conn, err := net.DialTimeout("tcp", server, timeOut)
	if err != nil {
		beego.Debug("tcp connect server error", server, err)
		return false
	}
	conn.Close()
	return true
}

// Int2IP Convert uint to net.IP
func Int2IP(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// IP2Int Convert net.IP to int64
func IP2Int(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

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

// String2NetIP Convert string to net.IP
func String2NetIP(ipstr string) net.IP {
	return net.ParseIP(ipstr)
}

// NetIP2String Convert net.IP to string
func NetIP2String(ip net.IP) string {
	return ip.String()
}

// BetweenIP generate ip string from one to another
func BetweenIP(ipsrt1 string, ipstr2 string) []string {
	ip1 := String2NetIP(ipsrt1)
	ip2 := String2NetIP(ipstr2)
	ipint1 := IP2Int(ip1)
	ipint2 := IP2Int(ip2)
	res := []string{}

	if ipint2 > ipint1 {
		for ipint1 <= ipint2 {
			res = append(res, Int2IP(ipint1).String())
			ipint1 = ipint1 + 1
		}
	}

	return res
}
