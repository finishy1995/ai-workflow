package common

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
	"github.com/finishy1995/go-library/storage"
)

type ConfigVersion struct {
	data.BaseConfigVersion
}

type Config struct {
	data.BaseConfig[*ConfigVersion]
}

const (
	ScriptType = "common"
)

func Init(st storage.Storage) {
	err := st.CreateTable(Config{}, base.ScriptTableName)
	if err != nil {
		panic(err)
	}
	core.ScriptInit(ScriptType, NewScript, &Config{})
}
