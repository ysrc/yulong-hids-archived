// +build windows

package monitor

/*
#cgo windows LDFLAGS:  -lWS2_32
#include <windows.h>
#include <stdio.h>
#include "define.h"



#define host "127.0.0.1"
#define port 65530
#pragma comment(lib, "ws2_32.lib")

int CapturePrecess()
{
    HANDLE        hDevice;
    int        status;
    HANDLE        m_hCommEvent;
    ULONG        dwReturn;
    char        outbuf[255];
    CHECKLIST    CheckList;

     SOCKET sock;
   WSADATA wsaData;
    struct sockaddr_in saddr;
    if (WSAStartup(MAKEWORD(2, 2), &wsaData) != 0)
    {
        printf("error");
    }
    saddr.sin_family = AF_INET;
    saddr.sin_port = htons(port);
    saddr.sin_addr.s_addr = inet_addr(host);
    sock = socket(AF_INET, SOCK_DGRAM, 0);
    hDevice = NULL;
    m_hCommEvent = NULL;
    hDevice = CreateFile("\\\\.\\MonitorProcess",
        GENERIC_READ | GENERIC_WRITE,
        FILE_SHARE_READ | FILE_SHARE_WRITE,
        NULL,
        OPEN_EXISTING,
        FILE_ATTRIBUTE_NORMAL,
        NULL);
    if (hDevice == INVALID_HANDLE_VALUE)
    {
        printf("createfile wrong\n");
        getchar();
        return 0;
    }

    m_hCommEvent = CreateEvent(NULL,
        0,
        0,
        NULL);
    printf("hEvent:%08x\n", m_hCommEvent);

    status = DeviceIoControl(hDevice,
        IOCTL_PASSEVENT,
        &m_hCommEvent,
        sizeof(m_hCommEvent),
        NULL,
        0,
        &dwReturn,
        NULL);
    if (!status)
    {
        printf("IO wrong+%d\n", GetLastError());
        getchar();
        return 0;
    }

    CheckList.ONLYSHOWREMOTETHREAD = TRUE;
    CheckList.SHOWTHREAD = TRUE;
    CheckList.SHOWTERMINATETHREAD = FALSE;
    CheckList.SHOWTERMINATEPROCESS = FALSE;
    status = DeviceIoControl(hDevice,
        IOCTL_PASSEVSTRUCT,
        &CheckList,
        sizeof(CheckList),
        NULL,
        0,
        &dwReturn,
        NULL);
    if (!status)
    {
        printf("IO wrong+%d\n", GetLastError());
        getchar();
        return 0;
    }

    while (1)
    {
        ResetEvent(m_hCommEvent);
        WaitForSingleObject(m_hCommEvent, INFINITE);
        status = DeviceIoControl(hDevice,
            IOCTL_PASSBUF,
            NULL,
            0,
            &outbuf,
            sizeof(outbuf),
            &dwReturn,
            NULL);
        if (!status)
        {
            printf("IO wrong+%d\n", GetLastError());
            getchar();
            return 0;
        }
        sendto(sock,outbuf,strlen(outbuf),0,(struct sockaddr *)&saddr,sizeof(saddr));


    }

    status = DeviceIoControl(hDevice,
        IOCTL_UNPASSEVENT,
        NULL,
        0,
        NULL,
        0,
        &dwReturn,
        NULL);
    if (!status)
    {
        printf("UNPASSEVENT wrong+%d\n", GetLastError());
        getchar();
        return 0;
    }

    status = CloseHandle(hDevice);
    status = CloseHandle(m_hCommEvent);
    getchar();
    return 0;
}
*/
import "C"
import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"yulong-hids/agent/common"

	"github.com/StackExchange/wmi"
)

type process struct {
	Name        string
	CommandLine *string
}

func getProcessInfo(pid string) (process, bool) {
	var dst []process
	err := wmi.Query(fmt.Sprintf("SELECT * FROM Win32_Process where ProcessID = %s", pid), &dst)
	if err == nil && len(dst) != 0 {
		return dst[0], true
	}
	return process{}, false
}

// StartProcessMonitor 开始进程监控
func StartProcessMonitor(resultChan chan map[string]string) {
	log.Println("StartProcessMonitor")
	var buf [255]byte
	// 开启进程监控提取线程
	go C.CapturePrecess()
	localaddress, _ := net.ResolveUDPAddr("udp", "127.0.0.1:65530")
	udplistener, err := net.ListenUDP("udp", localaddress)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	defer udplistener.Close()
	var resultdata map[string]string
	for {
		n, _, err := udplistener.ReadFromUDP(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}
		// TYPE|进程名|进程PID|父进程|父进程PID
		proList := strings.Split(string(buf[0:n-1]), "|")
		if s, _ := strconv.Atoi(proList[4]); s == os.Getpid() {
			continue
		}
		resultdata = make(map[string]string)
		resultdata["source"] = "process"
		resultdata["name"] = proList[1]
		resultdata["pid"] = proList[2]
		resultdata["parentname"] = proList[3]
		resultdata["ppid"] = proList[4]
		resultdata["command"] = ""
		resultdata["info"] = ""
		if processInfo, ok := getProcessInfo(proList[2]); ok {
			resultdata["command"] = *processInfo.CommandLine
			// 驱动获取的进程名长度最高只有14
			// if len(resultdata["name"]) == 14 {
			resultdata["name"] = processInfo.Name
			// }
		}
		if common.InArray(common.Config.Filter.Process, strings.ToLower(resultdata["name"]), true) ||
			common.InArray(common.Config.Filter.Process, strings.ToLower(resultdata["command"]), true) {
			continue
		}
		resultChan <- resultdata
	}
}
