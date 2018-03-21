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
import struct
import zipfile

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
    if name_base is 'web':
            name_base = os.path.join('web', name_base)
    if is_win():
        return '{}.{}'.format(name_base, "exe")
    else:
        return name_base


def mk_build_dir():
    dir_name = 'build_{date:%Y%m%d_%H%M%S}'.format(date=datetime.datetime.now())
    path = os.path.join(base_path, 'bin', dir_name)
    print('[*] mkdir {}'.format(path))
    os.mkdir(path, mode=0o755) # rwx
    return path


def start_package_name():
    base_pkg_name = '{}-{}'
    if is_win():
        platform = 'win'
    else:
        platform = 'linux'
    arch = 8 * struct.calcsize("P")
    return base_pkg_name.format(platform, arch)


def build(build_path):
    # 复制web文件目录到编译目录下
    web_path = os.path.join(build_path, 'web')
    os.mkdir(web_path, mode=0o755)
    for dir_name in web_dir_lst:
        dir_src = os.path.join(base_path, 'web', dir_name)
        dir_dst = os.path.join(web_path, dir_name)
        shutil.copytree(dir_src, dir_dst, symlinks=False)
    # 编译各个由go(或cgo)所编写的程序
    for type_ in main_function_lst:
        go_source = main_function_lst[type_]
        if is_win():
            exe_path = os.path.join(build_path, make_execute_name(type_))
            exe_path = exe_path.replace('\\', '\\\\')
            go_source = go_source.replace('/', '\\\\')
        else:
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
    # 'app-config-sample.conf' -> 'app.conf'
    shutil.move(
        os.path.join(web_path, 'conf', 'app-config-sample.conf'),
        os.path.join(web_path, 'conf', 'app.conf')
    )
    # 生成web的压缩包
    web_zip_path = os.path.join(build_path, 'web')
    shutil.make_archive(web_zip_path, 'zip', web_path)
    # 生成当前系统的上传包
    pkg_name = start_package_name()
    mk_start_pkg(build_path, pkg_name)
    if is_win():
        pkg_name_ = 'win-64'
        mk_start_pkg(build_path, pkg_name_)
    # 生成文档的压缩包
    doc_zip_path = os.path.join(build_path, 'doc')
    shutil.make_archive(doc_zip_path, 'zip', os.path.join(base_path, 'docs'))


def mk_start_pkg(build_path, name_):
    pkg_name = os.path.join(build_path, '{}.zip'.format(name_))
    if not os.path.exists(pkg_name):
        return
    with zipfile.ZipFile(pkg_name, 'w') as myzip:
        myzip.write(os.path.join(base_path, make_execute_name('agent')))
        myzip.write(os.path.join(base_path, make_execute_name('daemon')))
        myzip.write(os.path.join(base_path, 'bin', name_, 'data.zip'))


def project_path():
    return os.path.dirname(os.path.dirname(os.path.realpath(__file__)))


def check_gopath():
    """判断是否有GOPATH这个环境变量，并判断项目是否在GOPATH内。"""
    # TODO as __doc__
    pass


if __name__ == '__main__':
    base_path = project_path()
    build_path = mk_build_dir()
    build(build_path)
