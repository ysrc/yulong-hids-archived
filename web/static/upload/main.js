
var current_step = parseInt(document.getElementsByTagName('head')[0].getAttribute('currentstep'));
var request_token = document.getElementsByTagName('head')[0].getAttribute('requesttoken');

$(document).on('ajaxSend', function (elm, xhr, settings) {
    xhr.setRequestHeader('RequestToken', request_token);
});

var filelist = [];
var has_upload_filelist = [];
var file_msg_list = [
    "请选择agent文件 (用于运行与客户端收集信息)",
    "请选择daemon文件 (agent的守护进程)",
    "请选择data文件 (agent所需的资源,zip格式)",
    "文件已全部选择"
];

var filename = [
    "windows-32",
    "windows-64",
    "linux-64"
];

function count_length(arr) {
    count = 0;
    for (a in arr) {
        count = count + 1;
    }
    console.log(count);
    return count
}

String.prototype.url_update_query = function(key, value) {
    fake_link = document.createElement('a');
    fake_link.href = this.toString();
    hash = fake_link.hash;
    if (hash) {
        uri = this.slice(0, 0-(hash.length));
    } else {
        uri = this;
    }
    if (key) {
        var re = new RegExp("([?&])" + key + "=.*?(&|$)", "i");
        var separator = uri.indexOf('?') !== -1 ? "&" : "?";
        if (uri.match(re)) {
            return uri.replace(re, '$1' + key + "=" + value + '$2') + hash;
        }
        else {
            return uri + separator + key + "=" + value + hash;
        }
    }
    return uri.toString() + hash.toString();
}

if (!localStorage.step) {
    localStorage.step = 1;
}

for (var i = 1; i < current_step; i++) {
    btn = $("button.step-" + i);
    btn.addClass("disabled");
    btn.text("你已完成该动作");
}

for (var i = current_step + 1; i <= 4; i++) {
    li = $("li.step-" + i);
    li.addClass("disabled");
}

$(document).ready(function(){
    $('#init-database').click(function(){
        $.ajax({
            type: "POST",
            url: location.href.url_update_query("step", 1),
            success: function(data) {
                if (data.status) {
                    localStorage.step = 2;
                    location.search = "?step=2";
                } else {
                    swal("Oops...", data.msg, "error");
                }
            }
        });
    });

    $('#init-rules').click(function(){
        $.ajax({
            type: "POST",
            url: location.href.url_update_query("step", 2),
            data : $('textarea#rules_text').val(),
            success: function(data) {
                if (data.status) {
                    localStorage.step = 3;
                    location.search = "?step=3";
                } else {
                    swal("Oops...", data.msg, "error");
                }
            }
        });
    });

    filename.map(function(name){
        $('#' + name).fileupload({
            dataType: 'json',
            add: function (e, data) {
                filelist[filename.indexOf(name)] = data;
                len = filelist.length;
                $("div.file-container").append("<li>" + name + ": " + data.files[0].name + " (" + data.files[0].size + " b)");
                $('#' + name).attr('disabled','disabled');
                $('#' + name).parent().attr('disabled','disabled');
                if (len >= 1 && $('button.btn.btn-warning.col-xs-8.col-xs-offset-2').length == 0) {
                    btn = $('<button/>').text('开始上传(如果需要上传多个平台的文件，请全部选择再点击上传按钮)').addClass("btn btn-warning col-xs-8 col-xs-offset-2")
                        .appendTo($('.file-upload-div'));
                    if (current_step > 3) {
                        btn.addClass("disabled");
                    } else {
                        btn.click(function () {
                            for (index in filelist) {
                                file = filelist[index];
                                file.submit();
                            }
                        });
                    }
                }
            },
            done: function (e, data) {
                $("div.file-container").append('<li>'+ name +'上传成功...</li>');
                has_upload_filelist[filename.indexOf(name)] = data;
                if (count_length(has_upload_filelist) == count_length(filelist)){
                    location.search = "?step=4";
                }
            }
        });
    });

    $('#create-key').click(function(){
        $.ajax({
            type: "POST",
            url: location.href.url_update_query("step", 4).url_update_query("action", "createkey"),
            success: function(data) {
                if (data.status) {
                    if (data.cert) {
                        $('textarea.jsoncode.init-setting').text(data.cert);
                    } else {
                        swal("无法自动化生成证书请手动生成", data.msg, "warning");
                        $('textarea.jsoncode.init-setting').attr('placeholder', "请下载私钥内容，使用'openssl req -new -x509 -key private.pem -out cert.pem -days 3650'命令生成证书文件(cert.pem)，并把 cert.pem 的文件内容粘贴到该编辑框");
                    }
                    $("div.key-dl").text('openssl req -new -x509 -key private.pem -out cert.pem -days 3650');
                    ["private"].map(function(keytype){
                        $('<a/>')
                        .attr("href", location.href.url_update_query("download", keytype))
                        .addClass("download-link")
                        .text("点击下载" + keytype)
                        .appendTo($("div.key-dl"))
                        $('#init-setting').removeAttr('disabled');
                    })
                } else {
                    swal("Oops...", data.msg, "error");
                }
            }
        });
    });

    $('#init-setting').click(function(){
        cert = $("textarea#cert").val();
        if (!/^-----BEGIN CERTIFICATE-----([\n\r]+[^-]+)+\n-----END CERTIFICATE-----/.test(cert)) {
            swal("证书格式错误", "前后记得不要有换行或者空格哦", "error");
            return
        }
        $.ajax({
            type: "POST",
            url: location.href.url_update_query("step", 4).url_update_query("action", "addconfig"),
            data: JSON.stringify({
                "ip": $("input#ipaddress").val().split(","),
                "process": $("input#process").val().split(","),
                "cert": cert,
            }),
            success: function(data) {
                if (data.status) {
                    swal({
                        title: "安装步骤已经全部完成!!",
                        text: "基础配置已上传,HTTPS证书已经更新,请重启web程序作用新的HTTPS证书",
                        type: "success",
                        confirmButtonColor: '#3085d6'
                    },
                    function(ischeck){
                        location.reload();
                    })
                } else {
                    swal("Oops...", data.msg, "error");
                }
            }
        });
    });

})