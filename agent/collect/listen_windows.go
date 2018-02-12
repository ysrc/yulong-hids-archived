// +build windows

package collect

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/StackExchange/wmi"
)

var libiphlpapiDll = syscall.MustLoadDLL("iphlpapi.dll")

type MIB_TCPROW struct {
	DwState      uint32
	DwLocalAddr  uint32
	DwLocalPort  uint32
	DwRemoteAddr uint32
	DwRemotePort uint32
}

type MIB_TCPROW2 struct {
	DwState        uint32
	DwLocalAddr    uint32
	DwLocalPort    uint32
	DwRemoteAddr   uint32
	DwRemotePort   uint32
	DwOwningPid    uint32
	DwOffloadState uint32
}
type MIB_TCPTABLE2 struct {
	DwNumEntries uint32
	Table        [65535]MIB_TCPROW2
}
type MIB_TCPTABLE struct {
	DwNumEntries uint32
	Table        [65535]MIB_TCPROW
}
type listening struct {
	LocalAddr string // 监听地址
	LocalPort int    // 监听端口
	PID       int    // 监听进程的PID
}

func getProcessName(pid string) (string, bool) {
	var dst []process
	err := wmi.Query(fmt.Sprintf("SELECT * FROM Win32_Process where ProcessID = %s", pid), &dst)
	if err == nil && len(dst) != 0 {
		return dst[0].Name, true
	}
	return "", false
}

// getTCPTable 获取tcp表
func getTCPTable(ptb *MIB_TCPTABLE, dwSize *uint32, border uint32) uint32 {
	iphelpapiGetTCPTable := libiphlpapiDll.MustFindProc("GetTcpTable")
	ret, _, _ := iphelpapiGetTCPTable.Call(
		uintptr(unsafe.Pointer(ptb)),
		uintptr(unsafe.Pointer(dwSize)),
		uintptr(border),
	)
	return uint32(ret)
}

func getTCPTable2(ptb *MIB_TCPTABLE2, dwSize *uint32, border uint32) uint32 {
	iphelpapiGetTCPTable := libiphlpapiDll.MustFindProc("GetTcpTable2")
	ret, _, _ := iphelpapiGetTCPTable.Call(
		uintptr(unsafe.Pointer(ptb)),
		uintptr(unsafe.Pointer(dwSize)),
		uintptr(border),
	)
	return uint32(ret)
}

func getListening2003() (resultData []map[string]string) {
	var size uint32
	ptb := &MIB_TCPTABLE{}
	if getTCPTable(ptb, &size, 1) != 122 {
		return resultData
	}
	if getTCPTable(ptb, &size, 1) != 0 {
		return resultData
	}
	var tr listening
	var m map[string]string
	for _, tcp := range ptb.Table[0:int(ptb.DwNumEntries)] {
		if int(tcp.DwState) == 2 {
			m = make(map[string]string)
			if tr.LocalAddr == "127.0.0.1" || strings.HasPrefix(tr.LocalAddr, "169.254.") {
				continue
			}
			m["address"] = fmt.Sprintf("%s:%d", addrDecode(tcp.DwLocalAddr), portDecode(tcp.DwLocalPort))
			m["proto"] = "TCP"
			resultData = append(resultData, m)
		}
	}
	return resultData
}
func getListeningOther() (resultData []map[string]string) {
	var size uint32
	ptb := &MIB_TCPTABLE2{}
	if getTCPTable2(ptb, &size, 1) != 122 {
		return resultData
	}
	if getTCPTable2(ptb, &size, 1) != 0 {
		return resultData
	}
	var tr listening
	var m map[string]string
	for _, tcp := range ptb.Table[0:int(ptb.DwNumEntries)] {
		if int(tcp.DwState) == 2 {
			m = make(map[string]string)
			if tr.LocalAddr == "127.0.0.1" || strings.HasPrefix(tr.LocalAddr, "169.254.") {
				continue
			}
			m["address"] = fmt.Sprintf("%s:%d", addrDecode(tcp.DwLocalAddr), portDecode(tcp.DwLocalPort))
			m["pid"] = fmt.Sprintf("%d", int(tcp.DwOwningPid))
			m["name"] = ""
			if name, ok := getProcessName(m["pid"]); ok {
				m["name"] = name
			}
			m["proto"] = "TCP"
			resultData = append(resultData, m)
		}
	}
	return resultData
}

// GetListening 获取tcp端口监听端口
func GetListening() (resultData []map[string]string) {
	if _, err := os.Stat(os.Getenv("SystemDrive") + `/Users`); err != nil {
		fmt.Println("2003")
		resultData = getListening2003()
	} else {
		resultData = getListeningOther()
	}
	return resultData
}

func addrDecode(addr uint32) string {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, addr)
	return net.IP(buf.Bytes()).String()
}

func portDecode(port uint32) int {
	return int(binary.BigEndian.Uint16([]byte{byte(port), byte(port >> 8)}))
}
