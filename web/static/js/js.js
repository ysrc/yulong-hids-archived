
Array.prototype.remove = function () {
    var what, a = arguments, L = a.length, ax;
    while (L && this.length) {
        what = a[--L];
        while ((ax = this.indexOf(what)) !== -1) {
            this.splice(ax, 1);
        }
    }
    return this;
};

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

String.prototype.url_add_Paginator = function(page, limit) {
    if (page == undefined) {
        return this.toString();
    } 
    if (limit) {
        this.url_update_query("limit", limit);
    }
    result = this.url_update_query("page", page);
    return result.toString();
}

random_hex = function () {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  
    for (var i = 0; i < 32; i++)
      text += possible.charAt(Math.floor(Math.random() * possible.length));
  
    return text;
  }

timeformat = function(datestring) {
    if (datestring == undefined) {
        return (new Date()).toLocaleString();
    }
    return (new Date(datestring)).toLocaleString();
}

/**
 * 返回当前的时间，格式：2017-10-25 15:21:57
 */
current_time = function(){
    return timeformat();
}

format_monitor_msg = function(json) {
    data = json.data
    data_type = data._type;
    source = data._source;
    return window["format_"+data_type+"_msg"](source);
}

before_1minutes = (function(){
    var d = new Date();
    return d.setMinutes(d.getMinutes() - 1);
})()

str2Date = function(string){
    var d = new Date(string);
    return d
}

json_value_to_string = function(obj, separator){
    value_list = [];
    for (key in obj) {
        if (key != "object") {
            stri = obj[key];
            if(stri.length > 10){
                stri = stri.substring(0,30);
                stri = stri + " ...";
            }
            value_list.push(stri);
        }
    }
    return value_list.join(separator);
}

var request_token = document.getElementsByTagName('head')[0].getAttribute('requesttoken');
var is_twofactorauth = document.getElementsByTagName('head')[0].getAttribute('ispass');
var httpport = document.getElementsByTagName('head')[0].getAttribute('httpport');
var iswatchmode = document.getElementsByTagName('head')[0].getAttribute('iswatchmode');
var apihost = document.getElementsByTagName('head')[0].getAttribute('apihost');

var devmode = false
var login_url = apihost + '/login'
var api_base_url = apihost + '/json'

var client_url = api_base_url + "/client"
var info_url = api_base_url + "/info"
var notice_url = api_base_url + "/notice"
var config_url = api_base_url + "/config"
var task_url = api_base_url + "/tasks"
var file_url = api_base_url + "/file"
var monitor_url = api_base_url + "/monitor"
var analyze_url = api_base_url + "/analyze"
var statistics_url = api_base_url + "/statistics"
var rules_url = api_base_url + "/rules"
var logout_url = api_base_url + "/logout"

if (!localStorage.search_history) {
    localStorage.search_history = "";
}

if (!localStorage.host_filter_history) {
    localStorage.host_filter_history = "";
}

var lang = "zh_cn"; // 默认语言为简体中文

var hostw = angular.module(
    'hostw',
    ['ngRoute', 'ui-notification', 'moment-picker'],
    function ($interpolateProvider) {}
);

hostw.run(function($rootScope) {
    $rootScope.timeformat = timeformat;
});

hostw.filter('cut', function () {
    return function (value, wordwise, max, tail) {
        if (!value) return '';

        max = parseInt(max, 10);
        if (!max) return value;
        if (value.length <= max) return value;

        value = value.substr(0, max);
        if (wordwise) {
            var lastspace = value.lastIndexOf(' ');
            if (lastspace !== -1) {
              if (value.charAt(lastspace-1) === '.' || value.charAt(lastspace-1) === ',') {
                lastspace = lastspace - 1;
              }
              value = value.substr(0, lastspace);
            }
        }

        return value + (tail || ' …');
    };
});

hostw.service('XSRFInterceptor', [function () {
    var service = this;
    service.request = function (config) {
        config.headers.RequestToken = request_token;
        return config;
    };
}]);

hostw.factory('Redirect', function ($q) {
    return {
        response: function (response) {
            if (response.headers()['is-login-page'] == 'true') {
                location.href = login_url;
                return response;
            }
            if(response.headers()['content-type'] != "application/json; charset=utf-8"){
                swal('响应格式错误!', '服务端响应格式错误，请检查输入是否合理', 'error');
                return response;
            }
            return response;
        }
    };
});

hostw.filter('cutWords', function() {

    return function(input, length) {
        if (input) {
            if (input.length > length) {
                return input.substring(0, length) + ' ......'
            }
            return input;
        }
        return input;
    }

});

// 适应前后端分离的devmode
if (devmode) {
    hostw.service('devmode', [function () {
        var service = this;
        service.request = function (config) {
            config.headers["Content-Type"] = "application/x-www-form-urlencoded";
            return config;
        };
    }]);
}

if (!window.HostWData) {
    HostWData = {
        update: function (name, dict) {
            if (!HostWData[name]) {
                HostWData[name] = {};
            }
            for (key in dict) {
                HostWData[name][key] = dict[key];
            }
            if (HostWData.onload) HostWData.onload();
        }
    }
};

change_switchery_status = function(selector){
    ele = $(selector);
    ele.siblings().remove();
    ele.prop("checked", !ele.is(':checked'));
    return new Switchery(document.querySelector(selector), {size: 'small'});
}

add_history = function(key, query){
    history_list = localStorage[key].split(":::");
    for (h in history_list) {
        if (history_list[h] == query) {
            return
        }
    }
    if (history_list.length >= 6) {
        history_list.shift();
    }
    history_list.push(query);
    localStorage[key] = history_list.join(":::")
}

function escape_regexp(str) {
    return str.replace(/[\-\[\]\/\{\}\(\)\*\+\?\.\\\^\$\|]/g, "\\$&");
}

function string2regexp(str) {
    return "^" + escape_regexp(str) + "$";
}

function request_password(callback){
    if (is_twofactorauth == "true") {
        setTimeout(() => {
            swal({
                title: "请输入双因子验证密码",
                text: "您已经开启了双因子验证，请输入双因子验证密码",
                type: "input",
                showCancelButton: true,
                closeOnConfirm: true
            },
            function (inputValue) {
                if (inputValue) {
                    callback(inputValue)
                }
            });
        }, 500);
    } else {
        callback(0);
    }
}

function enable_pass(url, pass) {
    return url.url_update_query("pass", pass);
}

function modal_option_click() {
    enable_command_type = ['kill', 'delete', 'exec'];
    $($('div#newTaskModel input')[1]).click(function(){
        type_ = $('div#newTaskModel select').val();
        if (enable_command_type.indexOf(type_) < 0) {
            swal("Oops!", type_ + "类型的任务无需填写command", "warning");
        }
    })
}

$(document).ready(function(){
    ele_lst = $('.sidebar-menu .sidebar-nav a');
    ele_lst.each(function(a_ele){
        ele = $(ele_lst[a_ele]);
        ele.click(function(){
            nohref = $(ele_lst[a_ele]).attr('no-href');
            if (nohref) {
                if (location.href.endsWith(nohref)) {
                    location.reload();
                } else {
                    location.replace(nohref);
                }
            }
         });
    });
});
