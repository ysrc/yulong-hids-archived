// +build windows

// Package monitor 监控以下关键行为
// 文件操作、网络连接、命令执行
package monitor

/*
#cgo windows LDFLAGS:-lWS2_32 -liphlpapi
#include <stdio.h>
#include <tchar.h>
#include <windows.h>
#include <Tlhelp32.h>
#include <winsock.h>
#include <iphlpapi.h>


#pragma comment(lib, "ws2_32.lib")
#pragma comment(lib, "iphlpapi.lib")



typedef enum {
    TcpConnectionOffloadStateInHost,
    TcpConnectionOffloadStateOffloading,
    TcpConnectionOffloadStateOffloaded,
    TcpConnectionOffloadStateUploading,
    TcpConnectionOffloadStateMax
} TCP_CONNECTION_OFFLOAD_STATE, *PTCP_CONNECTION_OFFLOAD_STATE;



typedef struct _MIB_TCPROW2 {
    DWORD dwState;
    DWORD dwLocalAddr;
    DWORD dwLocalPort;
    DWORD dwRemoteAddr;
    DWORD dwRemotePort;
    DWORD dwOwningPid;
    TCP_CONNECTION_OFFLOAD_STATE dwOffloadState;
} MIB_TCPROW2, *PMIB_TCPROW2;


typedef struct _MIB_TCPTABLE2 {
    DWORD dwNumEntries;
    MIB_TCPROW2 table[ANY_SIZE];
} MIB_TCPTABLE2, *PMIB_TCPTABLE2;



typedef DWORD(WINAPI *_InternalGetTcpTable2)(
	PMIB_TCPTABLE2 pTcpTable_Vista,
	PULONG SizePointer,
	BOOL Order
	);
static _InternalGetTcpTable2 pGetTcpTable = NULL;


typedef struct tagMIB_TCPEXROW{
	DWORD dwState;              // 连接状态.
	DWORD dwLocalAddr;             // 本地地址.
	DWORD dwLocalPort;           // 本地端口.
	DWORD dwRemoteAddr;            // 远程地址.
	DWORD dwRemotePort;         // 远程端口.
	int dwProcessId;            //进程pid
} MIB_TCPEXROW, *PMIB_TCPEXROW;

typedef struct tagMIB_TCPEXTABLE{
	DWORD dwNumEntries;
	MIB_TCPEXROW table[100];
} MIB_TCPEXTABLE, *PMIB_TCPEXTABLE;


typedef DWORD(WINAPI *PALLOCATE_AND_GET_TCPEXTABLE_FROM_STACK)(
	PMIB_TCPEXTABLE *pTcpTable,
	BOOL bOrder,
	HANDLE heap,
	DWORD zero,
	DWORD flags
	);
static PALLOCATE_AND_GET_TCPEXTABLE_FROM_STACK pAllocateAndGetTcpExTableFromStack = NULL;


ULONG WINAPI GetTcpTable2 (PMIB_TCPTABLE2 TcpTable, PULONG SizePointer, BOOL Order);

ULONG WINAPI GetTcpTable2 (PMIB_TCPTABLE2 TcpTable, PULONG SizePointer, BOOL Order)
{


	return 111;
}



#define MALLOC(x) HeapAlloc(GetProcessHeap(), 0, (x))
#define FREE(x) HeapFree(GetProcessHeap(), 0, (x))
int filter(char *ip, DWORD port)
{
	DWORD i;

	PMIB_TCPTABLE2 TCPTable2ForWin7;
	ULONG ulSize = 0;
	DWORD dwRetVal = 0;
	char szLocalAddr[128];
	char szRemoteAddr[128];
	PMIB_TCPEXTABLE TCPExTable;
	struct in_addr IpAddr;
	HMODULE hIpDLL = LoadLibraryA("iphlpapi.dll");
	if (!hIpDLL)
	{
		printf("LoadLibrary error!\n");
		return 0;
	}
	    int return_code = 0;
		pAllocateAndGetTcpExTableFromStack =
		(PALLOCATE_AND_GET_TCPEXTABLE_FROM_STACK)
		GetProcAddress(hIpDLL, "AllocateAndGetTcpExTableFromStack");

	//vista
	if (pAllocateAndGetTcpExTableFromStack == NULL)
	{


		pGetTcpTable = (_InternalGetTcpTable2)GetProcAddress(hIpDLL, "GetTcpTable2");
		if (pGetTcpTable != NULL)
		{
			TCPTable2ForWin7 = (MIB_TCPTABLE2 *)HeapAlloc(GetProcessHeap(), 0, sizeof (MIB_TCPTABLE2));
			ulSize = sizeof (MIB_TCPTABLE);
			if (NULL == TCPTable2ForWin7)
			{
				printf("allocating memory Error\n");
				FreeLibrary(hIpDLL);
				return 0;
			}
			if ((dwRetVal = pGetTcpTable(TCPTable2ForWin7, &ulSize, TRUE)) ==
				ERROR_INSUFFICIENT_BUFFER) {
				FREE(TCPTable2ForWin7);
				TCPTable2ForWin7 = (MIB_TCPTABLE2 *)MALLOC(ulSize);
				if (TCPTable2ForWin7 == NULL) {
					printf("Error allocating memory\n");
					FreeLibrary(hIpDLL);
					return 0;
				}
			}

			if ((dwRetVal = pGetTcpTable(TCPTable2ForWin7, &ulSize, TRUE)) == NO_ERROR) {
				for (i = 0; i < (int)TCPTable2ForWin7->dwNumEntries; i++) {
					IpAddr.S_un.S_addr = (u_long)TCPTable2ForWin7->table[i].dwRemoteAddr;
					strcpy(szRemoteAddr, inet_ntoa(IpAddr));

					if (strcmp(szRemoteAddr, ip) == 0 && ntohs((u_short)TCPTable2ForWin7->table[i].dwRemotePort) == port)
					{
						return_code = TCPTable2ForWin7->table[i].dwOwningPid;
					}
				}
			}
		}
		FREE(TCPTable2ForWin7);

	}
	else
	{
		//FALSE or TRUE 表明数据是否排序
		if (pAllocateAndGetTcpExTableFromStack(&TCPExTable, FALSE, GetProcessHeap(), 2, AF_INET))
		{
			printf("AllocateAndGetTcpExTableFromStack Error!\n");
			FreeLibrary(hIpDLL);
			return 0;
		}
                int i;
		for (i= 0; i < TCPExTable->dwNumEntries; i++)
		{
			IpAddr.S_un.S_addr = (u_long)TCPExTable->table[i].dwRemoteAddr;
			strcpy(szRemoteAddr, inet_ntoa(IpAddr));
			if (strcmp(szRemoteAddr, ip) == 0 && ntohs((u_short)TCPExTable->table[i].dwRemotePort) == port)
			{
				return_code = TCPExTable->table[i].dwProcessId;
			}
		}
		FREE(TCPExTable);
	}
	FreeLibrary(hIpDLL);
	return return_code;
}
*/
import "C"
import (
	"fmt"
	"strings"
	"yulong-hids/agent/common"

	"github.com/akrennmair/gopcap"
)

