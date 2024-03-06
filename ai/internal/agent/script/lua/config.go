package lua

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/common"
	"github.com/finishy1995/effibot-core/library/data"
	"github.com/finishy1995/go-library/storage"
)

type ConfigVersion struct {
	common.ConfigVersion
	Script string `json:"script"`
}

type Config struct {
	data.BaseConfig[*ConfigVersion]
}

const (
	ScriptType = "lua"
)

func Init(st storage.Storage) {
	core.ScriptInit(ScriptType, NewScript, &Config{})
}
