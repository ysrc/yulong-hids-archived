<!-- rule.tpl -->
<script type="text/ng-template" id="rules">
    <div ng-controller="rules" class="rules">

        <div class="col-md-12 top-nav">
            <div class="col-md-6 left-title">
                规则引擎
            </div>
            <div class="col-md-6">
                <a type="button" href="{{ rules_url }}?action=download" class="btn btn-info btn-square pull-right" target="_blank" rel="noopener noreferrer">
                    <i class="fa fa-cloud-download" aria-hidden="true"></i>
                    导出规则
                </a>
                <button class="btn btn-primary btn-square pull-right" href="#model-add-rules" data-toggle="modal" onclick="$('h4.modal-title strong').text('新增规则（多条规则请用json列表）');">
                    <i class="fa fa-plus" aria-hidden="true"></i>
                    添加规则
                </button>
            </div>
        </div>
    
        <div class="modal fade" id="model-add-rules">
            <div class='modal-dialog'>
                <div class='modal-content'>
                    <div class='modal-header'>
                        <button type="button" class="close" data-dismiss="modal" aria-hidden="true">×</button>
                        <h4 class='modal-title'>
                            <strong>新增规则（多条规则请用json列表）</strong>
                        </h4>
                    </div>
                    <div class='modal-body'>
                        <textarea class="jsoncode" ng-model="new_rules"></textarea>
                    </div>
                    <div class='modal-footer'>
                        <button class="btn btn-primary btn-square pull-right" ng-click="add_rules()">
                            <i class="fa fa-cloud" aria-hidden="true"></i>
                            保存规则
                        </button>
                    </div>
                </div>
            </div>
        </div>


        <div class="col-md-12">
            <div class="panel-group" id="accordion" role="tablist" aria-multiselectable="true">

                <div class="panel panel-rules" ng-repeat="rule in rulelist track by $index">

                    <div class="panel-heading" role="tab" id="head-{{ rule._id }}">
                        <h4 class="panel-title">
                            <a style="cursor:pointer;" ng-click="open_body(rule._id)">{{ rule.meta.name }}</a>

                            <span class="{{ style.notice.level[rule.meta.level] }} diff">
                                <i class="fa fa-exclamation-circle" aria-hidden="true"></i>
                                {{ langtem.notice.data.level[rule.meta.level] }}
                            </span>
                            <span class="badge badge-info">
                                <i class="fa fa-comment" aria-hidden="true"></i>
                                {{ rule.meta.description }}
                            </span>
                            <span class="badge badge-success">
                                <i class="fa fa-desktop" aria-hidden="true"></i>
                                {{ rule.system }}
                            </span>
                            <span class="badge badge-info">
                                    <i class="fa fa-pencil-square" aria-hidden="true"></i>
                                    {{ rule.meta.author }}
                            </span>
                            <span class="badge badge-success">
                                    <i class="fa fa-bookmark-o" aria-hidden="true"></i>
                                    {{ rule.source }}
                            </span>
                            <span class="badge badge-danger" style="cursor: pointer" ng-click="delete(rule._id)">
                                <i class="fa fa-trash" aria-hidden="true"></i>
                                删除
                            </span>
                            <span class="badge badge-warning" style="cursor: pointer" ng-click="edit($index)">
                                <i class="fa fa-pencil" aria-hidden="true"></i>
                                编辑
                            </span>
                            <a role="button">
                                <i class="more-less fa fa-plus" ng-click="open_body(rule._id)"></i>
                            </a>
                            <div class="material-switch pull-right" title="是否启用该规则" ng-click="enable_toggle($index)">
                                <input type="checkbox" ng-checked="rule.enabled"/>
                                <label class="label-primary"></label>
                            </div>
                        </h4>
                    </div>

                    <div class="panel-collapse collapse" role="tabpanel" id="panel-{{ rule._id }}">
                        <div class="panel-body" id="body-{{ rule._id }}">
                            <div class="jsoncode" readonly>{{ pettyprint(rule, null, 2) }}</div>
                        </div>
                    </div>

                </div>

            </div><!-- panel-group -->
        </div>
    </div>
</script>