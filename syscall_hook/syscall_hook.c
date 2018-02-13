#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/syscalls.h>
#include <linux/delay.h>
#include <linux/file.h>
#include <asm/paravirt.h>
#include <asm/syscall.h>
#include <linux/sys.h>
#include <linux/slab.h>
#include <linux/kallsyms.h>
#include <linux/binfmts.h>
#include <linux/version.h>
#include <net/sock.h>
#include <net/netlink.h>

unsigned long **sys_call_table_ptr;
unsigned long original_cr0;
void *orig_sys_call_table [NR_syscalls];

struct sock *syshook_nl_sk = NULL;
#define SYSHOOK_NL_NUM  31

#if LINUX_VERSION_CODE >= KERNEL_VERSION(3, 10, 0)
    struct user_arg_ptr {
    #ifdef CONFIG_COMPAT
        bool is_compat;
    #endif
        union {
            const char __user *const __user *native;
    #ifdef CONFIG_COMPAT
            const compat_uptr_t __user *compat;
    #endif
        } ptr;
    };
    struct filename *(*tmp_getname)(const char __user * filename);
    void (*tmp_putname)(struct filename *name);
    typedef asmlinkage long (*func_execve)(const char __user *,
                                           const char __user * const __user *,
                                           const char __user *const  __user *);
    extern asmlinkage long monitor_stub_execve_hook (const char __user *,
                                                     const char __user *const __user *,
                                                     const char __user *const __user *);
#elif LINUX_VERSION_CODE == KERNEL_VERSION(2, 6, 32)
    typedef asmlinkage long (*func_execve)(const char __user *,
                                           const char __user * const __user *,
                                           const char __user *const  __user *,
                                           struct pt_regs *);
    extern asmlinkage long monitor_stub_execve_hook(const char __user *,
                                                    const char __user * const __user *,
                                                    const char __user *const  __user *,
                                                    struct pt_regs *);
#endif

func_execve orig_stub_execve;

unsigned long **find_sys_call_table(void) {
    unsigned long ptr;
    unsigned long *p;

    pr_err("Start found sys_call_table.\n");
    
    for (ptr = (unsigned long)sys_close;
         ptr < (unsigned long)&loops_per_jiffy;
         ptr += sizeof(void *)) {

        p = (unsigned long *)ptr;

        if (p[__NR_close] == (unsigned long)sys_close) {
            pr_err("Found the sys_call_table!!! __NR_close[%d] sys_close[%lx]\n"
                    " __NR_execve[%d] sct[__NR_execve][0x%lx]\n",
                    __NR_close,
                    (unsigned long)sys_close,
                    __NR_execve,
                    p[__NR_execve]);
            return (unsigned long **)p;
        }
    }
    
    return NULL;
}



#if LINUX_VERSION_CODE == KERNEL_VERSION(2, 6, 32)
static int tmp_count(char __user * __user * argv, int max)
{
    int i = 0;

    if (argv != NULL) {
        for (;;) {
            char __user * p;

            if (get_user(p, argv))
                return -EFAULT;
            if (!p)
                break;
            argv++;
            if (i++ >= max)
                return -E2BIG;

            if (fatal_signal_pending(current))
                return -ERESTARTNOHAND;
            cond_resched();
        }
    }
    return i;
}

asmlinkage long monitor_execve_hook(char __user *name,
                                   char __user * __user *argv,
                                   char __user * __user *envp, 
                                   struct pt_regs *regs)
{
    long error = 0;
    struct filename *path = NULL;
    char __user * native = NULL;
    int tmp_argc = 0, tmp_envpc =  0;
    int i = 0, len = 0, offset = 0, max_len = 0;
    int total_argc_len = 0, total_envpc_len = 0;
    char *total_argc_ptr = NULL, *total_envpc_ptr = NULL;
    char *per_envp = NULL;
    int nl_send_len = 0;
    struct sk_buff *skb = NULL;
    struct nlmsghdr *nlh = NULL;
    struct file *file = NULL;
    char *tmp = kmalloc(PATH_MAX, GFP_KERNEL);
    char *path1 = NULL;


    path = getname(name);
    error = PTR_ERR(path);
    if (IS_ERR(path)) {
        pr_err("get path failed.\n");
        goto err;
    }
    
    file = open_exec(path->name);
    if (!IS_ERR(file) && tmp) {
        memset(tmp, 0, PATH_MAX);
        path1 = d_path(&file->f_path, tmp, PATH_MAX);
        if (IS_ERR(path1)) {
            path1 = NULL;
        }

        fput(file);
    }
    
