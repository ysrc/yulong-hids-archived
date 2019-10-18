<!-- hostlist.tpl -->
<script type="text/ng-template" id="hostlist">
  <div ng-controller="hostlist" class="hostlist">
    <div class="col-md-12">
    <div class="col-md-12">
        <span class="badge badge-icon badge-text history" ng-repeat="h in history track by $index" ng-if="h != ''" ng-click="re_search(h)">
          <i class="fa fa-history" aria-hidden="true"></i>{{ h }}
        </span>
    </div>
    <div class="command col-md-12" ng-show="command_show">
      <h3>安装命令：</h3>
      <p class="red"><i class="fa fa-warning" aria-hidden="true"></i>注意(WARNING)：驭龙是一个安全HIDS系统，其Agent在安装和使用的过程中会hook系统内核，<b>请勿未做调研和实验就直接在生产环境或重要主机上安装Agent。</b></p>
      <li ng-repeat="(cmd_index, cmd) in command_list">
          <span class="highlight">{{ cmd_index }}</span> : {{ cmd }}
      </li>
    </div>
    <div class="col-md-6 col-lg-6 col-xl-4 hostitem" ng-repeat="host in hosts | filter: search track by $index">
      <div class="card card-banner card-green-light">
        <div class="card-body">
           <div ng-if="host.system.indexOf('indows')>-1">
            <i class="icon fa fa-windows fa-4x"></i>
            </div>
          <div ng-if="host.system.indexOf('indows')<=-1">
            <i class="icon fa fa-linux fa-4x"></i>
            </div>
          <div class="content">
            <span class="badge badge-primary badge-icon"><i class="fa fa-tasks" aria-hidden="true"></i><span>{{ host.system }}</span></span>
            <p>
              <span class="badge badge-primary badge-icon">
                <i class="fa fa-clock-o" aria-hidden="true"></i>
                <span>{{ timeformat(host.uptime) }}</span>
              </span>
              <span class="badge badge-info badge-icon" ng-if="host.type">
                <i class="fa fa-tag" aria-hidden="true"></i>
                <span>{{ host.type }}</span>
              </span>
              <span class="badge badge-icon" ng-class="health_data[host.health]['style']" ng-if="host.health != undefined">
                <i class="fa fa-leaf" aria-hidden="true"></i>
                <span>{{ health_data[host.health]['word'] }}</span>
              </span>
            </p>
            <div class="title">{{ host.hostname }}</div>
            <div class="value">{{ host.ip }}</div>
            <div class="btn-group">
              <a ng-href="{{ '/#!/info/'+ host.ip }}" class="btn btn-success">信息</a>
              <button data-toggle="modal" ng-click="show_monitor(host.ip)" data-target="#monitor-modal" class="btn btn-success">监控</button>
              <button class="btn  btn-success" data-host="{{ host.ip }}" data-toggle="modal" ng-click="SetHosttext(host.ip)" data-target="#newTaskModel">推送</button>
            </div>
          </div>
        </div>
        </a>
      </div>
    </div>

    <button class="icon append-icon" id="addbutton" ng-click="toggle_command()">
        <i class="fa fa-plus" aria-hidden="true" title="显示或隐藏安装命令"></i>
    </button>

    <button class="icon append-icon" id="searchbutton" ng-click="add_filter()">
        <i class="fa fa-search-minus" aria-hidden="true" title="添加过滤器"></i>
    </button>

    <button class="icon append-icon" id="fixedbutton" ng-click="get_more()">
        <i class="fa fa-caret-down" aria-hidden="true" title="加载更多"></i>
    </button>

  </div>

<!-- Modal -->
<div class="modal right fade" id="monitor-modal" tabindex="-1" role="dialog" aria-labelledby="side-modal-title">
    <div class="modal-dialog" role="document">
        <div class="modal-content modal-black">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title" id="side-modal-title"></h4>
            </div>
            <div class="modal-body" id="monitor_info_list">
                
            </div>

        </div><!-- modal-content -->
    </div><!-- modal-dialog -->
</div><!-- modal -->

  <div class="modal fade" id="newTaskModel" tabindex="-1" role="dialog" aria-labelledby="newTaskModelLabel" aria-hidden="true">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal" aria-hidden="true">
                    &times;
                </button>
          <h4 class="modal-title" id="newTaskModelLabel">
            新建任务
          </h4>
        </div>
        <div class="modal-body">
          <div class="row">
            <div class="col-md-6">
              <input type="text" class="form-control" placeholder="name" check-type="required" name="name">
            </div>
            <div class="col-md-6">
              <select name="type" class="form-control">
                <option value="">-请选择-</option>
                <option value="kill">kill</option>
                <option value="uninstall">uninstall</option>
                <option value="update">update</option>
                <option value="delete">delete</option>
                <option value="exec">exec</option>
                <option value="reload">reload</option>
                <option value="quit">quit</option>
              </select>
            </div>
            <div class="col-md-12">
              <input type="text" class="form-control" placeholder="command" check-type="required" name="command">
            </div>
            <div class="col-md-12">
              <textarea check-type="required" id="host_list" name="host_list" rows="3" class="form-control" placeholder="host_list;每个host以逗号隔开..."></textarea>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal" id="model-close">关闭
          </button>
          <button type="button" class="btn btn-info" id="newtask" onclick="newtask()">
            提交
          </button>
        </div>
      </div>
    </div>
  </div>
</div>
</div>
</script>