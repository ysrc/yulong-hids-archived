obj-m := syshook_execve.o

syshook_execve-objs := syscall_hook.o syscall_stub_hook.o

KERNEL := /lib/modules/`uname -r`/build

all:
	make -C $(KERNEL)   M=`pwd` modules

install:
	make -C $(KERNEL)   M=`pwd` modules_install
	depmod -A

clean:
	make -C $(KERNEL)   M=`pwd` clean