    error = 0;
    tmp_argc = tmp_count(argv, MAX_ARG_STRINGS);
    if(tmp_argc < 0) {
        error = tmp_argc;
        goto err;
    }
   
    for(i = 0; i < tmp_argc; i ++) {
        if(get_user(native, argv + i)) {
            error = -EFAULT;
            goto err;
        }

        len = strnlen_user(native, MAX_ARG_STRLEN);
        if(!len) {
            error = -EFAULT;
            goto err;
        }

        total_argc_len += len;
    }

    total_argc_ptr = kmalloc(total_argc_len + 16 * tmp_argc, GFP_ATOMIC);
    if(!total_argc_ptr) {
        error = -ENOMEM;
        goto err;
    }
    memset(total_argc_ptr, 0, total_argc_len + 16 * tmp_argc);
    
    for(i = 0; i < tmp_argc; i ++) {
        if(i == 0) {
            continue;
        }
        if(get_user(native, argv + i)) {
            error = -EFAULT;
            goto err;
        }

        len = strnlen_user(native, MAX_ARG_STRLEN);
        if(!len) {
            error = -EFAULT;
            goto err;
        }

        if(offset + len > total_argc_len + 16 * tmp_argc) {
            break;
        }

        if (copy_from_user(total_argc_ptr + offset, native, len)) {
            error = -EFAULT;
            goto err;
        }

        offset += len - 1;
        *(total_argc_ptr + offset) = ' ';
        offset += 1;
    }
    
    /*--------envp--------------*/
    len = 0;
    offset = 0;
    tmp_envpc = tmp_count(envp, MAX_ARG_STRINGS);
    if(tmp_envpc < 0) {
        error = tmp_envpc;
        goto err;
    }

    for(i = 0; i < tmp_envpc; i ++) {
        if(get_user(native, envp + i)) {
            error = -EFAULT;
            goto err;
        }

        len = strnlen_user(native, MAX_ARG_STRLEN);
        if(!len) {
            error = -EFAULT;
            goto err;
        }

        if(len > max_len) {
            max_len = len;
        }
        total_envpc_len += len;
    }

    per_envp = kmalloc(max_len + 16, GFP_KERNEL);
    if(!per_envp) {
        error = -ENOMEM;
        goto err;
    }
    
    total_envpc_ptr = kmalloc(total_envpc_len + 16 * tmp_envpc, GFP_ATOMIC);
    if(!total_envpc_ptr) {
        error = -ENOMEM;
        goto err;
    }
    memset(total_envpc_ptr, 0, total_envpc_len + 16 * tmp_envpc);
    
    for(i = 0; i < tmp_envpc; i ++) {
        if(get_user(native, envp + i)) {
            error = -EFAULT;
            goto err;
        }
        
        len = strnlen_user(native, MAX_ARG_STRLEN);
        if(!len) {
            error = -EFAULT;
            goto err;
        }

        if(offset + len > total_envpc_len + 16 * tmp_envpc) {
            break;
        }

        memset(per_envp, 0, max_len);
        if(copy_from_user(per_envp, native, len)) {
            error = -EFAULT;
            goto err;
        }

        if(!strstr(per_envp, "PWD") && !strstr(per_envp, "LOGNAME") && !strstr(per_envp, "USER")) {
            continue;
        }

        if (copy_from_user(total_envpc_ptr + offset, native, len)) {
            error = -EFAULT;
            goto err;
        }

        offset += len - 1;
        *(total_envpc_ptr + offset) = ' ';
        offset += 1;
    }
    
    nl_send_len = (path1 != NULL ? strlen(path1) : strlen(path->name)) + strlen(current->parent->comm) + 128;
    if(!ZERO_OR_NULL_PTR(total_envpc_ptr)) {
        nl_send_len += strlen(total_envpc_ptr);
    }
    if(!ZERO_OR_NULL_PTR(total_argc_ptr)) {
        nl_send_len += strlen(total_argc_ptr);
    }
    
    nl_send_len = nl_send_len < PATH_MAX + 2048 ? nl_send_len : PATH_MAX + 2048;
    skb = alloc_skb(NLMSG_SPACE(nl_send_len), GFP_ATOMIC);
    if(!skb) {
        error = -ENOMEM;
        goto err;
    }
    
