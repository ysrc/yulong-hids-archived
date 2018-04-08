// +build linux
package main

/*
#include <sys/socket.h>
#include <linux/netlink.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>

#define NETLINK_USER 31

#define MAX_PAYLOAD 2048
struct sockaddr_nl src_addr, dest_addr;
struct nlmsghdr *nlh = NULL;
struct iovec iov;
struct msghdr msg;

#define PORT  65530

int CapturePrecess()
{
	//udp sock
	int sock;
	int payload_max_len = 0;

	payload_max_len = pathconf("/", _PC_PATH_MAX);
	if(payload_max_len < 0) {
		return -1;
	}

	payload_max_len += MAX_PAYLOAD;

	sock = socket(AF_INET, SOCK_DGRAM, 0);
	if(sock < 0) {
		return -1;
	}
	struct sockaddr_in sockaddrin;
	memset(&sockaddrin, 0, sizeof(sockaddrin));
	sockaddrin.sin_family = AF_INET;
	sockaddrin.sin_port = htons(PORT);
	sockaddrin.sin_addr.s_addr = inet_addr("127.0.0.1");

	//netlink sock
	int sock_fd;
	sock_fd = socket(PF_NETLINK, SOCK_RAW, NETLINK_USER);
	if (sock_fd < 0) {
		return -1;
	}

	memset(&src_addr, 0, sizeof(src_addr));
    memset(&msg, 0, sizeof(msg));

    src_addr.nl_family = AF_NETLINK;
    src_addr.nl_pid = getpid();
    src_addr.nl_groups = 1;
    bind(sock_fd, (struct sockaddr*)&src_addr, sizeof(src_addr));
    memset(&dest_addr, 0, sizeof(dest_addr));
    nlh = (struct nlmsghdr *)malloc(NLMSG_SPACE(payload_max_len));
    memset(nlh, 0, NLMSG_SPACE(payload_max_len));

    iov.iov_base = (void *)nlh;
    iov.iov_len = NLMSG_SPACE(payload_max_len);
    msg.msg_name = (void *)&dest_addr;
    msg.msg_namelen = sizeof(dest_addr);
    msg.msg_iov = &iov;
    msg.msg_iovlen = 1;

	while (1)
	{
		recvmsg(sock_fd, &msg, 0);
		sendto(sock, (char *)NLMSG_DATA(nlh), strlen((char *)NLMSG_DATA(nlh)), 0, (struct sockaddr *)&sockaddrin, sizeof(sockaddrin));
		memset((char *)NLMSG_DATA(nlh), 0, strlen((char *)NLMSG_DATA(nlh)));
	}
	close(sock_fd);
	close(sock);
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
)

func main() {
	log.Println("StartProcessMonitor")
	var buf [255]byte
	//开启进程监控提取线程
	go func() {
		ok := C.CapturePrecess()
		if ok < 0 {
			log.Println("connect syshook netlink error")
		}
	}()
	localaddress, _ := net.ResolveUDPAddr("udp", "127.0.0.1:65530")
	udplistener, err := net.ListenUDP("udp", localaddress)
	if err != nil {
		log.Print(err.Error())
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
		//log.Println(string(buf[0 : n-1]))
		// name|command|pid|pname|ppid|info
		//进程名|参数|进程PID|父进程|父进程PID
		log.Println(string(buf[0:n]))
		proList := strings.Split(string(buf[0:n-1]), string(0x01))
		if len(proList) < 5 {
			log.Println(string(buf[0:n]))
			continue
		}
		//不记录agent执行的命令  || s == os.Getppid()
		if s, _ := strconv.Atoi(proList[4]); s == os.Getpid() {
			continue
		}
		resultdata = make(map[string]string)
		resultdata["source"] = "process"
		//resultdata["type"] = proList[0]
		resultdata["name"] = proList[0]
		resultdata["command"] = proList[1]
		resultdata["pid"] = proList[2]
		resultdata["parentname"] = proList[3]
		resultdata["ppid"] = proList[4]
		resultdata["info"] = proList[5]
		fmt.Println(resultdata)
		//fmt.Print(string(buf[0:n]))
	}
}
