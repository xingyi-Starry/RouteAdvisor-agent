package diary_writer

import (
	"context"
	"route-advisor-agent/utils/ai"

	"github.com/sirupsen/logrus"
)

const prompt = `你是一个游记写手，负责为给定的标题生成一篇游记。请根据以下标题生成一篇游记，游记的内容应该包括景点、活动、饮食等方面的信息，并且要有一定的情感色彩。游记的长度应该在500字左右。请确保生成的游记内容与标题相关，并且要有一定的逻辑性和连贯性。`

type DiaryWriter struct {
	ModelConfig *ai.ModelConfig
}

type DiaryConfig struct {
	Title      string `json:"title"`
	PictureNum int    `json:"picture_num"`
}

func NewDiaryWriter(modelConfig *ai.ModelConfig) *DiaryWriter {
	return &DiaryWriter{
		ModelConfig: modelConfig,
	}
}

func (agent *DiaryWriter) CreateDiary(ctx context.Context, config DiaryConfig) (string, error) {
	question := ai.AiQuestion{
		SystemPrompt: prompt,
		History:      []ai.HistoryEntry{},
		Question:     config.Title,
	}

	// 调用 LLM 接口
	ch, err := ai.InvokeLLM(ctx, question, *agent.ModelConfig)
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
