
hostw.service('globalVariable', function () {
    var notice_search = '';
    return {
        notice_search: notice_search
    };
});

hostw.config(function ($routeProvider, $httpProvider) {

    $httpProvider.interceptors.push('Redirect');

    if (!devmode) {
        $httpProvider.interceptors.push('XSRFInterceptor');
        $(document).on('ajaxSend', function (elm, xhr, settings) {
            xhr.setRequestHeader('requesttoken', request_token);
        });
    } else {
        $httpProvider.interceptors.push('devmode');
        $(document).on('ajaxSend', function (elm, xhr, settings) {
            if (settings.method == "post") {
                xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
            }
        });
    }

    $('div.icon.logout').click(function(){
        swal({
            title: "确定要退出登录吗?",
            type: 'warning',
            showCancelButton: true,
            closeOnConfirm: true
        },
        function () {
            $.post(logout_url, function(data){
                if (data.status) {
                    location.replace(login_url);
                }
            });
        });
    })

    $routeProvider.when("/", {
        controller: hostw.statistics,
        template: document.getElementById('statistics').text
    }).when("/info/:ip/", {
        controller: hostw.detailinfo,
        template: document.getElementById('detailinfo').text
    }).when("/notice", {
        controller: hostw.notice,
        template: document.getElementById('notice').text
    }).when("/config/", {
        controller: hostw.config,
        template: document.getElementById('config').text
    }).when("/tasks", {
        controller: hostw.tasks,
        template: document.getElementById('tasks').text
    }).when("/taskresult/:id/", {
        controller: hostw.taskresult,
        template: document.getElementById('taskresult').text
    }).when("/analyze", {
        controller: hostw.analyze,
        template: document.getElementById('analyze').text
    }).when("/rules", {
        controller: hostw.rules,
        template: document.getElementById('rules').text
    }).when("/hostlist", {
        controller: hostw.hostlist,
        template: document.getElementById('hostlist').text
    }).otherwise({
        redirectTo: "/"
    });
});

hostw.controller('analyze', function ($scope, $http, Notification, globalVariable) {
    var input = document.getElementById("search");

    $scope.history = localStorage.search_history.split(":::");
    $scope.re_search = function (history) {
        $("input#search")[0].value = history;
    }

    request_and_show = function (url) {
        $http.get(url).then(function (response) {
            $scope.analyze = response.data;
            if (response.data) {
                $scope.tasks = response.data;
            } else {
                Notification.error("已经到达最末页。");
            }
            $scope.thead_root_items = [];
            $scope.thead_data_items = [];
            $scope.data_is_array = false;
            $scope.json_value_to_string = json_value_to_string;
            $scope.show_sub_index = -1;
            $scope.sub_td_data = [];

            (function () {
                for (k in $scope.analyze[0]) {
                    if ((!k.startsWith('_')) && (k != "data")) {
                        $scope.thead_root_items.push(k);
                    }
                }
                if ($scope.analyze[0].data instanceof Array) {
                    for (k in $scope.analyze[0].data[0]) {
                        $scope.thead_data_items.push(k);
                    }
                    $scope.data_is_array = true;
                    $scope.data_index = 0;
                } else {
                    for (k in $scope.analyze[0].data) {
                        $scope.thead_data_items.push(k);
                    }
                    $scope.data_is_array = false;
                }
            })();
        });
    }

    $scope.click_more = function (index) {
        $scope.sub_td_data = $scope.analyze[index].data;
    };

    // paginator design
    $scope.current_page = 1;
    $scope.previous = function () {
        if ($scope.current_page > 1) {
            $scope.current_page = $scope.current_page - 1;
            $scope.get_result();
        }
    }
    $scope.next = function () {
        $scope.current_page = $scope.current_page + 1;
        $scope.get_result();
    }

    $scope.show_data = function () {
        $scope.current_page = 1;
        $scope.get_result();
    }

    $scope.append2search = function (end) {
        input = $("input#search");
        input[0].value = input[0].value + end;
        input.focus();
    }

    // get analyze result
    $scope.get_result = function () {
        value = $("input#search")[0].value;
        if (value == "") {
            return
        }
        add_history('search_history', value);

        query = encodeURIComponent(value);
        url = analyze_url
            .url_update_query("q", query)
            .url_add_Paginator($scope.current_page);

        if ($scope.ctrl != undefined && $scope.ctrl.minday != undefined) {
            minday = Date.parse($scope.ctrl.minday);
            if ($scope.ctrl.maxday == undefined) {
                maxday = "now";
            } else {
                maxday = Date.parse($scope.ctrl.maxday);
            }

            time_query = encodeURIComponent(minday + "-" + maxday);
            url = url.url_update_query("tq", time_query);
        }

        request_and_show(url);
    }

    if (globalVariable.notice_search) {
        $("input#search")[0].value = globalVariable.notice_search;
        globalVariable.notice_search = "";
        $scope.get_result();
    }

    $scope.search_start = $("input#search")[0].value;

    $('input#search').on('keydown', function (e) {
        $scope.search_sample = "";
        $scope.search_start = $("input#search")[0].value;
        if (e.keyCode == 9) {
            e.preventDefault();
            search = $(input).val();
            last = search.split("|").pop();
            if (last.endsWith(':')) {
                $.ajax({
                    url: analyze_url,
                    type: "POST",
                    data: {
                        keyword: last
                    },
                    success: function (data) {
                        $scope.search_sample = data;
                        $scope.$apply();
                    }
                });
            } else if (last == "") {
                res = /type\s*\:\s*(\w+)[\|\b]/.exec(search);
                console.log(res);
                if (res) {
                    type_text = res[1];
                    $scope.search_sample = {
                        "msg": "",
                        "samplelist": Object.keys(global_data_info[type_text])
                    }
                    $scope.$apply();
                }
            }
        } else if (e.keyCode == 13) {
            $scope.show_data();
        }
    });
});

