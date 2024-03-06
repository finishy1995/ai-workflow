package agent

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
	"github.com/finishy1995/go-library/storage"
)

func Init(st storage.Storage) {
	data.SetStorage(st)

	script.Init(st)
	err := st.CreateTable(core.Task{}, "")
	err = st.CreateTable(core.AgentConfig{}, base.AgentTableName)
	if err != nil {
		panic(err)
	}
	err = data.SetupPool(base.AgentDataName, core.NewBaseAiAgent, &core.AgentConfig{}, &data.LoaderConfig{
		Type:      data.StorageWithFileInit,
		FilePath:  base.DataPath + core.AgentDataPath,
		TableName: base.AgentTableName,
	})
	if err != nil {
		panic(err)
	}
}
