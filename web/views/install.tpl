<!DOCTYPE html>
<html>
<head currentstep="<<< .currentstep >>>" requesttoken="<<< .token >>>">
  <title>主页</title>
  <meta charset="utf-8">
  <!-- ext -->
  <link rel="stylesheet" type="text/css" href="/static/css/ext-bootstrap-3-3-7.css">
  <link rel="stylesheet" type="text/css" href="/static/css/ext-font-awesome.css">
  <link rel="stylesheet" type="text/css" href="/static/css/sweetalert.min.css">
  <link rel="stylesheet" type="text/css" href="/static/css/ysrc.css">
  <script src="/static/js/sweetalert.min.js"></script>
  <script src="/static/js/jquery.js"></script>
  <script src="/static/upload/jquery.ui.widget.js"></script>
  <script src="/static/upload/jquery.iframe-transport.js"></script>
  <script src="/static/upload/jquery.fileupload.js"></script>
  <script src="/static/upload/jquery.fileupload-process.js"></script>  
</head>
<body>

<div class="out">
    <div class="container">
        <div class="row form-group">
            <div class="col-xs-12">
                <ul class="nav nav-pills nav-justified thumbnail setup-panel">
                    <li class="<<< if eq 1 .step >>>active<<< end >>> step-1"><a href="?step=1">
                        <h4 class="list-group-item-heading">Step 1</h4>
                        <p class="list-group-item-text">数据库初始化</p>
                    </a></li>
                    <li class="<<< if eq 2 .step >>>active<<< end >>> step-2"><a href="?step=2">
                        <h4 class="list-group-item-heading">Step 2</h4>
                        <p class="list-group-item-text">初始化规则</p>
                    </a></li>
                    <li class="<<< if eq 3 .step >>>active<<< end >>> step-3"><a href="?step=3">
                        <h4 class="list-group-item-heading">Step 3</h4>
                        <p class="list-group-item-text">上传文件</p>
                    </a></li>
                    <li class="<<< if eq 4 .step >>>active<<< end >>> step-4"><a href="?step=4">
                        <h4 class="list-group-item-heading">Step 4</h4>
                        <p class="list-group-item-text">写入基础配置</p>
                    </a></li>
                </ul>
            </div>
        </div>
        <div class="row setup-content" id="step-1" style="display: <<< if ne 1 .step >>> none <<< end >>>" >
            <div class="col-xs-12">
                <div class="col-md-12 well text-center">
                    <h1> 数据库初始化 </h1>
                    <button id="init-database" class="btn btn-primary btn-lg step-1">
                        初始化数据库
                    </button>
                </div>
            </div>
        </div>
        <div class="row setup-content" id="step-2" style="display: <<< if ne 2 .step >>> none <<< end >>>">
            <div class="col-xs-12">
                <div class="col-md-12 well text-center">
                    <h1> 初始化规则 </h1>
                    <textarea class="jsoncode" id="rules_text" placeholder="此处不得为空，请查看文档填写, 默认规则请粘贴rules.json文件内容"></textarea>
                    <button id="init-rules" class="btn btn-primary btn-lg step-2">
                        初始化规则
                    </button>
                </div>
            </div>
        </div>
        <div class="row setup-content" id="step-3" style="display: <<< if ne 3 .step >>> none <<< end >>>">
            <div class="col-xs-12">
                <div class="col-md-12 well text-center file-upload-div">
                    <h1> 文件上传 </h1>
                    <p>你已经选择了以下文件:</p>
                    <div class="file-container">
                    </div>

                    <div class="file-uploader btn btn-primary">
                        <span> <i class="fa fa-windows" aria-hidden="true"></i> 
                            <span>
                                Windows 64位压缩包
                            </span>
                        </span>
                        <input type="file" class="fileele" id="windows-64" name="file" data-url="?step=3&system=windows&platform=64"/>
                    </div>

                    <div class="file-uploader btn btn-primary">
                        <span> <i class="fa fa-windows" aria-hidden="true"></i>
                            <span>Windows 32位压缩包</span>
                        </span>
                        <input type="file" class="fileele" id="windows-32" name="file" data-url="?step=3&system=windows&platform=32"/>
                    </div>

                    <div class="file-uploader btn btn-success">
                        <span> <i class="fa fa-linux" aria-hidden="true"></i>
                            <span>Linux 64位压缩包</span>
                        </span>
                        <input type="file" class="fileele" id="linux-64" name="file" data-url="?step=3&system=linux&platform=64"/>
                    </div>

                </div>
            </div>
        </div>

        <div class="row setup-content" id="step-4" style="display: <<< if ne 4 .step >>> none <<< end >>>">
            <div class="col-xs-12">
                <div class="col-md-12 well text-center">
                    <h1> 写入基础配置文件 </h1>
                    <div class="form-group">
                        <input class="form-control" placeholder="不纳入监控的IP或运维、管理人员常用IP，多个IP以逗号分开" id="ipaddress">
                    </div>
                    <div class="form-group">
                        <input class="form-control" placeholder="不纳入监控的程序或其他agent等，该项为正则，多个正则以逗号分开" id="process">
                    </div>
                    <div class="key-dl">
                    </div>
                    <textarea class="jsoncode init-setting" placeholder="此处填写证书文件内容。请先点击‘生成证书’按钮。" id="cert"></textarea>
                    <button id="create-key" class="btn btn-primary btn-lg step-4">
                        生成证书
                    </button>
                    <button id="init-setting" class="btn btn-primary btn-lg step-4" disabled>
                        完成
                    </button>
                </div>
            </div>
        </div>

    </div>
</div>
</body>

<style>
    div.out {
        margin-top: 5%;
    }
    textarea.jsoncode {
        height: 400px;
        overflow-y: scroll;
    }
    textarea.jsoncode.init-setting {
        height: 180px;
    }
    .bar {
        height: 18px;
        background: green;
    }
    div.select-file {
        margin-top: 10px;
    }

    button.btn.btn-warning.col-xs-8.col-xs-offset-2 {
        margin-top: 10px;
    }
    
    a.download-link {
        margin-left: 10px;
    }

</style>

<script src="/static/upload/main.js"></script>

</html>