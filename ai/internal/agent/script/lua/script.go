package lua

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/common"
	"github.com/finishy1995/effibot-core/ai/internal/agent/tools"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
	lua "github.com/yuin/gopher-lua"
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

	if cv.Script != "" {
		l := lua.NewState()
		defer l.Close()

		tools.SetLuaData(l, task, param)
		// 创建 AgentFunc 表
		agentFuncTable := l.NewTable()
		l.SetField(agentFuncTable, "SimpleGenerate", l.NewFunction(tools.SimpleGenerate))
		// 将 AgentFunc 表设置为全局变量
		l.SetGlobal("AgentFunc", agentFuncTable)

		// TODO: 脚本执行安全
		if err := l.DoString(cv.Script); err != nil {
			return err
		} else {
			task.ErrorCode = base.CodeOk
		}
	} else {
		return s.Script.Run(task, param)
	}
	return nil
}
