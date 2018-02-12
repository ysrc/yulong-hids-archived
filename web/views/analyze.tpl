<!-- analyze.tpl -->

<script type="text/ng-template" id="analyze">
  <div ng-controller="analyze" class="analyze">
    
    <nav class="navbar">
      <div class="app-heading">
        <div class="app-title">
          <div class="title">数据分析<span class="highlight ng-binding">&</span>溯源</div>
          <div class="description">可在搜索框内直接输入检索语句，格式为 "<span class="highlight">type:typename | key:value</span>"。使用 "<span class="highlight">TAB</span>" 键可获得搜索提示。
          <p>
            <span class="badge badge-icon badge-text" ng-repeat="h in history track by $index" ng-if="h != ''" ng-click="re_search(h)">
              <i class="fa fa-history" aria-hidden="true"></i>{{ h }}
            </span>
          </p>
          <p class="search-sample">
            {{ search_start }} ( {{ search_sample.msg }} : <span class="highlight" ng-repeat="s in search_sample.samplelist" ng-click="append2search(s)"> {{ s }} </span> )
          </p>
          </div>
        </div>
      </div>

      <div class="navbar-collapse collapse in">
        <ul class="nav navbar-nav navbar-left col-md-12">
          <li class="navbar-search">
            <input id="search" type="text" value="type:" autocomplete="off" autofocus>
            <button class="btn-search" ng-click="show_data()"><i class="fa fa-search"></i></button>
          </li>
        </ul>
        <div class="pull-right time-select">
          <span moment-picker="ctrl.minday" format="YYYY-MM-DD HH:mm">
            {{ ctrl.minday || '点击选择开始日期' }}
          </span>
          -
          <span moment-picker="ctrl.maxday" format="YYYY-MM-DD HH:mm">
            {{ ctrl.maxday || '点击选择结束日期' }}
          </span>
        </div>
      </div>
    </nav>

    <div class="col-md-12 analyze-result">
      <div class="card">
        <div class="card-header">
          检索结果
        </div>
        <div class="card-body no-padding">
          <div class="dataTables_wrapper form-inline dt-bootstrap no-footer table-responsive">
          <table class="datatable table table-striped primary dataTable no-footer" cellspacing="0" id="result">
              <thead>
                <tr class="row">
                  <th ng-repeat="key in thead_root_items">{{ key | uppercase }}</th>
                  <th ng-if="analyze[0]['data']">DATA</th>
                </tr>
              </thead>
              <tbody>
                <!-- 适配es数据和mongodb info数据的情况 -->
                <tr ng-repeat="data in analyze track by $index" class="row" ng-if="data['data'] && data['data'].length != 0">
                  <td ng-repeat="key in thead_root_items" title="{{ data[key] }}">
                    {{ key == 'uptime' || key == "time" ? timeformat(data[key]) : data[key] }}
                  </td>
                  <td>
                    <span class="badge badge-icon" ng-class="$index % 2 == 0 ? 'badge-info':'badge-success' " ng-repeat="key in thead_data_items track by $index" title="{{ data.data[key] }}">
                        <i class="fa" aria-hidden="true" ng-if="data_is_array">{{ key | uppercase }}</i>
                          {{ data.data[0][key] | cutWords:70 }}
                          <i class="fa" aria-hidden="true" ng-if="!data_is_array">{{ key | uppercase }}</i>
                          {{ data.data[key] | cutWords:70 }}
                    </span>
                    <span href="#model-sub-td" data-toggle="modal" class="badge badge-icon badge-primary" ng-click="click_more($index)" ng-if="data_is_array">
                      <i class="fa" aria-hidden="true">action</i> see more
                    </span>
                  </td>
                </tr>

                <!-- 适配mongodb statistics数据库 -->
                <tr ng-repeat="data in analyze track by $index" class="row" ng-if="data['server_list']">
                    <td> {{ data['count'] }} </td>
                    <td> {{ data['info'] }} </td>
                    <td>
                      <span class="badge badge-icon badge-success" ng-repeat="ipaddr in data['server_list']">
                          <i class="fa fa-map-marker" aria-hidden="true"></i>
                        {{ ipaddr }}
                      </span>
                    </td>
                    <td> {{ data['type'] }} </td>
                    <td> {{ timeformat(data['uptime']) }} </td>
                </tr>
              </tbody>
            </table>
            <div class="modal fade" id="model-sub-td">
              <div class="modal-dialog">
                <div class="modal-content">
                  <div class="modal-body">
                    <table class="table table-striped" id="tblGrid">
                      <thead id="tblHead">
                        <tr>
                          <th ng-repeat="key in thead_data_items">{{ key | uppercase }}</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr ng-repeat="item in sub_td_data">
                          <td ng-repeat="i in item">
                            {{ i }}
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            </div>

            <div class="bottom">
              <div class="dataTables_paginate paging_simple_numbers">
                <ul class="pagination">
                  <li class="paginate_button previous" ng-class="current_page==1 ? 'disabled':''" ng-click="previous()">
                    <a data-dt-idx="0" tabindex="0" >Previous</a></li>
                  <li class="paginate_button active">
                    <a data-dt-idx="1" tabindex="0">{{ current_page }}</a></li>
                  <li class="paginate_button next" ng-click="next()">
                    <a data-dt-idx="7" tabindex="0">Next</a></li>
                </ul>
              </div>
              <div class="clear"></div>
            </div>
            
          </div>
        </div>
      </div>
    </div>

  </div>
</script>