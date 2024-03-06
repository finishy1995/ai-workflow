package transformer

import "github.com/finishy1995/effibot-core/library/data"

type EnvConfig struct {
	Prompts []string `json:"prompts"`
	Answer  string   `json:"answer"`
}

const (
	Prompts string = "chat_prompts"
	Answer  string = "script_answer"
)

var (
	defaultEnvSettings = []*data.Env{
		{
			Name:        Answer,
			Description: "回答",
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
		Answer:  envDefault.Answer,
		Prompts: envDefault.Prompts,
	}
	if inter, ok := param[Answer]; ok {
		if item, ok := inter.(string); ok {
			result.Answer = item
		}
	}
	if inter, ok := param[Prompts]; ok {
		if item, ok := inter.([]string); ok {
			result.Prompts = item
		}
	}

	return result
}
