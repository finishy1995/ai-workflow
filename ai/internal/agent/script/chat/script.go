package chat

import (
	"context"
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/common"
	"github.com/finishy1995/effibot-core/ai/internal/agent/tools"
	"github.com/finishy1995/effibot-core/ai/internal/gpt"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
	"github.com/zeromicro/go-zero/core/logx"
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
	system := tools.FillParamForString(params, env.System)
	task.Input = append(task.Input, prompts...)

	logx.Infof("task: %+v, system: [%+v], prompts: [%+v]", task, system, task.Input)

	completion, err := gpt.ChatCompletion.CreateChatCompletion(
		context.Background(),
		system,
		task.Input,
		gpt.WithModel(env.Model),
		gpt.WithMaxToken(env.MaxToken),
		gpt.WithTemperature(env.Temperature))
	if err != nil {
		logx.Errorf("gpt chat completion error: %s, task: %+v", err.Error(), task)
		return err
	}
	task.Output = []string{completion}
	task.ErrorCode = base.CodeOk

	return nil
}
