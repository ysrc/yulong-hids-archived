package models

import "gopkg.in/mgo.v2/bson"

type TaskForm struct {
	Name      string   `form:"name"        valid:"Required"`
	Tag       string   `form:"tag"         valid:"Required"`
	Type      string   `form:"type"        valid:"Required"`
	Command   string   `form:"command"     valid:"Required"`
	Host_List []string `form:"host_list"   valid:"Required"`
}

type StatusForm struct {
	Id     string `form:"id"    valid:"Required"`
	Info   string `form:"info"    valid:""`
	Type   string `form:"type" valid:""`
	Status int    `form:"status" valid:""`
}

type EditCfgForm struct {
	Id    string `form:"id"    valid:"Required"`
	Key   string `form:"key" valid:"Required"`
	Input string `form:"input" valid:"Required"`
}

type CodeInfo struct {
	Status int               `json:"status"`
	Msg    string            `json:"msg"`
	Data   map[string]string `json:data`
}

type NewTaskInfo struct {
	Status int           `json:"status"`
	Msg    string        `json:"msg"`
	TaskId bson.ObjectId `json:_id`
}

type SearchForm struct {
	Keyword string `json:"keyword"`
}

func NewErrorInfo(info string) *CodeInfo {
	return &CodeInfo{0, info, make(map[string]string)}
}

func NewNormalInfo(info string) *CodeInfo {
	return &CodeInfo{1, info, make(map[string]string)}
}