hostw.controller('tasks', function ($scope, $http, Notification) {
    modal_option_click();
    $scope.table_key_list = ['name', 'type', 'command', 'host_list', 'time'];

    $scope.current_page = 1;
    $scope.previous = function () {
        if ($scope.current_page > 1) {
            $scope.current_page = $scope.current_page - 1;
            $scope.get_result();
        }
    }
    $scope.next = function () {
        $scope.current_page = $scope.current_page + 1;
        $scope.get_result();
    }

    $scope.get_result = function () {
        $http.get(task_url.url_add_Paginator($scope.current_page)).then(
            function (response) {
                if (response.data) {
                    $scope.tasks = response.data;
                    Notification.success("成功加载 " + response.data.length + " 条数据。");
                } else {
                    Notification.error("没有其它数据了。");
                }
            });
    }

    $scope.get_result();

});

hostw.controller('taskresult', function ($scope, $http, Notification, $rootScope, $routeParams) {
    modal_option_click();
    // paginator design
    $scope.table_key_list = ['ip', 'status', 'data', 'time'];
    $scope.task_id = $routeParams.id;

    $scope.current_page = 1;
    $scope.previous = function () {
        if ($scope.current_page > 1) {
            $scope.current_page = $scope.current_page - 1;
            $scope.get_result();
        }
    }
    $scope.next = function () {
        $scope.current_page = $scope.current_page + 1;
        $scope.get_result();
    }

    $scope.get_result = function () {
        $http.get(task_url.url_add_Paginator($scope.current_page).url_update_query('tid', $routeParams.id)).then(
            function (response) {
                if (response.data) {
                    $scope.taskresult = response.data;
                    Notification.success("成功加载 " + response.data.length + " 条数据。");
                } else {
                    Notification.error("没有其它数据了。");
                }
            });
    }

    $scope.get_result();

});

hostw.controller('file', function ($scope, $http) {
    $http.get(file_url).then(
        function (response) {
            $scope.upload_text = "upload"
            $scope.files = response.data
        });

    change = function (index) {
        console.log($('#file' + index)[0].files[0])
        var formData = new FormData();
        formData.append("system", $scope.files[index].system)
        formData.append("platform", $scope.files[index].platform)
        formData.append("type", $scope.files[index].type)
        formData.append('file', $('#file' + index)[0].files[0]);
        request_password(function (password) {
            file_url = file_url.url_update_query('pass', password);
            $.ajax({
                url: file_url,
                type: 'PUT',
                cache: false,
                data: formData,
                processData: false,
                contentType: false
            }).done(function (res) {
                ajaxcallback(res)
            }).fail(function (res) {
                ajaxcallback(res)
            });
        });
    };

    ShowUpload = function () {
        $("#con-close-modal").modal();
    }
    changeFile = function () {
        var PLATFORM = $("input[name='PLATFORM']:checked").val();
        var SYSTEM = $("input[name='SYSTEM']:checked").val();
        var TYPE = $("input[name='TYPE']:checked").val();
        console.log(PLATFORM + SYSTEM + TYPE);

        var formData = new FormData();
        formData.append("system", SYSTEM)
        formData.append("platform", PLATFORM)
        formData.append("type", TYPE)
        formData.append('file', $('#fileUpload')[0].files[0]);
        request_password(function (password) {
            file_url = file_url.url_update_query('pass', password);
            $.ajax({
                url: file_url,
                type: 'PUT',
                cache: false,
                data: formData,
                processData: false,
                contentType: false
            }).done(function (res) {
                $("#con-close-modal").modal('hide')
                ajaxcallback(res)
            }).fail(function (res) {
                ajaxcallback(res)
            });
        });
    }
});

