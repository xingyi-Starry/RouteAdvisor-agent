package config

import (
	"route-advisor-agent/utils"

	"github.com/sirupsen/logrus"
)

type EnvConfig struct {
	BindAddr     string
	ApiUrl       string
	ApiKey       string
	DefaultModel string
}

func LoadEnvConfig() *EnvConfig {
	ans := &EnvConfig{
		BindAddr:     utils.ParseEnvConfig("PORT", "localhost:8000"),
		ApiUrl:       utils.ParseEnvConfig("API_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1"),
		ApiKey:       utils.ParseEnvConfig("API_KEY", ""),
		DefaultModel: utils.ParseEnvConfig("DEFAULT_MODEL", "qwen-turbo"),
	}
	if ans.ApiKey == "" {
		logrus.Fatal("env API_KEY is required!")
	}
	return ans
}
