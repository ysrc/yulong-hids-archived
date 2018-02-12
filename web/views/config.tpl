<!-- config.tpl -->

<script type="text/ng-template" id="config">
<div ng-controller="config" id="config-div" class="config">

<div class="col-md-12 top-nav">
    <div class="col-md-6 left-title">
        设置面板
    </div>
</div>

<div class="col-lg-12 col-md-12 col-sm-12 col-xs-12">
  <div class="card card-tab card-mini">
    <div class="card-header">
      <ul class="nav nav-tabs tab-stats">
        <li role="{{ config.type }}" id="{{ config.type }}_tab" ng-repeat="config in configs| filter: search track by $index" ng-class="{active:type == config.type}"  data-tooltip="Refresh">
          <a aria-controls="{{ config.type }}" role="tab" data-toggle="tab" aria-expanded="false">{{ config.type }}</a>
        </li>

        <li role="update" id="update_tab" class="ng-scope">
          <a aria-controls="update" role="tab" data-toggle="tab" aria-expanded="false">update</a>
        </li>
      </ul>
    </div>
    <div class="card-body tab-content">
      <div role="tabpanel" class="tab-pane" id="{{ config.type }}" ng-repeat="(cindex, config) in configs" ng-class="{active:type == config.type}">
        <div class="card-body no-padding table-responsive row">
          <div class="table card-table col-md-4" ng-repeat="(k,items) in config.dic" id="{{ k }}">
            <div class="config-key col-md-2 right black">
              {{ config_first_child(config.type, k) }}
              <i class="fa fa-question-circle-o" aria-hidden="true"></i>
            </div>
            <div ng-if="isArray(items)" class="config-item col-md-10">
              <span class="badge badge-success badge-icon" ng-repeat="item in items" ng-mouseover="showdel($event)" ng-click="del($event, config._id, cindex, k)"
                ng-mouseout="hidedel($event)">
                <i class="fa fa-check" aria-hidden="true"></i>
                <span>{{ item }}</span>
              </span>
              <span class="badge badge-info badge-icon additem" ng-click="add(cindex, k)">
                <i class="fa fa-plus" aria-hidden="true"></i>
                <span>New in {{ k | uppercase }}</span>
              </span>
            </div>
            <div ng-if="!isArray(items)" class="config-item col-md-10">
              <div ng-if="typeof(items) == 'boolean'">
                <input type="checkbox" class="js-switch" data-k="{{ k }}" data-cindex="{{cindex}}" ng-change="UpdateBool(cindex, k)" ng-model="chk1"
                  id="chk{{configs[cindex].type}}-{{ k }}" ng-checked="{{ items }}" />
              </div>
              <div ng-if="typeof(items) == 'string' || typeof(items) == 'number'">
                <div class="cfgvalue" id="span_{{ k }}" >
                  <pre class="config-text">{{ items }}</pre>
                </div>
                <span class="badge badge-icon edit" id="edit_{{ k }}" ng-click="edit(cindex, k)" ng-show="disable_edit_list.indexOf(k) < 0">
                  <i class="fa fa-pencil" aria-hidden="true"></i>
                  <span>Edit</span>
                </span>
                <span class="badge badge-icon save" id="save_{{ k }}" ng-click="saveInfo(cindex, k)">
                  <i class="fa fa-pencil" aria-hidden="true"></i>
                  <span>Save</span>
                </span>
              </div>
            </div>

          </div>
        </div>
      </div>

      <div class="modal fade" id="model-select-learn">
          <div class='modal-dialog'>
              <div class='modal-content'>
                <div class='modal-header'>
                    <button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button>
                    <h4 class='modal-title'>
                        <strong>单击高频(>1)信息以添加到白名单,单击低频(=1)信息以忽略该信息</strong>
                    </h4>
                </div>
                <div class='modal-body'>
                    <div class="div-container">
                        <span class="badge badge-warning badge-icon" ng-repeat="notice in learn_data" ng-mouseover="showadd($event)" ng-mouseout="hideadd($event)" ng-click="add_white($event, notice._id.type, notice._id.info, notice.count)" ng-if="!(notice._id.type == 'process' && notice._id.info.includes('|'))">
                            <i class="fa {{ learn_data_fa(notice._id.type) }}" aria-hidden="true"> {{ langtem.notice.data.type[notice._id.type] }} </i>
                            <i class="count" aria-hidden="true"> {{ notice.count }} </i>
                            <span>{{ notice._id.info }}</span>
                        </span>
                  </div>
                </div>
                <div class='modal-footer'>
                    <button class="btn btn-danger btn-square pull-right" ng-click="close_learn()">
                        <i class="fa fa-cloud" aria-hidden="true"></i>
                        关闭观察模式
                    </button>
                    <button class="btn btn-warning btn-square pull-right" ng-click="clean_notice()">
                        <i class="fa fa-fast-forward" aria-hidden="true"></i>
                        移动剩余告警到未处理
                    </button>
                </div>
              </div>
          </div>
      </div>

      <div role="tabpanel" class="tab-pane" id="update" ng-class="{active:type == config.type}">
        <div class="card-body no-padding table-responsive filelist" style="height:400px;">
          <h5>
            当前选择的类型为 <span class="highlight">{{ agent_type.platform }}</span> 位 <span class="highlight">{{ agent_type.system }}</span> 系统, 可上传更新该类型agent。
          </h5>

          <div class="btn-group">
            <button type="button" class="btn btn-default dropdown-toggle" data-toggle="dropdown">
                <i class="fa fa-upload" aria-hidden="true"></i>选择agent类型并上传agent文件  <span class="caret"></span>
            </button>
            <ul class="dropdown-menu" role="menu">
              <li ng-repeat="agent_type in agent_type_lst" ng-click="select_file_type(agent_type)"><a>{{ agent_type.system }} - {{ agent_type.platform }}</a></li>
            </ul>

            <span class="upload-btn">
              <button class="btn btn-info btn-square" onclick="$('form input.file').click()"> {{ fileinfo }} </button>
              <button class="btn btn-success btn-square" onclick="$('form input.submit').click()">上传文件</button>
            </span>

            <form action="{{ upload_url }}" method="post" enctype="multipart/form-data" onsubmit="ajax_submit(this); return false;" id="file-upload">
              <input type="file" name="file" class="file" onchange="change_fileinfo()"></label>
              <input type="submit" value="上传文件" class="submit"/>
            </form>

          </div>

          <button class="btn btn-warning col-xs-1" ng-click="pull_all(upload_system)" ng-show="upload_system">快速推送</button>

        </div>
      </div>
    </div>
  </div>
</div>

<style>
  /* 一个历史原因的无奈之举 */
  div.config div.table.card-table.col-md-4 {
    padding-bottom: 10px;
    border-bottom: 1px rgba(141,146,147,0.3) dotted;
  }

  div.config cfgvalue {
    float: left;
  }

  div.config badge badge-icon edit {
    margin-left: 5px;
  }

  div.config modal fade{
    display: none;
  }
</style>

</script>