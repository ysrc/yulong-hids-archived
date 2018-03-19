#!/usr/bin/env python3
# coding=utf-8
# by nearg1e (nearg1e.com@gmail[dot]com)
"""
驭龙hids的编译脚本，一键打包成需要的 release。

该脚本会在 bin 目录下生成一个 release 文件夹, 文件名为:release_+可读的时间字符串。
release 目录下文件如下:

├── agent|agent.exe    当前平台的agent可执行文件
├── daemon|daemon.exe  当前平台的daemon可执行文件
├── server|server.exe  当前平台的server可执行文件
├── rules.json         默认规则，供web安装时上传
├── web                web文件文件夹，包含web可执行文件和静态文件
├── linux-64.zip       包含着agent，daemon和依赖文件的压缩包，供web安装时上传
├── doc.zip            项目文档
└── web.zip            web文件文件夹的压缩包

"""

import os
import datetime
import subprocess
import shlex
import shutil

base_path = ''
base_command = 'go build -o {out} --ldflags="-w -s" {source}'
main_function_lst = {
    'agent': 'agent/agent.go',
    'daemon': 'daemon/daemon.go',
    'server': 'server/server.go',
    'web': 'web/main.go'
}
web_dir_lst = ['conf', 'https_cert', 'static', 'upload_files', 'views']


def is_win():
    return os.name == 'nt'


def make_execute_name(name_base):
    if is_win():
        return '{}.{}'.format(name_base, "exe")
    else:
        return name_base


def mk_build_dir():
    dir_name = 'build_{date:%Y%m%d_%H%M%S}'.format(date=datetime.datetime.now())
    path = os.path.join(base_path, 'bin', dir_name)
    os.mkdir(dir_name, mode=755) # rwx
    return path


def build(build_path):
    # 编译各个由go(或cgo)所编写的程序
    for type_ in main_function_lst:
        go_source = main_function_lst[type_]
        exe_path = os.path.join(build_path, make_execute_name(type_))
        command = base_command.format(out=exe_path, source=go_source)
        print('[*] run command:', command)
        command_args = shlex.split(command)
        output = subprocess.check_output(
            command_args,
            cwd=base_path,
            stderr=subprocess.STDOUT
        )
        print('[*] stdout & stderr:', output.decode())
    # 移动默认规则到编译目录下
    rule_src = os.path.join(base_path, 'default_rules.json')
    rule_dst = os.path.join(build_path, 'rules.json')
    shutil.copyfile(rule_src, rule_dst, follow_symlinks=False)
    # 复制web文件目录到编译目录下
    web_path = os.path.join(build_path, 'web')
    os.mkdir(web_path, mode=755)
    for dir_name in web_dir_lst:
        dir_src = os.path.join(base_path, 'web', dir_name)
        dir_dst = os.path.join(web_path, dir_name)
        shutil.copytree(dir_src, dir_dst, symlinks=False)
    # 'app-config-sample.conf' -> 'app.conf'
    shutil.move(
        os.path.join(web_path, 'conf', 'app-config-sample.conf'),
        os.path.join(web_path, 'conf', 'app.conf')
    )
    shutil.copyfile(
        os.path.join(build, make_execute_name('web')),
        os.path.join(web_path, make_execute_name('web')),
        follow_symlinks=False
    )


def project_path():
    return os.path.dirname(os.path.dirname(os.path.realpath(__file__)))


def check_gopath():
    """判断是否有GOPATH这个环境变量，并判断项目是否在GOPATH内。"""
    # TODO as __doc__
    pass


if __name__ == '__main__':
    base_path = project_path()
    build_path = mk_build_dir()