    nlh = (struct nlmsghdr *)skb->data;
    nlh->nlmsg_len = NLMSG_SPACE(nl_send_len);
    nlh->nlmsg_pid = 0; 
    nlh->nlmsg_flags = 0;
    nlh = nlmsg_put(skb, 0, 0, 0, NLMSG_SPACE(nl_send_len) - sizeof (struct nlmsghdr), 0);
    if(!nlh) {
        kfree_skb(skb);
        pr_err("nlh get failed.\n");
        goto err;
    }
    
    snprintf(NLMSG_DATA(nlh), nl_send_len, "%s%c%s%c%u%c%s%c%d%c%s", path1 != NULL ? path1 : path->name, 0x1, (uint64_t)total_argc_ptr == 0x10 ? "N/A" : total_argc_ptr, 0x1, current->tgid, 0x1, current->parent->comm, 0x1, current->parent->tgid, 0x1, (uint64_t)total_envpc_ptr == 0x10 ? "N/A" : total_envpc_ptr);
    NETLINK_CB(skb).pid = 0;
    NETLINK_CB(skb).dst_group = 1;
    error = netlink_broadcast(syshook_nl_sk, skb, 0, 1, GFP_KERNEL);
    if(error != 0 && error != -3) {
        pr_err("send nl broadcast failed.\n");
        goto err;
    }
    
    //pr_err("%s|%s|%u|%s|%d|%s\n", path1 != NULL ? path1 : path->name, total_argc_ptr, current->tgid, current->parent->comm, current->parent->tgid, total_envpc_ptr);

err:
    if(tmp) {
        kfree(tmp);
        tmp = NULL;
    }
    if(total_envpc_ptr) {
        kfree(total_envpc_ptr);
        total_envpc_ptr = NULL;
    }
    if(per_envp) {
        kfree(per_envp);
        per_envp = NULL; 
    }
    if(total_argc_ptr) {
        kfree(total_argc_ptr);
        total_argc_ptr = NULL;
    }
    putname(path);
    return 0;
}
#elif LINUX_VERSION_CODE >= KERNEL_VERSION(3, 10, 0)
static const char __user *get_user_arg_ptr(struct user_arg_ptr argv, int nr) 
{
    const char __user *native;

#ifdef CONFIG_COMPAT
    if (unlikely(argv.is_compat)) {
        compat_uptr_t compat;

        if (get_user(compat, argv.ptr.compat + nr))
            return ERR_PTR(-EFAULT);

        return compat_ptr(compat);
    }   
#endif

    if (get_user(native, argv.ptr.native + nr))
        return ERR_PTR(-EFAULT);

    return native;
}

static int tmp_count(struct user_arg_ptr argv, int max)
{
    int i = 0;

    if (argv.ptr.native != NULL) {
        for (;;) {
            const char __user *p = get_user_arg_ptr(argv, i); 

            if (!p)
                break;

            if (IS_ERR(p))
                return -EFAULT;

            if (i >= max)
                return -E2BIG;
            ++i;

            if (fatal_signal_pending(current))
                return -ERESTARTNOHAND;
            cond_resched();
        }   
    }   
    return i;
}

