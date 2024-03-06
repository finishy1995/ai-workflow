package chat

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/common"
	"github.com/finishy1995/effibot-core/library/data"
	"github.com/finishy1995/go-library/storage"
)

type ConfigVersion struct {
	common.ConfigVersion
	EnvConfig
}

type Config struct {
	data.BaseConfig[*ConfigVersion]
}

const (
	ScriptType = "chat"
)

func Init(st storage.Storage) {
	core.ScriptInit(ScriptType, NewScript, &Config{})
}
