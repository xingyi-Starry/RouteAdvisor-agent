package main

import (
	"route-advisor-agent/agents/diary_writer"
	"route-advisor-agent/config"
	"route-advisor-agent/router"
	"route-advisor-agent/services"
	"route-advisor-agent/utils"
	"route-advisor-agent/utils/ai"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	utils.InitLogger()
	app := fiber.New()

	env := config.LoadEnvConfig()
	model := InitModel(env)
	controller := InitController(model)
	router.InitRouter(app, controller)

	err := app.Listen(env.BindAddr)
	if err != nil {
		logrus.Fatal("Failed to start server:", err)
	}
}

// 初始化模型配置
func InitModel(env *config.EnvConfig) ai.LLM {
	modelConfig := ai.DefaultLLMConfig()
	modelConfig.Url = env.ApiUrl
	modelConfig.ModelName = env.DefaultModel
	modelConfig.ApiKey = env.ApiKey
	modelConfig.MaxConn = 10
	modelConfig.MaxQps = 10
	modelConfig.MaxQueue = 10
	model := ai.NewModelLLM(modelConfig)
	model.Run()
	return model
}

func InitController(model ai.LLM) *services.Controller {
	diaryWriter := services.NewDiaryWriterService(diary_writer.NewDiaryWriter(model))
	controller := services.NewController(diaryWriter)
	return controller
}