asmlinkage long monitor_execve_hook(const char __user *filename, 
                          const char __user *const __user *argv,
                          const char __user *const __user *envp)
{
    int error = 0, i = 0, len = 0, offset = 0, max_len = 0;
    struct filename *path = NULL;
    const char __user * native = NULL;
    char *total_argc_ptr = NULL;
    char *total_envpc_ptr = NULL;
    char *per_envp = NULL;
    int tmp_argc = 0, total_argc_len = 0;
    int tmp_envpc = 0, total_envpc_len = 0;
    struct user_arg_ptr argvx = { .ptr.native = argv };
    struct user_arg_ptr envpx = { .ptr.native = envp };
    int nl_send_len = 0;
    struct sk_buff *skb = NULL;
    struct nlmsghdr *nlh = NULL;
    struct file *file = NULL;
    char *tmp = kmalloc(PATH_MAX, GFP_KERNEL);
    char *path1 = NULL;

    path = tmp_getname(filename);
    error = PTR_ERR(path);
    if (IS_ERR(path)) {
        goto err;
    }
    
    file = open_exec(path->name);
    if (!IS_ERR(file) && tmp) {
        memset(tmp, 0, PATH_MAX);
        path1 = d_path(&file->f_path, tmp, PATH_MAX);
        if (IS_ERR(path1)) {
            path1 = NULL;
        }
        fput(file);
    }

    error = 0;
    tmp_argc = tmp_count(argvx, MAX_ARG_STRINGS);
    if(tmp_argc < 0) {
        error = tmp_argc;
        goto err;
    }
    
    for(i = 0; i < tmp_argc; i ++) {
        native = get_user_arg_ptr(argvx, i);
        if(IS_ERR(native)) {
            error = -EFAULT;
            goto err;
        }
        
        len = strnlen_user(native, MAX_ARG_STRLEN);
        if(!len) {
            error = -EFAULT;
            goto err;
        }
        
        total_argc_len += len;
    }

    total_argc_ptr = kmalloc(total_argc_len + 16 * tmp_argc, GFP_ATOMIC);
    if(!total_argc_ptr) {
        error = -ENOMEM;
        goto err;
    }
    memset(total_argc_ptr, 0, total_argc_len + 16 * tmp_argc);

    for(i = 0; i < tmp_argc; i ++) {
        if(i == 0) {
            continue;
        }
        native = get_user_arg_ptr(argvx, i);
        if(IS_ERR(native)) {
            error = -EFAULT;
            goto err;
        }

        len = strnlen_user(native, MAX_ARG_STRLEN);
        if(!len) {
            error = -EFAULT;
            goto err;
        }
        
        if(offset + len > total_argc_len + 16 * tmp_argc) {
            break;
        }

        if (copy_from_user(total_argc_ptr + offset, native, len)) {
            error = -EFAULT;
            goto err;
        }
        offset += len - 1;
        *(total_argc_ptr + offset) = ' ';
        offset += 1;
    }
    
    /*--------envpx--------------*/
    len = 0;
    offset = 0;
    tmp_envpc = tmp_count(envpx, MAX_ARG_STRINGS);
    if(tmp_envpc < 0) {
        error = tmp_envpc;
        goto err;
    }
    
    for(i = 0; i < tmp_envpc; i ++) {
        native = get_user_arg_ptr(envpx, i);
        if(IS_ERR(native)) {
            error = -EFAULT;
            goto err;
        }
        
        len = strnlen_user(native, MAX_ARG_STRLEN);
        if(!len) {
            error = -EFAULT;
            goto err;
        }
        
        if(len > max_len) {
            max_len = len;
        }

        total_envpc_len += len;
    }
    
    per_envp = kmalloc(max_len + 16, GFP_KERNEL);
    if(!per_envp) {
        error = -ENOMEM;
        goto err;
    }

    total_envpc_ptr = kmalloc(total_envpc_len + 16 * tmp_envpc, GFP_KERNEL);
    if(!total_envpc_ptr) {
        error = -ENOMEM;
        goto err;
    }
    memset(total_envpc_ptr, 0, total_envpc_len + 16 * tmp_envpc);
    
    for(i = 0; i < tmp_envpc; i ++) {
        native = get_user_arg_ptr(envpx, i);
        if(IS_ERR(native)) {
            error = -EFAULT;
            goto err;
        }

        len = strnlen_user(native, MAX_ARG_STRLEN);
        if(!len) {
            error = -EFAULT;
            goto err;
        }
        
        if(offset + len > total_envpc_len + 16 * tmp_envpc) {
            break;
        }
        
        memset(per_envp, 0, max_len);
        if(copy_from_user(per_envp, native, len)) {
            error = -EFAULT;
            goto err;
        }
        
        if(!strstr(per_envp, "PWD") && !strstr(per_envp, "LOGNAME") && !strstr(per_envp, "USER")) {
            continue;
        }

        if (copy_from_user(total_envpc_ptr + offset, native, len)) {
            error = -EFAULT;
            goto err;
        }
        offset += len - 1;
        *(total_envpc_ptr + offset) = ' ';
        offset += 1;
    }
    
    nl_send_len = (path1 != NULL ? strlen(path1) : strlen(path->name)) + strlen(current->parent->comm) + 128;
    if(!ZERO_OR_NULL_PTR(total_envpc_ptr)) {
        nl_send_len += strlen(total_envpc_ptr);
    }   
    if(!ZERO_OR_NULL_PTR(total_argc_ptr)) {
        nl_send_len += strlen(total_argc_ptr);
    }
    nl_send_len = nl_send_len < PATH_MAX + 2048 ? nl_send_len : PATH_MAX + 2048;
    skb = alloc_skb(NLMSG_SPACE(nl_send_len), GFP_ATOMIC);
    if(!skb) {
        error = -ENOMEM;
        goto err;
    }

    nlh = (struct nlmsghdr *)skb->data;
    nlh->nlmsg_len = NLMSG_SPACE(nl_send_len);
    nlh->nlmsg_pid = 0;
    nlh->nlmsg_flags = 0;
    nlh = nlmsg_put(skb, 0, 0, 0, NLMSG_SPACE(nl_send_len) - sizeof (struct nlmsghdr), 0);
    if(!nlh) {
        kfree_skb(skb);
        pr_err("nlh get failed.\n");
        goto err;
    }

    snprintf(NLMSG_DATA(nlh), nl_send_len, "%s%c%s%c%u%c%s%c%d%c%s", path1 != NULL ? path1 : path->name, 0x1, (uint64_t)total_argc_ptr == 0x10 ? "N/A" : total_argc_ptr, 0x1, current->tgid, 0x1, current->parent->comm, 0x1, current->parent->tgid, 0x1, (uint64_t)total_envpc_ptr == 0x10 ? "N/A" : total_envpc_ptr);
    NETLINK_CB(skb).portid = 0;
    NETLINK_CB(skb).dst_group = 1;
    error = netlink_broadcast(syshook_nl_sk, skb, 0, 1, 0);
    if(error != 0 && error != -3) {
        pr_err("send nl broadcast failed.\n");
        goto err;
    }
    
    //pr_err("%s|%s|%u|%s|%d|%s\n", path->name, total_argc_ptr, current->tgid, current->parent->comm, current->parent->tgid, total_envpc_ptr);

err:
    if(tmp) {
        kfree(tmp);
        tmp = NULL;
    }
    if(per_envp) {
        kfree(per_envp);
        per_envp = NULL;
    }
    if(total_argc_ptr) {
        kfree(total_argc_ptr);
        total_argc_ptr = NULL;
    }
    if(total_envpc_ptr) {
        kfree(total_envpc_ptr);
        total_envpc_ptr = NULL;
    }

    tmp_putname(path);
    return 0;
}
#else
asmlinkage long monitor_execve_hook(void)
{
    return 0;
}
#endif

