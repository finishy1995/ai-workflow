package core

import (
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
)

type AgentConfigVersion struct {
	data.BaseConfigVersion
	Script string                 `json:"script"`
	Env    map[string]interface{} `json:"env"`
	script AiScript
}

// AgentConfig AI Agent 初始化所需数据
type AgentConfig struct {
	data.BaseConfig[*AgentConfigVersion]
}

func (acv *AgentConfigVersion) Init() bool {
	ok := acv.BaseConfigVersion.Init()
	if !ok {
		return false
	}
	if acv.Env == nil {
		acv.Env = make(map[string]interface{})
	}

	initScriptItem, ok := data.GetOneItem(base.ScriptDataName, acv.Script)
	if !ok {
		return false
	}
	acv.script, ok = initScriptItem.(AiScript)
	if !ok {
		return false
	}

	return true
}