hostw.controller('rules', function ($scope, $http, Notification) {
    $http.get(rules_url).then(function (response) {
        $scope.rules_url = rules_url;
        $scope.rulelist = response.data;
        $scope.rulelist.sort(function(item1, item2){
            name1 = item1.meta.name;
            name2 = item2.meta.name;
            if (name1 > name2) return -1;
            if (name1 < name2) return 1;
            return 0;
        });
        $scope.langtem = HostWData[lang];
        $scope.style = HostWData.style;

        rule_template = {
            "and": true,
            "enabled": false,
            "meta": {
                "author": "",
                "description": "",
                "level": 0,
                "name": ""
            },
            "rules": {
                "name": {
                    "data": "",
                    "type": ""
                }
            },
            "source": "",
            "system": ""
        }

        $scope.edit = function(index) {
            current = jQuery.extend({}, $scope.rulelist[index]);
            // current['//'] = '点击编辑规则并保存并不会覆盖原本规则，请调试新规则无误之后再删除原本规则。';
            $scope.new_rules = $scope.pettyprint(current, null, 2);
            $('[href="#model-add-rules"]').click();
            $('h4.modal-title strong').text(
                "编辑规则 (点击编辑规则并保存并不会覆盖原本规则, 而会新增一条规则，请调试新规则无误之后再删除原本规则)"
            );
        }

        $scope.pettyprint = function (obj, ...args) {
            res = Object.assign({}, obj);
            res._id = undefined;
            return JSON.stringify(res, ...args);
        }

        $scope.delete = function (id) {
            obj = { "id": id };
            swal({
                title: "删除操作",
                text: "该动作会删除这条规则，且不可复原。",
                showCancelButton: true,
                type: "warning",
                confirmButtonColor: "#DD6B55"
            },
            function () {
                request_password(function (password) {
                    rules_url = rules_url.url_update_query('pass', password);
                    $http.post(
                        rules_url.url_update_query('action', 'del'),
                        obj
                    ).then(function (response) {
                        res = response.data;
                        if (res.status) {
                            Notification.success('成功删除规则!');
                            location.reload();
                        }
                        if (response.data.status == false) {
                            ajaxcallback(response.data);
                        }
                    })
                })
            });
        }

        $scope.new_rules = $scope.pettyprint(rule_template, null, 2);

        $scope.open_body = function (id) {
            $('#panel-' + id).toggle();
        }
        setTimeout(function () {
            if ($scope.rulelist.length > 0) {
                $scope.open_body($scope.rulelist[0]._id);
            }
        }, 10);

        $scope.add_rules = function () {
            obj = JSON.parse($scope.new_rules);
            if (!(obj instanceof Array)) {
                obj = [obj];
            }
            swal({
                title: "是否确定添加规则？",
                text: "该动作会添加 [" + (obj.length ? obj.length : 1) + "] 条规则",
                showCancelButton: true,
                type: "warning",
                confirmButtonColor: "#DD6B55"
            },
                function () {
                    request_password(function (password) {
                        rules_url = rules_url.url_update_query('pass', password);
                        $http.post(
                            rules_url.url_update_query('action', 'add'),
                            obj
                        ).then(function (response) {
                            res = response.data;
                            if (res.length) {
                                Notification.success('成功添加' + res.length + '条规则!');
                            }
                            if (response.data.status == false) {
                                ajaxcallback(response.data);
                            }
                        })
                    })
                }
            )
        }

        $scope.enable_toggle = function (index) {
            current = $scope.rulelist[index];
            current_status = current.enabled;
            word_dict = {
                true: "关闭",
                false: "开启"
            }
            swal({
                title: "是否确定" + word_dict[current_status] + "该规则？",
                text: "该动作会" + word_dict[current_status] + "\"" + current.meta.name + "\"规则",
                showCancelButton: true,
                type: "warning",
                showCancelButton: true,
                confirmButtonColor: "#DD6B55"
            },
                function () {
                    request_password(function (password) {
                        rules_url = rules_url.url_update_query('pass', password);
                        $http.post(
                            rules_url.url_update_query('action', 'enable'),
                            { "id": current._id, "enable": !current_status }
                        ).then(function (response) {
                            res = response.data;
                            if (res.enable != undefined) {
                                $scope.rulelist[index].enabled = res.enable;
                                location.reload();
                                // $scope.$apply();
                            }
                            if (response.data.status == false) {
                                ajaxcallback(response.data);
                            }
                        })
                    })
                }
            );
        }
    });
});

