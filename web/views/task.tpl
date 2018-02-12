<!-- task.tpl -->

<script type="text/ng-template" id="tasks">
<div ng-controller="tasks">
    <div class="col-md-12 top-nav">
        <div class="col-md-6 left-title">
            任务管理
        </div>
    </div>
    <div class="col-md-12">
      <div class="card">
        <div class="card-header">
          任务列表
        </div>
        <div class="card-body no-padding">
          <div class="dataTables_wrapper form-inline dt-bootstrap no-footer">
          <table class="datatable table table-striped primary dataTable no-footer" cellspacing="0" role="grid" id="tasktable">
            <thead>
            <tr role="row">
                <th ng-repeat="k in table_key_list">{{ k | uppercase }}</th>
                <th><i class="fa fa-cogs" aria-hidden="true"></i></th>
            </tr>
            </thead>
            <tbody>
                <tr role="row" class="odd" ng-repeat="task in tasks" id="{{ task._id }}">
                    <td ng-repeat="key in table_key_list" title="{{ task[key].toString() }}" ng-if="!key.startsWith('_')">
                        {{ key == 'time' ? timeformat(task[key]) : task[key].toString() }}
                    </td>
                    <td>
                      <a type="button" class="mb-sm btn-xs btn btn-primary" href="/#!/taskresult/{{ task._id }}/">查看结果</a>
                  </td>
                </tr>
            </tbody>
            </table>

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
    <button class="icon" id="fixedbutton" data-toggle="modal" data-target="#newTaskModel">
        <i class="fa fa-plus" aria-hidden="true"></i>
    </button>
</div>
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
          <button type="button" class="btn btn-default" data-dismiss="modal" id="model-close">关闭</button>
          <button type="button" class="btn btn-info" id="newtask" onclick="newtask()">
              提交
          </button>
        </div>
      </div>
    </div>
</div>
</script>

