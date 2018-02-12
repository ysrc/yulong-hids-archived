package models

type TaskResult struct {
	baseModel
}

func NewTaskResult() TaskResult {
	mdl := TaskResult{}
	mdl.collectionName = "task_result"
	return mdl
}