static int __init monitor_execve_init(void)
{
    int i = 0;

    if (!(sys_call_table_ptr = find_sys_call_table())){
        pr_err("Get sys_call_table failed.\n");
        return -1;
    }

#if LINUX_VERSION_CODE == KERNEL_VERSION(2, 6, 32)
    /*NetLink do not recv from userSpace*/
    syshook_nl_sk = netlink_kernel_create(&init_net, SYSHOOK_NL_NUM, 0, NULL, NULL, THIS_MODULE);
    if(!syshook_nl_sk) {
        pr_err("syshook: can not create netlink socket.\n");
        return -EIO;
    }
#elif LINUX_VERSION_CODE >= KERNEL_VERSION(3, 10, 0)
    syshook_nl_sk = netlink_kernel_create(&init_net, SYSHOOK_NL_NUM, NULL);
    if(!syshook_nl_sk) {
        pr_err("syshook: can not create netlink socket.\n");
        return -EIO;
    }
#endif
    pr_err("syshook: create netlink success.\n");


#if LINUX_VERSION_CODE >= KERNEL_VERSION(3, 10, 0)
    tmp_getname = (void *)kallsyms_lookup_name("getname");
    if(!tmp_getname) {
        pr_err("unknow Symbol getname\n");
        return -1;
    }

    tmp_putname = (void *)kallsyms_lookup_name("putname");
    if(!tmp_putname) {
        pr_err("unknow Symbol putname\n");
        return -1;
    }
#endif

    original_cr0 = read_cr0();
    write_cr0(original_cr0 & ~0x00010000);
    pr_err("Loading module monitor_execve, sys_call_table at %p\n", sys_call_table_ptr);
    
    for(i = 0; i < NR_syscalls - 1; i ++) {
        orig_sys_call_table[i] = sys_call_table_ptr[i];
    }

    orig_stub_execve = (void *)(sys_call_table_ptr[__NR_execve]);
    sys_call_table_ptr[__NR_execve]= (void *)monitor_stub_execve_hook;

    write_cr0(original_cr0);
    return 0;
}

static void __exit monitor_execve_exit(void)
{
    netlink_kernel_release(syshook_nl_sk);

    if (!sys_call_table_ptr){
        return;
    }

    write_cr0(original_cr0 & ~0x00010000);
    sys_call_table_ptr[__NR_execve] = (void *)orig_stub_execve;
    write_cr0(original_cr0);
    
    sys_call_table_ptr = NULL;
    pr_err("unload syshook_execve succ.\n");
}

module_init(monitor_execve_init);
module_exit(monitor_execve_exit);

MODULE_LICENSE("GPL");
MODULE_AUTHOR("mlsm <454667707@qq.com>");
MODULE_DESCRIPTION("Monitor Syscall sys_execve");
