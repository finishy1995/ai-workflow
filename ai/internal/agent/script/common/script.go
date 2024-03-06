package common

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
)

type Script struct {
	data.BaseItem
}

func NewScript() data.Item {
	return new(Script)
}

func (s *Script) Run(task *core.Task, _ map[string]interface{}) error {
	task.ErrorCode = base.CodeOk
	task.Output = []string{task.Input[len(task.Input)-1]}

	return nil
}
