<!-- file.tpl -->
<script type="text/ng-template" id="file">
<div ng-controller="file">
  <div class="col-md-12 col-lg-12 col-xl-12">
      <div class="card">
        <div class="card-header">
          文件管理
        </div>
        <div class="card-body no-padding">
          <div id="DataTables_Table_0_wrapper" class="dataTables_wrapper form-inline dt-bootstrap no-footer">
          <table class="datatable table table-striped primary dataTable no-footer" cellspacing="0" id="DataTables_Table_0" role="grid" aria-describedby="DataTables_Table_0_info" id="tasktable">
            <thead>
            <tr role="row">
                <th ng-repeat="(key,item) in files[0]" ng-if="!key.startsWith('_')">{{ key | uppercase}}</th>
                <th>UPDATE</th>
            </tr>
            </thead>
            <tbody>
                <tr role="row" class="odd" ng-repeat="file in files" id="{{ file._id }}">
                    <td>{{ file.platform }}</td>
                    <td>{{ file.system }}</td>
                    <td>{{ file.type }}</td>
                    <td>{{ file.hash }}</td>
                    <td>{{ file.uptime | date:'yyyy-MM-dd' }}</td>
                    <td>
                      <div class="fileUpload btn btn-info">
                          <span>{{ upload_text }}</span>
                          <input type="file" class="fileele" onchange="change($(this).attr('index'))" index="{{ $index }}" id="file{{ $index }}"/>
                      </div>
                    </td>
                </tr>
            </tbody>
            </table>
            </div>
        </div>
      </div>
  </div>
</div>
</script>