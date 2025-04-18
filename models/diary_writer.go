package models

type DiaryConfig struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DiaryResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    DiaryConfig `json:"data"`
}
