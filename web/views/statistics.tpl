<script type="text/ng-template" id="statistics">
  <div ng-controller="statistics">
    <div class="row tile_count">
      <div class="col-md-2 col-sm-4 col-xs-6 tile_stats_count">
        <span class="count_top">
        <i class="fa fa-user"></i>主机数（在线/全部）
      </span>
        <div class="count">
          <span class="green">{{totaldata.hostdata[0]}}</span> / 
          <span class="yellow">{{totaldata.hostdata[1]}}</span>
        </div>
      </div>
      <div class="col-md-2 col-sm-4 col-xs-6 tile_stats_count">
        <span class="count_top">
        <i class="fa fa-clock-o"></i>告警（确认/全部/忽略）
      </span>
        <div class="count">
          <span class="green">{{totaldata.alarmdata[0]}}</span> / 
          <span class="yellow">{{totaldata.alarmdata[1]}}</span> / 
          <span class="red">{{totaldata.alarmdata[2]}}</span>
        </div>
      </div>
      <div class="col-md-2 col-sm-4 col-xs-6 tile_stats_count">
        <span class="count_top">
        <i class="fa fa-user"></i>任务（执行/全部/失败）
      </span>
        <div class="count">
          <span class="green">{{totaldata.taskdata[0]}}</span> /
          <span class="yellow">{{totaldata.taskdata[1]}}</span> /
          <span class="red">{{totaldata.taskdata[2]}}</span>
        </div>
      </div>
      <div class="col-md-2 col-sm-4 col-xs-6 tile_stats_count">
        <span class="count_top">
        <i class="fa fa-user"></i>服务端（在线/负载）
      </span>
        <div class="count green">
          <span class="green">{{totaldata.servicedata[0]}}</span> /
          <span class="yellow">{{totaldata.servicedata[1]}}%</span>
        </div>
      </div>
      <div class="col-md-3 col-sm-6 col-xs-6 tile_stats_count">
        <span class="count_top">
        <i class="fa fa-user"></i>数据总览（信息/行为）
      </span>
        <div class="count green">
          <span class="green">{{totaldata.totaldata[0]}}</span> /
          <span class="yellow">{{totaldata.totaldata[1]}}</span>
        </div>
      </div>
    </div>
    <div class="row">
      <div class="row col-md-7">
        <div class="col-md-12 ">
          <div class="x_panel tile  overflow_hidden" style="height: 300px;">
            <div class="x_title">
              <h2>告警日期分布</h2>
              <div class="clearfix"></div>
            </div>
            <div class="x_content" style="margin: 0px 10px 10px 0px; width: 100%; height: 220px;" id="chartLine">
            </div>
          </div>
        </div>
        <div class="col-md-12">
          <div class="x_panel tile overflow_hidden" id="div-time-chart">
            <div class="x_title">
              <h2>告警时间分布</h2>
              <div class="clearfix"></div>
            </div>
            <div class="x_content" id="chartTwo">
            </div>
          </div>
        </div>
      </div>
      <div class="row col-md-5">

        <div class="col-md-12">
          <div class="x_panel tile overflow_hidden" style="height: 340px;">
            <div class="x_title">
              <h2>告警类型统计</h2>
              <div class="clearfix"></div>
            </div>
            <div class="x_content">
              <div id="chartOne" style="margin: 0px 10px 10px 0px; width: 450px; height: 270px;"></div>
            </div>
          </div>
        </div>
        <div class="col-md-12 top-notice">
          <div class="x_panel" style="height: 390px;">
            <div class="x_title">
              <h2>告警信息 </h2>
              <div class="clearfix"></div>
            </div>
            <div class="x_content">
              <ol class="group full-height" id="olmsgList">
                <li class="message" ng-repeat="msg in data.listdata">
                  <a href="/#!/notice" target="_blank" rel="noopener noreferrer">
                  <label>{{ msg[0] | cutWords:25 }} ...</label>
                  <span class="badge badge-danger badge-icon">
                    <i class="fa fa-calendar" aria-hidden="true"></i>{{ msg[3] }}
                  </span>
                  <span class="badge badge-info badge-icon">
                    <i class="fa fa-check-square-o" aria-hidden="true"></i>{{ msg[1] }}
                  </span>
                  <span class="badge {{ style.notice.level[msg[2]] }} badge-icon">
                      <i class="fa fa-exclamation-triangle" aria-hidden="true"></i>
                      {{ langtem.notice.data.level[msg[2]] }}
                  </span>
                  </td>
                </a>
                </li>
              </ol>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</script>