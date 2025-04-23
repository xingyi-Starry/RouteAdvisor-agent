package task_assigner

import (
	"context"
	"encoding/json"
	"route-advisor-agent/models"
	"route-advisor-agent/utils/ai"

	"github.com/sirupsen/logrus"
)

// 纯纯课程幽默需求，真有人会把系统架构、任务分配这样的活交给ai干吗

const prompt = `你是一名资深程序架构师，负责根据程序需求和团队成员的技能，合理设计使用的技术栈和分配任务。请根据以下信息，优先考虑前后端分离架构，给出合理的任务分配方案。如果没有符合的技术栈，就在任务中加入技术学习任务。请使用JSON格式返回结果，每一个成员包含"name"和"task"两个字段。注意，实际的输入和输出会远比示例更复杂，请给出更细节的任务安排，如模块任务划分，模块完成时间等等。确保输出的 JSON 格式正确，不要添加任何额外的文本或解释，不要用markdown的代码块包裹输出。
<输入示例>
{"members":[{"name":"Alice","skill":"精通Java"},{"name":"Bob","skill":"熟练掌握Python"},{"name":"Charlie","skill":"掌握JavaScript和前端技术"}],"requirement":"开发一个在线购物平台"}
<输入示例 />
<对应输出示例>
[{"name":"Alice","task":"负责后端服务的设计和实现，使用Java和Spring Boot框架。"},{"name":"Bob","task":"负责数据处理和分析，使用Python进行数据挖掘和机器学习。"},{"name":"Charlie","task":"负责前端页面的设计和实现，使用JavaScript和React框架。"}]
<对应输出示例 />
`

type TaskAssigner struct {
	model ai.LLM
}

func NewTaskAssigner(model ai.LLM) *TaskAssigner {
	return &TaskAssigner{
		model: model,
	}
}

func (agent *TaskAssigner) AssignTask(ctx context.Context, config models.TaskAssignConfig) ([]models.TaskConfig, error) {

	jsonConfig, err := json.Marshal(config)
	if err != nil {
		logrus.Error(err)
		return []models.TaskConfig{}, err
	}

	question := ai.Messages{
		SystemPrompt: prompt,
		History:      []ai.HistoryEntry{},
		Question:     string(jsonConfig),
	}

	// 调用 LLM 进行推理
	questionChan := ai.InvokeLLM(ctx, question, agent.model)

	answer := ""
	for msg := range questionChan {
		answer += msg
	}
	logrus.Debug(answer) // debug
	result := []models.TaskConfig{}
	err = json.Unmarshal([]byte(answer), &result)
	if err != nil {
		logrus.Error(err)
		return []models.TaskConfig{}, err
	}

	return result, nil
}
