// +build linux
package monitor

/*
#include <arpa/inet.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <fcntl.h>
#include <pwd.h>
#include <errno.h>
#include <dirent.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/stat.h>

typedef union iaddr iaddr;

union iaddr {
    unsigned u;
    unsigned char b[4];
};

#define PRG_LOCAL_ADDRESS "local_address"
#define PRG_INODE	 "inode"
#define PRG_SOCKET_PFX    "socket:["
#define PRG_SOCKET_PFXl (strlen(PRG_SOCKET_PFX))
#define PRG_SOCKET_PFX2   "[0000]:"
#define PRG_SOCKET_PFX2l  (strlen(PRG_SOCKET_PFX2))

#ifndef LINE_MAX
#define LINE_MAX 4096
#endif

#define PATH_PROC	   "/proc"
#define PATH_FD_SUFF	"fd"
#define PATH_FD_SUFFl       strlen(PATH_FD_SUFF)
#define PATH_PROC_X_FD      PATH_PROC "/%s/" PATH_FD_SUFF
#define PATH_CMDLINE	"cmdline"
#define PATH_CMDLINEl       strlen(PATH_CMDLINE)

#undef  DIRENT_HAVE_D_TYPE_WORKS


#define ADDR_LEN INET6_ADDRSTRLEN + 1 + 5 + 1

static void addr2str(int af, const void *addr, unsigned port, char *buf)
{
    if (inet_ntop(af, addr, buf, ADDR_LEN) == NULL) {
        *buf = '\0';
        return;
    }
}
#define PROGNAME_WIDTH 20

#define PRG_HASH_SIZE 211
#define PRG_HASHIT(x) ((x) % PRG_HASH_SIZE)

static char finbuf[PROGNAME_WIDTH];

static void extract_type_1_socket_inode(const char lname[], long * inode_p) {

    if (lname[strlen(lname) - 1] != ']') *inode_p = -1;
	else {
		char inode_str[strlen(lname + 1)];
		const int inode_str_len = strlen(lname) - PRG_SOCKET_PFXl - 1;
		char *serr;

		strncpy(inode_str, lname + PRG_SOCKET_PFXl, inode_str_len);
		inode_str[inode_str_len] = '\0';
		*inode_p = strtol(inode_str, &serr, 0);
		if (!serr || *serr || *inode_p < 0 || *inode_p >= 2147483647)
			*inode_p = -1;
	}
}


static char * pget(unsigned uid,long inode) {
    DIR *d = opendir("/proc");
    if(NULL == d) return "";
	long inode_p;
    char statline[1024];
    char cmdline[1024];
    struct dirent *de;
    struct stat stats;
    struct passwd *pw;
	int fd, i = 0; int procfdlen, fd_last, cmdllen, lnamelen;
	char line[LINE_MAX], eacces = 0;
	char lname[30], cmdlbuf[512];
	DIR *dirproc = NULL, *dirfd = NULL;
	struct dirent *direproc, *direfd;
	cmdlbuf[sizeof(cmdlbuf)-1] = '\0';
	const char *cs, *cmdlp;
	if (!(dirproc = opendir(PATH_PROC)))
	{
		printf("error");
		return "";
	}
	while (direproc = readdir(dirproc)) {
		for (cs = direproc->d_name; *cs; cs++)
		if (!isdigit(*cs))
			break;
		if (*cs)
			continue;

		procfdlen = snprintf(line, sizeof(line), PATH_PROC_X_FD, direproc->d_name);
		if (procfdlen <= 0 || procfdlen >= sizeof(line)-5)
			continue;
		errno = 0;
		dirfd = opendir(line);
		if (!dirfd) {
			if (errno ==1 )
				eacces = 1;
			continue;
		}
		line[procfdlen] = '/';
		cmdlp = NULL;

		while ((direfd = readdir(dirfd))) {
			if (procfdlen + 1 + strlen(direfd->d_name) + 1>sizeof(line))
				continue;
			memcpy(line + procfdlen - PATH_FD_SUFFl, PATH_FD_SUFF "/",
				PATH_FD_SUFFl + 1);
			strcpy(line + procfdlen + 1, direfd->d_name);
			lnamelen = readlink(line, lname, sizeof(lname)-1);
			lname[lnamelen] = '\0';
			extract_type_1_socket_inode(lname, &inode_p);
			if (inode_p == inode)
			{
				if (!cmdlp) {
					if (procfdlen - PATH_FD_SUFFl + PATH_CMDLINEl >=
						sizeof(line)-5)
						continue;
					strcpy(line + procfdlen - PATH_FD_SUFFl, PATH_CMDLINE);
					fd = open(line, O_RDONLY);
					if (fd < 0)
						continue;
					cmdllen = read(fd, cmdlbuf, sizeof(cmdlbuf)-1);
					if (close(fd))
						continue;
					if (cmdllen == -1)
						continue;
					if (cmdllen < sizeof(cmdlbuf)-1)
						cmdlbuf[cmdllen] = '\0';
					if ((cmdlp = strrchr(cmdlbuf, '/')))
						cmdlp++;
					else
						cmdlp = cmdlbuf;
				}

				snprintf(finbuf, sizeof(finbuf), "%s/%s", direproc->d_name, cmdlp);
				return finbuf;
			}
		}
	}
    return "";
}
char *s;
static char * filter(char *host, int port) {
	char *filename = "/proc/net/tcp";
	char *label = "tcp";
    memset(finbuf,0,PROGNAME_WIDTH);
    FILE *fp = fopen(filename, "r");
    if (fp == NULL) return;
	long inode;
    char buf[BUFSIZ];
    fgets(buf, BUFSIZ, fp);
    while (fgets(buf, BUFSIZ, fp)){
        char lip[ADDR_LEN];
        char rip[ADDR_LEN];
		char more[512];
        iaddr laddr, raddr;
        unsigned lport, rport, state, txq, rxq, num, tr, tm_when, retrnsmt, uid;
		int timeout;
        int n = sscanf(buf, " %d: %x:%x %x:%x %x %x:%x %x:%x %x %d %d %ld %512s",
                       &num, &laddr.u, &lport, &raddr.u, &rport,
					   &state, &txq, &rxq, &tr, &tm_when, &retrnsmt, &uid, &timeout, &inode, more);
        if (n == 15) {
            addr2str(AF_INET, &laddr, lport, lip);
            addr2str(AF_INET, &raddr, rport, rip);
			if ( (int)rport == port)
			{
				pget(uid, inode);
				fclose(fp);
				return finbuf;
			}
        }
    }
    fclose(fp);
	return "";
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
		var port int
		var ip string
		var localPort int
		resultdata = map[string]string{
			"source":   "",
			"dir":      "",
			"protocol": "",
			"remote":   "",
			"local":    "",
			"pid":      "",
			"name":     "",
		}
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
				conName := C.GoString(C.filter(C.CString(ip), C.int(port)))
				if conName != "" {
					resultdata["pid"] = strings.SplitN(conName, "/", 2)[0]
					resultdata["name"] = strings.SplitN(conName, "/", 2)[1]
				}
				resultdata["local"] = fmt.Sprintf("%s:%d", common.LocalIP, localPort)
			} else {
				continue
			}
			resultChan <- resultdata
		}
	}
}
