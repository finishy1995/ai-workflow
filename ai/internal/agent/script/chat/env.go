package chat

import "github.com/finishy1995/effibot-core/library/data"

type EnvConfig struct {
	MaxToken    int      `json:"max_token"`
	Temperature float32  `json:"temperature"`
	Model       string   `json:"model"`
	System      string   `json:"system"`
	Prompts     []string `json:"prompts"`
}

const (
	MaxToken    string = "chat_max_token"
	Temperature string = "chat_temperature"
	Model       string = "chat_model"
	System      string = "chat_system"
	Prompts     string = "chat_prompts"
)

var (
	defaultEnvSettings = []*data.Env{
		{
			Name:        MaxToken,
			Description: "最大返回值长度",
			Type:        data.EnvTypeInt,
		},

		{
			Name:        Temperature,
			Description: "返回值温度",
			Type:        data.EnvTypeFloat,
		},

		{
			Name:        Model,
			Description: "GPT 模型",
			Type:        data.EnvTypeString,
		},

		{
			Name:        System,
			Description: "系统设定",
			Type:        data.EnvTypeString,
		},

		{
			Name:        Prompts,
			Description: "提示词",
			Type:        data.EnvTypeStringSlice,
		},
	}
)

func getEnvConfig(param map[string]interface{}, envDefault *EnvConfig) *EnvConfig {
	result := &EnvConfig{
		MaxToken:    envDefault.MaxToken,
		Temperature: envDefault.Temperature,
		Model:       envDefault.Model,
		System:      envDefault.System,
		Prompts:     envDefault.Prompts,
	}
	if inter, ok := param[MaxToken]; ok {
		if item, ok := inter.(int); ok {
			result.MaxToken = item
		}
	}
	if inter, ok := param[Temperature]; ok {
		if item, ok := inter.(float32); ok {
			result.Temperature = item
		}
	}
	if inter, ok := param[Model]; ok {
		if item, ok := inter.(string); ok {
			result.Model = item
		}
	}
	if inter, ok := param[System]; ok {
		if item, ok := inter.(string); ok {
			result.System = item
		}
	}
	if inter, ok := param[Prompts]; ok {
		if item, ok := inter.([]string); ok {
			result.Prompts = item
		}
	}

	return result
}