hostw.controller('config', function ($scope, $http, $rootScope, $routeParams, Notification) {
    $http.get(config_url).then(function (response) {
        $scope.langtem = HostWData[lang];
        $scope.configs = response.data;
        $scope.type = $routeParams.type
        $scope.isArray = Array.isArray
        $scope.typeof = function (var_) {
            return typeof (var_)
        }
        $scope.type = $scope.configs[0].type
        $scope.white_index = null;
        $scope.white_list_switchery = null;
        $scope.disable_edit_list = ["cert", "privatekey", "publickey"];

        $scope.config_first_child = function(type, key){
            dic = $scope.langtem.config[type];
            if (dic == undefined ) {
                return key
            }
            word = dic[key];
            if (word == undefined ) {
                return key
            }
            split_lst = word.split(' ');
            return split_lst[0];
        }

        // add html5tooltips
        $(document).ready(function(){
            language_dict = $scope.langtem.config;
            Object.keys(language_dict).forEach(function (config_type) {
                html5tooltips({
                    animateFunction: "spin",
                    color: "sky",
                    contentText: language_dict[config_type].type_description,
                    stickTo: "bottom",
                    targetSelector: "#"+config_type+'_tab'
                });
                
                Object.keys($scope.langtem.config[config_type]).forEach(function(config_sub_key){
                    selector = '#'+config_type + ' #'+config_sub_key + ' div.right i';
                    html5tooltips({
                        animateFunction: "foldin",
                        color: "lilac",
                        contentText: language_dict[config_type][config_sub_key],
                        stickTo: "right",
                        targetSelector: selector
                    });
                });
            });
        });

        (function () {
            for (i_ in $scope.configs) {
                config = $scope.configs[i_];
                if (config.type == "whitelist") {
                    $scope.white_index = i_;
                    return
                }
            }
        })();

        $scope.edit = function (index, key) {
            origin_value = $('#'+$scope.configs[index].type).find('#'+key).find('pre.config-text').text();
            get_input(title = key.toUpperCase(), text = 'Input a new value!', origin_value, function (inputValue) {
                if (inputValue === "") {
                    swal.showInputError("You need to write something!");
                    return false
                }
                if (inputValue) {
                    json = {
                        id: $scope.configs[index]._id,
                        key: key,
                        input: inputValue
                    }
                    currenttab = $('.card-header ul li.active').attr('id')
                    edit_config_item(json, currenttab = currenttab)
                }
            })
        }

        $scope.add = function (index, key) {
            get_input(title = key.toUpperCase(), text = 'Add a new item!', "", function (inputValue) {
                if (inputValue === "") {
                    swal.showInputError("You need to write something!");
                    return false
                }
                if ($scope.configs[index]['dic'][key].indexOf(inputValue) != -1) {
                    swal.showInputError("Your input already exists in " + key);
                    return false
                }
                if (inputValue) {
                    json = {
                        id: $scope.configs[index]._id,
                        key: key,
                        input: inputValue
                    }
                    request_password(function (password) {
                        config_url = config_url.url_update_query('pass', password);
                        req = {
                            method: 'PUT',
                            url: config_url,
                            data: json
                        }
                        $http(req).then(function (response) {
                            if (response.data.status) {
                                $scope.configs[index]['dic'][key].push(inputValue);
                            }
                            ajaxcallback(response.data);
                        })
                    });
                }
            })
        }

        $scope.del = function (event, id, index, k) {
            swal({
                title: "Are you sure?",
                text: "Delete the \"" + ele.text().trim() + "\"?",
                type: "warning",
                showCancelButton: true,
                confirmButtonColor: "#DD6B55",
                confirmButtonText: "comfirm!",
                closeOnConfirm: false
            },
            function () {
                input = ele.text().trim()
                ele = $(event.currentTarget)
                request_password(function (password) {
                    config_url = config_url.url_update_query('pass', password);
                    req = {
                        method: 'DELETE',
                        url: config_url,
                        data: { input: input, key: k, id: id }
                    }
                    $http(req).then(function (response) {
                        ajaxcallback(response.data);
                        if (response.data.status) {
                            $scope.configs[index]['dic'][k].remove(input)
                        }
                    })
                });
            });
        }

        $scope.showdel = function (event) {
            ele = $(event.currentTarget)
            ele.addClass('badge-danger').removeClass('badge-success')
            ele.find('i').addClass('fa-times').removeClass('fa-check')
        }
        $scope.hidedel = function (event) {
            ele = $(event.currentTarget)
            ele.addClass('badge-success').removeClass('badge-danger')
            ele.find('i').addClass('fa-check').removeClass('fa-times')
        }
        $scope.showadd = function (event) {
            ele = $(event.currentTarget)
            ele.addClass('badge-primary').removeClass('badge-warning')
        }
        $scope.hideadd = function (event) {
            ele = $(event.currentTarget)
            ele.addClass('badge-warning').removeClass('badge-primary')
        }
        $scope.add_white = function (event, type, info, count) {

            if (count == 1) {
                $.ajax({
                    url: notice_url,
                    method: "post",
                    dataType: "json",
                    data: JSON.stringify({
                        id: "learn",
                        type: type,
                        info: info
                    }),
                    success: function (data) {
                        if (data.status) {
                            ele = $(event.currentTarget);
                            ele.addClass('badge-success').removeClass('badge-primary');
                            ele.find('i.count').text('已忽略');
                            Notification.success($scope.langtem.notice.data.type[type] + ':' + info + ' 已经忽略');
                        }
                    }
                });
                return
            }

            keydict = {
                "process": "process",
                "loginlog": "ip",
                "connection": "ip"
            }
            info_ = info;

            if (type == 'loginlog') {
                ara = /[\d\.]+/.exec(info);
                if (ara) {
                    info_ = ara[0];
                }
            }

            if (type == "process") {
                info_ = string2regexp(info_)
            }

            addw_json = {
                id: $scope.configs[$scope.white_index]._id,
                key: keydict[type],
                input: info_
            }

            req = {
                method: 'PUT',
                url: config_url,
                data: addw_json
            }

            $http(req).then(function (response) {
                if (response.data.status) {
                    ele = $(event.currentTarget);
                    ele.addClass('badge-success').removeClass('badge-primary');
                    ele.find('i.count').text('已添加');
                    Notification.success($scope.langtem.notice.data.type[type] + ':' + info + ' 已添加到白名单');
                    delete_();
                }
                if (response.data.status == false) {
                    ajaxcallback(response.data);
                }
            });

            function delete_() {
                req_del = {
                    method: 'DELETE',
                    url: notice_url.url_update_query('type', type).url_update_query('info', info)
                }
                $http(req_del).then(function (response) {
                    if (response.data.status) {
                        Notification.warning('删除已经添加到白名单的信息:' + info);
                    }
                });
            }

        }

        $scope.learn_data_fa = function (type) {
            json = {
                "process": "fa-file",
                "loginlog": "fa-sign-in",
                "connection": "fa-plug"
            }
            return json[type]
        }

        $scope.clean_notice = function () {
            data_ = {
                id: 'all_learn',
                status: 0
            }
            $http({
                method: 'POST',
                url: notice_url,
                data: data_
            }).then(function (response) {
                if (response.data.status) {
                    Notification.warning('所有在观察模式下的未处理告警已经设置为未处理告警。')
                }
            });
        }

        $scope.UpdateBool = function (index, key) {
            var inputValue = $("#chk" + $scope.configs[index].type + "-" + key).is(':checked') == true ? 'true' : 'false';
            if (inputValue == 'false' && key == "learn") {
                $scope.white_list_switchery.setPosition(true);
                $http.get(notice_url.url_update_query('status', 'learn')).then(function (response) {
                    $scope.learn_data = response.data;
                });
                $('#model-select-learn').modal('toggle');

                $scope.close_learn = function () {
                    json = {
                        id: $scope.configs[index]._id,
                        key: key,
                        input: "false"
                    }
                    currenttab = $('.card-header ul li.active').attr('id');
                    change_switchery_status("#chk" + $scope.configs[index].type + "-" + key);
                    edit_config_item(json, currenttab = currenttab);

                    $('#model-select-learn').modal('toggle');
                }
                return
            }
            if (inputValue) {
                json = {
                    id: $scope.configs[index]._id,
                    key: key,
                    input: inputValue
                }
                currenttab = $('.card-header ul li.active').attr('id')
                edit_config_item(json, currenttab = currenttab)
            }
        }

        $scope.saveInfo = function (index, key) {
            var inputValue = $("#span_" + key).find("pre").text();
            if (inputValue) {
                json = {
                    id: $scope.configs[index]._id,
                    key: key,
                    input: inputValue
                }
                currenttab = $('.card-header ul li.active').attr('id')
                edit_config_item(json, currenttab = currenttab)
                $("#save_" + key).css("display", "none");
                $("#edit_" + key).css("display", "inline-block");
            }
        }
        $(document).ready(function () {
            var elems = Array.prototype.slice.call(document.querySelectorAll('.js-switch'));
            elems.forEach(function (html) {
                var switchery = new Switchery(html, { size: 'small' });
                if (switchery.element.id == "chkserver-learn") {
                    $scope.white_list_switchery = switchery;
                }
            });
            elems.forEach(function (html) {
                html.onchange = function () {
                    var k = $(html).attr("data-k")
                    var cindex = $(html).attr("data-cindex")
                    $scope.UpdateBool(cindex, k);
                };
            });
        });

        add_tab_click_event()
    });

    $scope.upload_system = "";
    $scope.agent_type = { "system": "null", "platform": "null" };
    $scope.agent_type_lst = [
        { "system": "windows", "platform": 32 },
        { "system": "windows", "platform": 64 },
        { "system": "linux", "platform": 64 }
    ]
    $scope.filename = "";
    $scope.select_file_type = function (type_) {
        $scope.agent_type = type_;
        create_file_form();
    }
    $scope.pull_all = function (system) {
        json = {
            "name": "更新所有 " + system + "下的agent",
            "command": "",
            "host_list": [system],
            "type": "update"
        };
        request_password(function (password) {
            $.ajax({
                url: enable_pass(task_url, password),
                method: "post",
                contentType: "application/json; charset=utf-8",
                data: JSON.stringify(json),
                success: function (data) {
                    if (data.status) {
                        setTimeout(() => {
                            swal('Updated!', '已为所有' + system + '机器添加更新任务.', 'success');
                        }, 100);
                    }
                    if (data.status == false) {
                        ajaxcallback(data);
                    }
                },
                error: function (data) {
                    console.log(data)
                }
            })
        });
    }

    window.ajax_success = function () {
        response = JSON.parse(this.responseText);
        if (response.status) {
            $scope.upload_system = $scope.agent_type.system;
            $scope.fileinfo = "文件已上传";
            $scope.$apply();
            setTimeout(function(){
                swal('文件已上传', '你已经上传了新的agent，可点击快速更新，为所有agent添加更新任务。', 'success');
            }, 100);
        } else {
            $scope.fileinfo = "文件上传失败";
            $scope.$apply();
            ajaxcallback(response);
        }
    }

    window.ajax_submit = function (formElement) {
        url = formElement.action;
        request_password(function (password){
            url = url.url_update_query('pass', password);
            $scope.fileinfo = "文件上传中...";
            $scope.$apply();
            (function () {
                if (!url) { return; }
                var xhr = new XMLHttpRequest();
                xhr.onload = ajax_success;
                xhr.open("post", url);
                xhr.setRequestHeader('requesttoken', request_token);
                xhr.send(new FormData(formElement));
            })();
        });
    }

    $scope.fileinfo = "请选择文件";
    change_fileinfo = function() {
        if ($('#file-upload input.file')[0].value) {
            $scope.fileinfo = $('#file-upload input.file')[0].value;
        } else {
            $scope.fileinfo = "文件选择错误，请重新选择"
        }
        $scope.$apply();
    }

    create_file_form = function () {
        $scope.upload_url = file_url
            .url_update_query("system", $scope.agent_type.system)
            .url_update_query("platform", $scope.agent_type.platform);
        $('span.upload-btn').show();
    }
});

