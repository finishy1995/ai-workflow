package workflow

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
)

type AiWorkflow interface {
	data.Item
	// Process AI Workflow 运行并响应请求
	Process(task *core.Task)
}

// BaseAiWorkflow Ai Workflow 最基础的实现
type BaseAiWorkflow struct {
	data.BaseItem
}

func NewBaseAiWorkflow() data.Item {
	return new(BaseAiWorkflow)
}

func (b *BaseAiWorkflow) Process(task *core.Task) {
	c := b.GetConfig().GetVersionByUserID(task.UserId)
	if c == nil {
		task.ErrorCode = base.CodeAiWorkflowProcessError
		return
	}
	acv, ok := c.(*ConfigVersion)
	if !ok {
		task.ErrorCode = base.CodeAiWorkflowProcessError
		return
	}

	if task.Input == nil {
		task.Input = make([]string, 0)
	}
	task.ErrorCode = base.CodeOk

	for _, agent := range acv.agents {
		agent.Process(task)
		if task.ErrorCode != base.CodeOk {
			return
		}
	}
}
