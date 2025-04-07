package ai

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/sirupsen/logrus"
)

type ModelConfig struct {
	Url    string
	Model  string
	ApiKey string
}

type HistoryEntry struct {
	UserContent string `json:"user_content"`
	AiContent   string `json:"ai_content"`
}

type AiQuestion struct {
	SystemPrompt string         `json:"system_prompt"`
	History      []HistoryEntry `json:"history"`
	Question     string         `json:"question"`
}

func InvokeLLM(ctx context.Context, question AiQuestion, model ModelConfig) (<-chan string, error) {
	client := openai.NewClient(
		option.WithBaseURL(model.Url),
		option.WithAPIKey(model.ApiKey),
	)

	// 构造对话
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(question.SystemPrompt),
	}
	for _, h := range question.History {
		messages = append(messages, openai.UserMessage(h.UserContent))
		messages = append(messages, openai.SystemMessage(h.AiContent))
	}
	messages = append(messages, openai.UserMessage(question.Question))

	// 构造请求参数
	params := openai.ChatCompletionNewParams{
		Model:    model.Model,
		Seed:     openai.Int(1),
		Messages: messages,
	}

	// 流式请求
	stream := client.Chat.Completions.NewStreaming(ctx, params)
	ch := make(chan string)
	go func() {
		defer close(ch)
		for stream.Next() {
			evt := stream.Current()
			if len(evt.Choices) > 0 {
				ch <- evt.Choices[0].Delta.Content
				// fmt.Print(evt.Choices[0].Delta.Content) // debug
			}
		}
	}()
	logrus.Debugf("streaming started")
	return ch, nil
}