hostw.controller('hostlist', function ($scope, $http, Notification) {

    modal_option_click();
    $scope.before_1minutes = before_1minutes;
    $scope.str2Date = str2Date;
    $scope.current_page = 1;
    $scope.hosts = [];
    $scope.monitor_list = [];
    $scope.monitor_ip = "";
    $scope.current_monitor_url = "";
    $scope.timer = 0;
    $scope.has_show_ids = [];
    $scope.search_filter = '';
    $scope.command_show = false;
    $scope.health_data = {
        0: {'style': 'badge-success', 'word': 'ONLINE'},
        1: {'style': '', 'word': 'OFFLINE'},
        2: {'style': 'badge-warning', 'word': '无法连接任务推送端口'}
    }

    http_domain = document.domain + ":" + httpport;
    $scope.command_list = {
        "linux-64": `wget -O /tmp/daemon http://${http_domain}/json/download?type=daemon\\&system=linux\\&platform=64\\&action=download;chmod +x /tmp/daemon;/tmp/daemon -install -netloc ${location.host}`,
        "windows-64": `cd %SystemDrive% & certutil -urlcache -split -f http://${http_domain}/json/download?type=daemon^&system=windows^&platform=64^&action=download daemon.exe & daemon.exe -install -netloc ${location.host}`,
        "windows-32": `cd %SystemDrive% & certutil -urlcache -split -f http://${http_domain}/json/download?type=daemon^&system=windows^&platform=32^&action=download daemon.exe & daemon.exe -install -netloc ${location.host}`,
        "windows-32-powershell": `[System.Net.ServicePointManager]::ServerCertificateValidationCallback={$true};(New-Object System.Net.WebClient).DownloadFile("https://${location.host}/json/download?type=daemon&system=windows&platform=32&action=download", "C:\\daemon.exe");C:\\daemon.exe -install -netloc ${location.host};`,
        "windows-64-powershell": `[System.Net.ServicePointManager]::ServerCertificateValidationCallback={$true};(New-Object System.Net.WebClient).DownloadFile("https://${location.host}/json/download?type=daemon&system=windows&platform=64&action=download", "C:\\daemon.exe");C:\\daemon.exe -install -netloc ${location.host};`
    }

    $scope.toggle_command = function () {
        $scope.command_show = !$scope.command_show;
    }
    // hosts 为空时，默认显示安装 agent 命令
    if ($scope.hosts.length == 0) {
        $scope.toggle_command();
    }

    default_history = ["windows", "linux", "web", "db", "offline", "online", "can-not-push"]
    $scope.history = default_history.concat(localStorage.host_filter_history.split(":::"));
    $scope.re_search = function (history) {
        $scope.search_filter = history;
        $scope.hosts = [];
        $scope.get_result();
    }

    $scope.SetHosttext = function (ip) {
        $('div.modal-content div.col-md-12 textarea').text(ip);
    }

    $scope.add_filter = function () {
        swal({
            title: "主机列表过滤器",
            text: "请输入需要搜索的ip, 主机名等信息",
            type: "input",
            showCancelButton: true,
            inputPlaceholder: "filter"
        },
            function (inputValue) {
                if (inputValue) {
                    add_history("host_filter_history", inputValue)
                    $scope.search_filter = inputValue;
                    $scope.hosts = [];
                    $scope.get_result();
                } else {
                    $scope.search_filter = "";
                }
            }
        );
    }

    $scope.get_result = function () {
        $http.get(
            client_url.url_add_Paginator($scope.current_page).url_update_query('q', $scope.search_filter).url_update_query('limit', 24)
        ).then(
            function (response) {
                if (response.data) {
                    $scope.hosts.push.apply($scope.hosts, response.data);
                    Notification.success("成功加载 " + response.data.length + " 条数据。");
                } else {
                    Notification.error("没有其它数据了。");
                }

            });
    }

    $scope.get_more = function () {
        $scope.current_page = $scope.current_page + 1;
        $scope.get_result()
    }

    $scope.get_result();

    $('#monitor-modal').on('hidden.bs.modal', function () {
        console.log('kill the interval timer : ' + $scope.timer);
        clearInterval($scope.timer);
    })

    monitor_scroll = function () {
        url = $scope.current_monitor_url.url_update_query("t", 7);
        $.ajax({
            url: url,
            type: "GET",
            success: function (data) {
                monitor_append(data);
            }
        });
    }

    monitor_append = function (lst) {
        side = $('#monitor_info_list');
        for (key in lst) {
            data = lst[key];
            if ($scope.has_show_ids.indexOf(data._id) < 0) {
                side.prepend($('<p>').html(format_msg(data)));
                $scope.has_show_ids.push(data._id);
            }
        }
    }

    $scope.show_monitor = function (ip) {
        $("#side-modal-title").text(ip);
        $scope.has_show_ids = [];
        clearInterval($scope.timer);
        $scope.current_monitor_url = client_url.url_update_query("ip", ip)
        $('#monitor_info_list').empty();
        $.ajax({
            url: $scope.current_monitor_url,
            type: "GET",
            success: function (data) {
                monitor_append(data.reverse());
            }
        });
        $scope.timer = setInterval("monitor_scroll()", 3 * 1000);
    }

});