// StartNetSniff 开始网络行为监控
func StartNetSniff(resultChan chan map[string]string) {
	var pkt *pcap.Packet
	var resultdata map[string]string
	h, err := getPcapHandle(common.LocalIP)
	if err != nil {
		return
	}
	for {
		pkt = h.Next()
		if pkt == nil {
			continue
		}
		pkt.Decode()
		resultdata = map[string]string{
			"source":   "",
			"dir":      "",
			"protocol": "",
			"remote":   "",
			"local":    "",
			"pid":      "",
			"name":     "",
		}
		var port int
		var ip string
		var localPort int
		//不记录跟安全中心的连接记录
		if pkt.IP != nil && (common.LocalIP == pkt.IP.SrcAddr() || common.LocalIP == pkt.IP.DestAddr()) &&
			!common.InArray(common.ServerIPList, pkt.IP.SrcAddr(), false) &&
			!common.InArray(common.ServerIPList, pkt.IP.DestAddr(), false) {
			resultdata["source"] = "connection"
			if common.LocalIP == pkt.IP.SrcAddr() {
				ip = pkt.IP.DestAddr()
				resultdata["dir"] = "out"
			} else {
				ip = pkt.IP.SrcAddr()
				resultdata["dir"] = "in"
			}
			if common.ServerInfo.Type == "web" && resultdata["dir"] == "in" {
				continue
			}
			//如果内网记录为关闭则进行IP判断
			if !common.Config.LAN && isLan(ip) {
				continue
			}
			//白名单
			if common.InArray(common.Config.Filter.IP, ip, false) {
				continue
			}
			if pkt.IP.Protocol == UDP {
				if common.Config.UDP == false {
					continue
				}
				resultdata["protocol"] = "udp"
				//resultdata["new"] = "1"
				if resultdata["dir"] == "out" {
					port = int(pkt.UDP.DestPort)
					localPort = int(pkt.UDP.SrcPort)
				} else {
					port = int(pkt.UDP.SrcPort)
					localPort = int(pkt.UDP.DestPort)
				}
				if isFilterPort(port) || isFilterPort(localPort) {
					continue
				}
				resultdata["remote"] = fmt.Sprintf("%s:%d", ip, port) //UDP
				resultdata["local"] = fmt.Sprintf("%s:%d", common.LocalIP, localPort)
			} else if pkt.IP.Protocol == TCP && strings.Contains(pkt.String(), "[syn]") {
				resultdata["protocol"] = "tcp"
				if resultdata["dir"] == "out" {
					port = int(pkt.TCP.DestPort)
					localPort = int(pkt.TCP.SrcPort)
				} else {
					port = int(pkt.TCP.SrcPort)
					localPort = int(pkt.TCP.DestPort)
				}
				if isFilterPort(port) || isFilterPort(localPort) {
					continue
				}
				resultdata["remote"] = fmt.Sprintf("%s:%d", ip, port) //TCP
				if pid := C.filter(C.CString(ip), C.DWORD(port)); pid != 0 {
					resultdata["pid"] = fmt.Sprintf("%d", pid)
					if processInfo, ok := getProcessInfo(resultdata["pid"]); ok {
						resultdata["name"] = processInfo.Name
					}
				}
				resultdata["local"] = fmt.Sprintf("%s:%d", common.LocalIP, localPort)
			} else {
				continue
			}
			resultChan <- resultdata
		}
	}
}
