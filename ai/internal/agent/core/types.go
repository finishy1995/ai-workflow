package core

// HookFunc AI Agent 钩子函数
type HookFunc func(task *Task)

// Role 代表 AI Agent 在工作流（AI Workflow）中的角色
type Role uint8

const (
	// IntermediateAgent 中间 AI Agent，处理的结果为中间产物
	IntermediateAgent Role = iota
	// FinalAgent 最终 AI Agent，输出最终结果
	FinalAgent
)

// HookType AI Agent 钩子触发节点类型
type HookType uint8

const (
	// BeforeProcess 在 AI Agent 处理前触发
	BeforeProcess HookType = iota
	// AfterProcess 在 AI Agent 处理后触发
	AfterProcess
)

const (
	AgentDataPath  = "agent.json"
	ScriptDataPath = "script.json"
)
