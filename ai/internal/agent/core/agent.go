package core

import (
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
	"github.com/finishy1995/go-library/routine"
	"github.com/zeromicro/go-zero/core/logx"
)

// AiAgent AI Agent 接口
type AiAgent interface {
	data.Item
	// RegisterHook 允许其他 AiAgent 注册 Hook，f 函数将在 typ 节点触发时异步执行
	RegisterHook(registrantId string, typ HookType, f HookFunc)
	// Process AI Agent 运行并响应请求
	Process(task *Task)

	// CanProcessingResult 当前 AI Agent 是否能作为一个 workflow 的中间 AI Agent
	CanProcessingResult() bool
	// CanFinalResult 当前 AI Agent 是否能作为一个 workflow 的结束 AI Agent
	CanFinalResult() bool
}

// BaseAiAgent Ai Agent 最基础的实现
type BaseAiAgent struct {
	data.BaseItem
	hookMap map[HookType]map[string]HookFunc
}

func NewBaseAiAgent() data.Item {
	return new(BaseAiAgent)
}

func (b *BaseAiAgent) Init(config data.Config) bool {
	b.BaseItem.Init(config)
	b.hookMap = map[HookType]map[string]HookFunc{}

	return true
}

func (b *BaseAiAgent) RegisterHook(registrantId string, typ HookType, f HookFunc) {
	if m, ok := b.hookMap[typ]; ok {
		m[registrantId] = f
	} else {
		b.hookMap[typ] = map[string]HookFunc{registrantId: f}
	}
}

func (b *BaseAiAgent) TriggerHook(typ HookType, task *Task) {
	if m, ok := b.hookMap[typ]; ok {
		for _, f := range m {
			err := routine.Run(true, func() {
				f(task)
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func (b *BaseAiAgent) Process(task *Task) {
	c := b.GetConfig().GetVersionByUserID(task.UserId)
	if c == nil {
		task.ErrorCode = base.CodeAiAgentProcessError
		return
	}
	acv, ok := c.(*AgentConfigVersion)
	if !ok {
		task.ErrorCode = base.CodeAiAgentProcessError
		return
	}

	if task.Input == nil {
		task.Input = make([]string, 0)
	}
	task.ErrorCode = base.CodeOk

	b.TriggerHook(BeforeProcess, task)
	err := acv.script.Run(task, acv.Env)
	if err != nil {
		logx.Errorf("Agent [%s] process failed, error: %s", b.GetID(), err.Error())
		task.ErrorCode = base.CodeAiAgentProcessError
		return
	}
	b.TriggerHook(AfterProcess, task)
}

func (b *BaseAiAgent) CanProcessingResult() bool {
	return false
}

func (b *BaseAiAgent) CanFinalResult() bool {
	return true
}
