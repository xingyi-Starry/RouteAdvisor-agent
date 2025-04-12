package ai

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

type ModelLLM struct {
	Client    *openai.Client
	ModelName string
	Sem       chan struct{}
	Limiter   rate.Limiter
	Queue     chan task
}

type ModelLLMConfig struct {
	Url       string
	ModelName string
	ApiKey    string
	MaxConn   int
	MaxQps    float64
	MaxQueue  int
}
type HistoryEntry struct {
	UserContent string `json:"user_content"`
	AiContent   string `json:"ai_content"`
}

type Messages struct {
	SystemPrompt string
	History      []HistoryEntry
	Question     string
	Seed         int64
}

type task struct {
	Msg Messages
	Ch  chan string
	Ctx context.Context
}

func DefaultLLMConfig() *ModelLLMConfig {
	return &ModelLLMConfig{
		Url:       "https://dashscope.aliyuncs.com/compatible-mode/v1",
		ModelName: "qwen-turbo",
		ApiKey:    "",
		MaxConn:   10,
		MaxQps:    10,
		MaxQueue:  100,
	}
}

func NewModelLLM(config *ModelLLMConfig) *ModelLLM {
	sem := make(chan struct{}, config.MaxConn)
	for range config.MaxConn {
		sem <- struct{}{}
	}

	client := openai.NewClient(
		option.WithBaseURL(config.Url),
		option.WithAPIKey(config.ApiKey),
	)
	return &ModelLLM{
		Client:    &client,
		ModelName: config.ModelName,
		Sem:       sem,
		Limiter:   *rate.NewLimiter(rate.Limit(config.MaxQps), max(int(config.MaxQps), 1)),
		Queue:     make(chan task, config.MaxQueue),
	}
}

func (model *ModelLLM) Run() {
	go func() {
		for {
			t := <-model.Queue
			// 先控制频率
			if err := model.Limiter.Wait(t.Ctx); err != nil {
				logrus.Error(err)
				continue
			}

			// 再控制并发
			<-model.Sem

			logrus.Debug("task started")

			// 执行请求
			go func(t task) {
				defer func() {
					logrus.Debug("task finished")
					model.Sem <- struct{}{}
				}() // 释放信号量
				defer close(t.Ch) // 关闭通道

				// 构造对话
				messages := []openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(t.Msg.SystemPrompt),
				}
				for _, h := range t.Msg.History {
					messages = append(messages, openai.UserMessage(h.UserContent))
					messages = append(messages, openai.SystemMessage(h.AiContent))
				}
				messages = append(messages, openai.UserMessage(t.Msg.Question))

				// 构造请求参数
				params := openai.ChatCompletionNewParams{
					Model:    model.ModelName,
					Seed:     openai.Int(t.Msg.Seed),
					Messages: messages,
				}

				// 流式请求
				stream := model.Client.Chat.Completions.NewStreaming(t.Ctx, params)
				for stream.Next() {
					evt := stream.Current()
					if len(evt.Choices) > 0 {
						t.Ch <- evt.Choices[0].Delta.Content
						// fmt.Print(evt.Choices[0].Delta.Content) // debug
					}
				}
			}(t)
		}
	}()
}

func InvokeLLM(ctx context.Context, msgs Messages, model *ModelLLM) <-chan string {
	ch := make(chan string, 5) // 响应通道

	// 任务加入队列
	t := task{
		Msg: msgs,
		Ch:  ch,
		Ctx: ctx,
	}
	model.Queue <- t

	return ch
}
