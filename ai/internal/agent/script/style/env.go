package style

import "github.com/finishy1995/effibot-core/library/data"

type EnvConfig struct {
	PromptPrefix []string `json:"prompt_prefix"`
	PromptSuffix []string `json:"prompt_suffix"`
}

const (
	PromptPrefix string = "style_prompt_prefix"
	PromptSuffix string = "style_prompt_suffix"
)

var (
	defaultEnvSettings = []*data.Env{
		{
			Name:        PromptPrefix,
			Description: "风格前缀",
			Type:        data.EnvTypeStringSlice,
		},

		{
			Name:        PromptSuffix,
			Description: "风格后缀",
			Type:        data.EnvTypeStringSlice,
		},
	}
)

func getEnvConfig(param map[string]interface{}, envDefault *EnvConfig) *EnvConfig {
	result := &EnvConfig{
		PromptPrefix: envDefault.PromptPrefix,
		PromptSuffix: envDefault.PromptSuffix,
	}
	if inter, ok := param[PromptPrefix]; ok {
		if item, ok := inter.([]string); ok {
			result.PromptPrefix = item
		}
	}
	if inter, ok := param[PromptSuffix]; ok {
		if item, ok := inter.([]string); ok {
			result.PromptSuffix = item
		}
	}

	return result
}
