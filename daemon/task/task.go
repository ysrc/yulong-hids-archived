package task

import (
	"encoding/json"
	"os"
	"yulong-hids/daemon/common"
)

// Task 接收任务结构
type Task struct {
	Type    string            // 任务类型
	Command string            // 任务内容
	Result  map[string]string // 返回结果
}

func (t *Task) Run() []byte {
	switch t.Type {
	case "reload":
		t.reload()
	case "quit":
		t.quit()
	case "kill":
		t.kill()
	case "uninstall":
		t.uninstall()
	case "update":
		t.update()
	case "delete":
		t.delete()
		// case "exec":
		// 	t.exec()
	}
	var sendResult []byte
	if b, err := json.Marshal(t.Result); err == nil {
		msg := string(b) + "\n"
		sendResult = []byte(msg)
	}
	return sendResult
}
func (t *Task) reload() {
	t.Result["status"] = "true"
	if err := common.KillAgent(); err != nil {
		t.Result["status"] = "false"
		t.Result["data"] = err.Error()
	}
}
func (t *Task) quit() {
	if common.AgentStatus {
		common.Cmd.Process.Kill()
	}
	panic(1)
}
func (t *Task) kill() {
	if redata := KillProcess(t.Command); redata != "" {
		t.Result["status"] = "true"
		t.Result["data"] = redata
	}
}
func (t *Task) uninstall() {
	UnInstallALL()
}
func (t *Task) update() {
	if ok, err := agentUpdate(common.ServerIP, common.InstallPath, common.Arch); err == nil {
		if ok {
			t.Result["status"] = "true"
			t.Result["data"] = "更新完毕"
		} else {
			t.Result["status"] = "true"
			t.Result["data"] = "已经是最新版本"
		}
	} else {
		t.Result["data"] = err.Error()
	}
}
func (t *Task) delete() {
	if err := os.Remove(t.Command); err == nil {
		t.Result["status"] = "true"
		t.Result["data"] = t.Command + " 删除成功"
	} else {
		t.Result["data"] = err.Error()
	}
}
func (t *Task) exec() {
	if dat, err := common.CmdExec(t.Command); err == nil {
		t.Result["status"] = "true"
		t.Result["data"] = dat
	} else {
		t.Result["data"] = err.Error()
	}
}
