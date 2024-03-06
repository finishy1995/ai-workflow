package transformer

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/common"
	"github.com/finishy1995/effibot-core/ai/internal/agent/tools"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
)

type Script struct {
	common.Script
}

func NewScript() data.Item {
	return new(Script)
}

func (s *Script) Run(task *core.Task, param map[string]interface{}) error {
	task.ErrorCode = base.CodeAiAgentProcessError
	c := s.GetConfig().GetVersionByUserID(task.UserId)
	if c == nil {
		return core.ErrConfigSettingInvalid
	}
	cv, ok := c.(*ConfigVersion)
	if !ok {
		return core.ErrConfigSettingInvalid
	}

	params := tools.MergeParam(task.Param, param)
	env := getEnvConfig(param, &cv.EnvConfig)
	prompts := tools.FillParam(params, env.Prompts)
	answer := tools.FillParamForString(params, env.Answer)
	task.Input = append(task.Input, prompts...)
	task.Input = append(task.Input, answer)

	task.Output = []string{answer}
	task.ErrorCode = base.CodeOk

	return nil
}
