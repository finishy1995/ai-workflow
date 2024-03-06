package style

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

	if styleInterface, ok := task.Param[core.AIKeyStyle]; ok {
		style, ok := styleInterface.(string)
		if !ok {
			task.ErrorCode = base.CodeOk
			return nil
		}

		if typ, ok := cv.styleTypeMap[style]; ok {
			cType := typ.GetConfig().GetVersionByUserID(task.UserId)
			if cType == nil {
				return core.ErrConfigSettingInvalid
			}
			tcv, ok := cType.(*TypeConfigVersion)
			if !ok {
				return core.ErrConfigSettingInvalid
			}

			task.ErrorCode = base.CodeOk

			params := tools.MergeParam(task.Param, param)
			env := getEnvConfig(params, &cv.EnvConfig)

			results := append(env.PromptPrefix, tcv.Prompt...)
			results = append(results, env.PromptSuffix...)
			task.Input = append(task.Input, tools.FillParam(params, results)...)
		}
	}
	return nil
}