hostw.controller('detailinfo', function ($scope, $http, $rootScope, $routeParams) {
    $scope.start = 0;
    $scope.monitorlist = [];
    $scope.type = "";
    $http.get(info_url + '/' + $routeParams.ip + '/').then(function (response) {
        $scope.hostip = $routeParams.ip
        $scope.infolist = response.data
        $scope.type = $scope.infolist.infodata[0].type
        add_tab_click_event()
    });
    $scope.appendNewTwenty = function (type) {
        $http.get(monitor_url + '/' + $routeParams.ip + '/' + type + '/' + $scope.start).then(function (response) {
            newlist = response.data;
            Array.prototype.push.apply($scope.monitorlist, newlist);
            if (newlist) {
                $scope.start = $scope.start + newlist.length;
            }
        });
    }
});

hostw.controller('notice', function ($scope, $http, $window, Notification, globalVariable) {
    $scope.noticelist = [];
    $scope.langtem = HostWData[lang];
    $scope.style = HostWData.style;
    $scope.currenttab = 0;
    $scope.currenttype = 'undeal';
    $scope.currentpage = 1;
    $scope.iswatchmode = iswatchmode;
    $scope.is_close_messagebox = false;
    $scope.close_msgbox = function(){$scope.is_close_messagebox = true};

    $scope.notice_type_words = {
        'undeal':"未处理",
        'dealed':"已处理",
        'ignore':"已忽略",
    }

    $scope.search_analyze = function (notice) {
        keylist = ["info", 'ip', "type"]
        query_list = [];
        keylist.forEach(function (key) {
            value = notice[key];
            if (key == "info") {
                key_ = "all";
                if (value.indexOf('|') > 0) {
                    value = value.split('|').sort(function (a, b) { return b.length - a.length; })[0];
                }
            } else {
                key_ = key;
            }
            query_list.push(key_ + ":" + value);
        });
        globalVariable.notice_search = query_list.join(" | ");
        $window.location.href = "/#!/analyze";
    }

    $scope.add2config = function (type, info, config_type) {
        if (info.indexOf('|') > 0) {
            swal("0ops!!!!", "该信息无法自动化添加到黑白名单，请手动添加~", "error");
            return
        }
        keydict = {
            "process": "process",
            "loginlog": "ip",
            "connection": "ip"
        }

        if (type == 'loginlog') {
            ara = /[\d\.]+/.exec(info);
            if (ara) {
                info_ = ara[0];
            }
        } else {
            info_ = info;
        }

        if (type in keydict) {
            key = keydict[type];
        } else {
            key = 'other';
        }

        if ((key == "process") || (key == "other")) {
            info_ = string2regexp(info_)
        }

        json = {
            id: config_type,
            key: key,
            input: info_
        }
        request_password(function (password) {
            config_url = config_url.url_update_query('pass', password);
            req = {
                method: 'PUT',
                url: config_url,
                data: json
            }
            $http(req).then(function (response) {
                if (response.data.status) {
                    swal('添加成功!', '该信息已经添加到设置中，部分设置可能需要到config页面修改', 'success');
                    return
                }
                if (response.data.status == false) {
                    ajaxcallback(response.data);
                }
            })
        })
    }

    $scope.kill_all = function (process) {
        if (process.indexOf('|') > -1) {
            swal('无法添加', '该信息无法直接添加到阻断任务里，请手动添加kill任务。', 'error');
            return;
        }
        json = {
            "name": "阻断操作 - 对象：" + process,
            "command": process,
            "host_list": ["all"],
            "type": "kill"
        };
        swal({
            title: '确定阻断进程：' + process,
            text: "是否要在所有服务器中删除进程：" + process,
            type: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'KILL IT!'
        },
            function (isconfirm) {
                if (isconfirm) {
                    request_password(function (password) {
                        $.ajax({
                            url: enable_pass(task_url, password),
                            method: "post",
                            contentType: "application/json; charset=utf-8",
                            data: JSON.stringify(json),
                            success: function (data) {
                                if (data.status) {
                                    swal('Killed!', '添加阻断任务成功.', 'success');
                                }
                                if (response.data.status == false) {
                                    ajaxcallback(response.data);
                                }
                            },
                            error: function (data) {
                                console.log(data)
                            }
                        })
                    });
                };
            });
    }

    $scope.add_filter = function () {
        swal(
            {
                title: "告警信息过滤器",
                text: "请输入需要搜索的ip, 告警类型等信息",
                type: "input",
                showCancelButton: true,
                closeOnConfirm: true,
                inputPlaceholder: "Filter"
            },
            function (inputValue) {
                if (inputValue) {
                    $scope.filter = inputValue;
                    $scope.noticelist = [];
                    $scope.get_result($scope.currenttype, 1);
                } else {
                    $scope.filter = "";
                }
            }
        );
    }

    $scope.get_result = function (t, p) {
        url = notice_url.url_update_query("status", t).url_add_Paginator(p);
        if ($scope.filter != "" && $scope.filter != undefined) {
            url = url.url_update_query("q", $scope.filter);
        }
        $http.get(url).then(function (response) {
            if (response.data) {
                $scope.noticelist.push.apply($scope.noticelist, response.data);
                Notification.success("成功加载 " + response.data.length + " 条数据。");
            } else {
                Notification.error("没有其它数据了。");
            }
        });
    }

    $scope.change_type = function (ntype) {
        $scope.currenttype = ntype;
        $scope.noticelist = [];
        $scope.get_result($scope.currenttype, 1);
    }

    $scope.get_result($scope.currenttype, 1);

    $scope.get_more = function () {
        $scope.currentpage = $scope.currentpage + 1;
        $scope.get_result($scope.currenttype, $scope.currentpage);
    }

    $scope.onClickTab = function (tab) {
        $scope.currenttab = tab
    }

    $scope.change_status = function (id, status) {
        notice = $scope.noticelist[id]
        $.ajax({
            url: notice_url,
            method: "post",
            dataType: "json",
            data: JSON.stringify({
                id: notice._id,
                status: status
            }),
            success: function (data) {
                console.log(data);
                if (data.status) {
                    notice.status = -1;
                    $scope.currenttab = id + 1;
                    $scope.$apply();
                }
            }
        })
    }

    $scope.change_all = function (status) {
        request_password(function (password) {
            notice_url = notice_url.url_update_query('pass', password);
            $.ajax({
                url: notice_url,
                method: "post",
                dataType: "json",
                data: JSON.stringify({
                    id: 'all',
                    status: status
                }),
                success: function (data) {
                    if (data.status) {
                        Notification.success("操作成功。");
                        location.reload();
                    } else {
                        setTimeout(() => {
                            swal('操作失败', data.msg, 'error');
                        }, 100);
                    }
                }
            });
        });
    }

});

