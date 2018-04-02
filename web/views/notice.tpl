<!-- notice.tpl -->
<script type="text/ng-template" id="notice">
  <div ng-controller="notice" class="notice">

      <div class="alert alert-warning alert-normal-warning" ng-if="iswatchmode == 'true' && !is_close_messagebox">
          <button type="button" class="close" ng-click="close_msgbox()">×</button>
            观察模式已开启，所有新告警将不会被显示到告警列表。关于观察模式的更多详细信息，请查看文档。
      </div>

      <div class="col-md-12 top-nav">
          <div class="col-md-6 left-title">
              告警列表
          </div>
          <div class="col-md-6">
            <button class="btn btn-danger btn-square pull-right" ng-click="change_all(2)">
                  <i class="fa fa-exclamation-circle" aria-hidden="true"></i>
                  忽略所有未处理
              </button>
            <button class="btn btn-warning btn-square pull-right" ng-click="change_all(1)">
                  <i class="fa fa-exclamation-triangle" aria-hidden="true"></i>
                  把所有未处理警告设置为已处理
              </button>
          </div>
      </div>
      <div class="app-messaging-container">
        <div class="app-messaging collapse in" id="collapseMessaging" aria-expanded="true">
          <div class="chat-group col-lg-4 col-md-4 col-sm-12 col-xs-12">
            <div class="heading item-list">
              告警列表 |
              <span>
                {{ filter }}
                <a><i ng-click="add_filter()" class="fa fa-pencil-square-o" aria-hidden="true"></i></a>
              </span>
              <div class="notice-type pull-right">
                <a ng-click="change_type('undeal')"> {{ notice_type_words['undeal'] }}</a> |
                <a ng-click="change_type('dealed')"> {{ notice_type_words['dealed'] }}</a> |
                <a ng-click="change_type('ignore')"> {{ notice_type_words['ignore'] }}</a>
              </div>
            </div>
            <ul class="group full-height">
              <li class="section">{{ notice_type_words[currenttype] }}</li>
              <li class="message" ng-repeat="notice in noticelist | filter: search track by $index" title="{{ notice.info }}" ng-class='{active:$index == currenttab}' id="item{{ $index }} " ng-click="onClickTab($index)" ng-if="notice.status > -1">
                <a data-toggle="collapse" href="" aria-expanded="true" aria-controls="collapseMessaging">
                  <span class="{{ style.notice.level[notice.level] }} pull-right">{{ langtem.notice.data.level[notice.level] }}</span>
                  <span class="badge badge-success pull-right">{{ notice.ip }}</span>
                  <span class="badge badge-success pull-right">{{ langtem.notice.data.type[ notice.type ] }}</span>
                  <div class="message">
                    <div class="content">
                      <div class="title">
                        告警信息: <b>{{ notice.info | cutWords:30 }}</b>
                      </div>
                    </div>
                  </div>
                </a>
              </li>
              <li class="more" ng-click="get_more()">点击加载更多告警</li>
            </ul>
          </div>
          <div class="messaging col-lg-12 col-md-12 col-sm-12 col-xs-12" ng-repeat="notice in noticelist" ng-if="currenttab == $index" ng-hide="currenttab != $index">
            <div class="heading">
              <div class="title">
                <b>告警详情:</b>  ({{ notice.info | cutWords:50 }})
                <span class="{{ style.notice.status[notice.status] }} badge-icon">
                  <i class="fa fa-circle" aria-hidden="true"></i>
                  <span>{{ langtem.notice.data.status[ notice.status ] }}</span>
                </span>
              </div>
              <div class="action">
              </div>
            </div>
            <div class="col-md-12 messagebody">
              <table class="table table-striped table-bordered table-hover">
                <tbody>
                  <tr>
                    <td class="key">{{ langtem.notice.key.info }}</td>
                    <td class="v">{{ notice["info"] }}</td>
                  </tr>
                  <tr>
                    <td class="key">{{ langtem.notice.key.description }}</td>
                    <td class="v">{{ notice["description"] }}</td>
                  </tr>
                  <tr>
                    <td class="key">{{ langtem.notice.key.time }}</td>
                    <td class="v">{{ timeformat(notice["time"]) }}</td>
                  </tr>
                  <tr>
                    <td class="key">{{ langtem.notice.key.type }}</td>
                    <td class="v">{{ langtem.notice.data.type[notice["type"]] }}</td>
                  </tr>
                  <tr>
                    <td class="key">{{ langtem.notice.key.ip }}</td>
                    <td class="v">{{ notice["ip"] }}</td>
                  </tr>
                  <tr>
                    <td class="key">{{ langtem.notice.key.source }}</td>
                    <td class="v">{{ notice["source"] }}</td>
                  </tr>
                  <tr>
                    <td class="key">{{ langtem.notice.key.level }}</td>
                    <td class="v">
                        <span class="badge {{ style.notice.level[notice['level']] }} badge-icon">
                            <i class="fa fa-exclamation-triangle" aria-hidden="true"></i>
                            {{ langtem.notice.data.level[notice['level']] }}
                        </span>
                      </td>
                  </tr>
                  <tr>
                    <td class="key">{{ langtem.notice.key.raw }}</td>
                    <td class="v">
                      {{ notice["raw"] }}
                    </td>
                  </tr>

                </tbody>
              </table>
            </div>
            <div class="col-md-12 done-notice">
              <div class="dropdown pull-right">
                <button class="btn btn-default dropdown-toggle" type="button" id="menu1" data-toggle="dropdown">
                  状态处理
                  <span class="caret"></span>
                </button>
                <ul class="dropdown-menu" role="menu" aria-labelledby="menu1">
                  <li role="presentation">
                    <a role="menuitem" tabindex="-1" ng-click="change_status($index, 1)">已处理</a>
                  </li>
                  <li role="presentation">
                    <a role="menuitem" tabindex="-1" ng-click="change_status($index, 2)">忽略</a>
                  </li>
                </ul>
              </div>
              <div class="dropdown pull-right">
                <button class="btn btn-default dropdown-toggle" type="button" id="menu1" data-toggle="dropdown">
                  黑白名单
                  <span class="caret"></span>
                </button>
                <ul class="dropdown-menu" role="menu" aria-labelledby="menu1">
                  <li role="presentation">
                    <a role="menuitem" tabindex="-1" ng-click="add2config(notice.type, notice.info, 'blacklist')">加入黑名单</a>
                  </li>
                  <li role="presentation">
                    <a role="menuitem" tabindex="-1" ng-click="add2config(notice.type, notice.info, 'whitelist')">加入白名单</a>
                  </li>
                </ul>
              </div>
              <div class="dropdown pull-right">
                <button class="btn btn-default dropdown-toggle" type="button" id="menu1" data-toggle="dropdown" ng-class=" notice['type'] == 'process' ? '' : 'disabled' ">
                  阻断操作
                  <span class="caret"></span>
                </button>
                <ul class="dropdown-menu" role="menu" aria-labelledby="menu1">
                  <li role="presentation">
                    <a role="menuitem" tabindex="-1" ng-click="kill_all(notice['info'])">阻断</a>
                  </li>
                </ul>
              </div>
              <button class="btn btn-info pull-right" type="button" ng-click="search_analyze(notice)">
                  搜索分析
              </button>
            </div>
          </div>
        </div>
      </div>
  </div>
</script>