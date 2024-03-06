package workflow

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
)

type ConfigVersion struct {
	data.BaseConfigVersion
	Agents []string `json:"agents"`
	agents []core.AiAgent
}

// Config AI Workflow 初始化所需数据
type Config struct {
	data.BaseConfig[*ConfigVersion]
}

func (cv *ConfigVersion) Init() bool {
	ok := cv.BaseConfigVersion.Init()
	if !ok {
		return false
	}

	if cv.Agents == nil {
		cv.Agents = make([]string, 0)
	}
	cv.agents = make([]core.AiAgent, 0, len(cv.Agents))
	for _, agentId := range cv.Agents {
		agentItem, ok := data.GetOneItem(base.AgentDataName, agentId)
		if !ok {
			return false
		}
		agent, ok := agentItem.(core.AiAgent)
		if !ok {
			return false
		}
		cv.agents = append(cv.agents, agent)
	}

	return true
}