hostw.controller('statistics', function ($scope, $http, Notification) {
    Notification.info("正在加载数据，请耐心等待...");
    $scope.langtem = HostWData[lang]
    $scope.style = HostWData.style
    $http.get(statistics_url + "?type=topmsg").then(
        function (response) {
            $scope.data = response.data;
            console.log(response.data.listdata);
        });
    $http.get(statistics_url + "?type=total").then(
        function (response) {
            $scope.totaldata = response.data;
        });

    $(function () {
        var myChart = echarts.init(document.getElementById('chartOne'));
        $.get(statistics_url + "?type=pie", function (data) {
            option = {
                tooltip: {
                    trigger: 'item',
                    formatter: "{a} <br/>{b}: {c} ({d}%)"
                },
                legend: {
                    orient: 'vertical',
                    x: 'left',
                    data: data.xdata
                },
                series: [
                    {
                        name: '风险类型',
                        type: 'pie',
                        selectedMode: 'single',
                        radius: [0, '30%'],

                        label: {
                            normal: {
                                position: 'inner'
                            }
                        },
                        labelLine: {
                            normal: {
                                show: false
                            }
                        },
                        data: data.listdataInner
                    },
                    {
                        name: '策略名称',
                        type: 'pie',
                        radius: ['40%', '55%'],

                        data: data.listdata
                    }
                ]
            };
            myChart.setOption(option);
        });
    });

    $(function () {
        var myChartTwo = echarts.init(document.getElementById('chartTwo'));
        $.get(statistics_url + "?type=time", function (data) {

            var hours = ['0', '1', '2', '3', '4', '5', '6',
                '7', '8', '9', '10', '11',
                '12', '13', '14', '15', '16', '17',
                '18', '19', '20', '21', '22', '23'];
            var days = data.xdata;

            count_lst = [Array.from(Array(24), () => 0), Array.from(Array(24), () => 0)];
            if (data.listdata == null) {
                data.listdata = [[0, 0, 0]];
            }

            data.listdata.forEach(function (item) {
                count_lst[item[0]][item[1]] = item[2];
            })

            option = {
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'shadow'
                    }
                },
                legend: {
                    data: ['昨天', '今天']
                },
                grid: {
                    left: '3%',
                    right: '4%',
                    bottom: '3%',
                    containLabel: true
                },
                xAxis: [
                    {
                        type: 'category',
                        data: hours
                    }
                ],
                yAxis: [
                    {
                        type: 'value'
                    }
                ],
                series: [
                    {
                        name: '昨天',
                        type: 'bar',
                        data: count_lst[0]
                    },
                    {
                        name: '今天',
                        type: 'bar',
                        data: count_lst[1]
                    }
                ]
            };

            myChartTwo.setOption(option);
        });
    });

    $(function () {
        var myChart = echarts.init(document.getElementById('chartLine'));
        $.get(statistics_url + "?type=line", function (data) {
            option = {
                title: {
                    text: ''
                },
                tooltip: {
                    trigger: 'axis'
                },
                legend: {
                    data: ['危险', '可疑', '风险']
                },
                grid: {
                    left: '3%',
                    right: '4%',
                    bottom: '3%',
                    containLabel: true
                },
                toolbox: {
                    feature: {
                        saveAsImage: {}
                    }
                },
                xAxis: {
                    type: 'category',
                    boundaryGap: false,
                    data: data.xdata
                },
                yAxis: {
                    type: 'value'
                },
                series: data.listdata
            };
            myChart.setOption(option);
        });
    });
});


