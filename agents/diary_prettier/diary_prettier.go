package diary_prettier

import (
	"context"
	"route-advisor-agent/utils/ai"

	"github.com/sirupsen/logrus"
)

const prompt = `你是一个游记写手，负责为给定的游记内容进行润色，或者是根据片段补完全文。请根据以下游记内容进行润色，确保语句通顺、逻辑清晰，并且增加一些细节和情感色彩。请确保生成的游记内容与原始内容相关，并且要有一定的逻辑性和连贯性。`

type DiaryPrettier struct {
	ModelConfig *ai.ModelLLM
}

type DiaryConfig struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewDiaryPrettier(modelConfig *ai.ModelLLM) *DiaryPrettier {
	return &DiaryPrettier{
		ModelConfig: modelConfig,
	}
}

func (agent *DiaryPrettier) PrettifyDiary(ctx context.Context, config DiaryConfig) (string, error) {
	question := ai.Messages{
		SystemPrompt: prompt,
		History:      []ai.HistoryEntry{},
		Question:     "标题：" + config.Title + "\n内容" + config.Content,
	}

	// 调用 LLM 接口
	ch, err := ai.InvokeLLM(ctx, question, agent.ModelConfig)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	// 处理返回的结果
	result := ""
	for msg := range ch {
		result += msg
	}
	logrus.Info(result)

	return result, nil
}
