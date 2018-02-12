#include <sys/stat.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <string.h>
#include <asm/types.h>
#include <linux/netlink.h>
#include <linux/socket.h>

#define NETLINK_TEST 31
#define MAX_PAYLOAD  (2048) /* maximum payload size*/
struct sockaddr_nl src_addr, dest_addr;
struct nlmsghdr *nlh = NULL;
struct iovec iov;
int sock_fd;
struct msghdr msg;

int main(int argc, char* argv[]) 
{
    
    int max_payload_len = 0;
    
    if((max_payload_len = pathconf("/", _PC_PATH_MAX)) == -1) {
        return -1;
    }
    
    printf("%d\n", max_payload_len);
    sock_fd=socket(PF_NETLINK, SOCK_RAW, NETLINK_TEST);
    if(sock_fd < 0) {
        printf("create nl failed.\n");
        return -1;
    }
    memset(&src_addr, 0, sizeof(src_addr));
    memset(&msg, 0, sizeof(msg));

    src_addr.nl_family = AF_NETLINK;
    src_addr.nl_pid = getpid();  /* self pid */
    /* interested in group 1<<0 */
    src_addr.nl_groups = 1;
    bind(sock_fd, (struct sockaddr*)&src_addr, sizeof(src_addr));
    memset(&dest_addr, 0, sizeof(dest_addr));
    nlh = (struct nlmsghdr *)malloc(NLMSG_SPACE(MAX_PAYLOAD));
    memset(nlh, 0, NLMSG_SPACE(MAX_PAYLOAD));

    iov.iov_base = (void *)nlh;
    iov.iov_len = NLMSG_SPACE(MAX_PAYLOAD);
    msg.msg_name = (void *)&dest_addr;
    msg.msg_namelen = sizeof(dest_addr);
    msg.msg_iov = &iov;
    msg.msg_iovlen = 1;


    while (1) {
    /* Read message from kernel */
    recvmsg(sock_fd, &msg, 0);

    printf("----%s\n", NLMSG_DATA(nlh));
    }
    close(sock_fd);         

    return 0;
}