// main.js
warning = function (title) {
    swal({
        title: title,
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#DD6B55",
        closeOnConfirm: true
    });
}

ajaxcallback = function (data) {
    setTimeout(() => {
        if (data.status) {
            swal("成功", "该动作已经作用到服务器", "success")
        } else if (data.msg) {
            swal("Oops...", data.msg, "error")
        } else {
            swal("Oops...", "something error", "error")
        }
    }, 200)
}

newtask = function () {
    json = {}
    submit = true
    $("#newTaskModel input:text, textarea").each(function () {
        input = $(this)
        if (input.val() == "" && input.attr("name") != "command") {
            warning(input.attr("name") + "不得为空！")
            submit = false
        }
        if (input.attr("name") == "host_list") {
            json[input.attr("name")] = input.val().split(',')
        } else {
            json[input.attr("name")] = input.val()
        }
    })

    $("#newTaskModel select").each(function () {
        input = $(this)
        json[input.attr("name")] = input.val()
    });
    if (submit) {
        $('#model-close').click();
        request_password(function (password) {
            $.ajax({
                url: enable_pass(task_url, password),
                method: "post",
                contentType: "application/json; charset=utf-8",
                data: JSON.stringify(json),
                success: function (data) {
                    if (data.status) {
                        location.reload()
                    }
                    ajaxcallback(data)
                },
                error: function (data) {
                    console.log(data)
                }
            })
        })
    }
    return json
}

add_tab_click_event = function () {
    $(document).ready(function () {
        $("div.card-header ul li").click(function () {
            currenttab = $.trim(this.textContent)
            $("div.tab-pane").removeClass("active")
            $("div.tab-pane#" + currenttab).addClass("active")
        })
    });
}

get_input = function (title, text, origin_value, callback) {
    // if (title.indexOf('_') == (title.length - 1)) {
    //     $("#span_" + title.toLowerCase()).find("pre").attr("contenteditable", "plaintext-only").css("color", "#20a3b9");
    //     $("#save_" + title.toLowerCase()).css("display", "inline-block");
    //     $("#edit_" + title.toLowerCase()).css("display", "none");
    //     return;
    // }
    swal({
        title: "You are editing \"" + title + "\"",
        text: text,
        type: "input",
        showCancelButton: true,
        closeOnConfirm: false
    },
        function (inputValue) {
            callback(inputValue)
        }
    );
    jQuery('.sweet-alert input[type=text]:first' ).val(origin_value);
}

edit_config_item = function (json, currenttab) {
    request_password(function (password) {
        config_url = config_url.url_update_query('pass', password);
        $.ajax({
            url: config_url,
            method: "post",
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify(json),
            success: function (data) {
                if (data.status) {
                    $('#' + json.key + ' tbody td span.cfgvalue').text(json.input)
                }
                ajaxcallback(data);
            },
            error: function (data) {
                console.log(data);
            }
        })
    })
}


