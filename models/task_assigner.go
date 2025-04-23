package models

type TaskAssignConfig struct {
	Members     []Member `json:"members"`
	Requirement string   `json:"requirement"`
}

type Member struct {
	Name  string `json:"name"`
	Skill string `json:"skill"`
}

type TaskConfig struct {
	Name string `json:"name"`
	Task string `json:"task"`
}

type TaskAssignVO struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    []TaskConfig `json:"data"`
}
