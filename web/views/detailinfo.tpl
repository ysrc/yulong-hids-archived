<!-- detailinfo.tpl -->
<script type="text/ng-template" id="detailinfo">
<div ng-controller="detailinfo" class="detailinfo">
<div class="col-lg-12 col-md-12 col-sm-12 col-xs-12">
  <div class="row">
  <div class="col-lg-12">
    <div class="app-heading">
      <div class="app-title">
        <div class="title">IP : '<span class="highlight">{{ infolist.ip }}</span>'</div>
        <div class="description">
          显示主机的相关信息
        </div>
        <p>
            <span class="badge badge-text badge-icon" ng-repeat="p in infolist.path">
               <i class="fa fa-folder" aria-hidden="true"></i>{{ p }}
            </span>
            <span class="badge badge-text badge-icon" ng-if="infolist.hostname">
              <i class="fa fa-server" aria-hidden="true"></i>{{ infolist.hostname }}
            </span>
            <span class="badge badge-text badge-icon" ng-if="infolist.system">
              <i ng-class="infolist.system.indexOf('indows') > -1 ? 'fa fa-windows' : 'fa fa-linux'" aria-hidden="true"></i>{{ infolist.system }}
            </span>
            <span class="badge badge-text badge-icon" ng-if="infolist.type">
              <i class="fa fa-file" aria-hidden="true"></i>{{ infolist.type }}
            </span>
            <span class="badge badge-text badge-icon" ng-if="infolist.uptime">
              <i class="fa fa-clock-o" aria-hidden="true"></i>{{ infolist.uptime }}
            </span>
          </p>
      </div>
    </div>
  </div>
  </div>

  <div class="card card-tab card-mini">
    <div class="card-header">
      <ul class="nav nav-tabs tab-stats">
        <li role="{{ info.type }}" ng-repeat="info in infolist.infodata | filter: search track by $index" ng-class="{active:type == info.type}">
          <a aria-controls="{{ info.type }}" role="tab" data-toggle="tab" aria-expanded="false">{{ info.type }}</a>
        </li>
      </ul>
    </div>
    <div class="card-body tab-content">
      <div role="tabpanel" class="tab-pane" id="{{ info.type }}" ng-repeat="info in infolist.infodata" ng-class="{active:type == info.type}">
          <div class="card-body no-padding table-responsive">
            <table class="table card-table">
              <thead>
                <tr>
                  <th class="right" ng-repeat="(k,d) in info.data[0]">{{ k }}</th>
                </tr>
              </thead>
              <tbody>
                <tr ng-repeat="item in info.data | filter: search track by $index">
                  <td ng-repeat="d in item" title="{{ d }}">{{ d | cutWords:100 }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
  </div>
</div>
</div>
</script>