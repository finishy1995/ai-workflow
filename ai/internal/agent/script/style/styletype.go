package style

import (
	"github.com/finishy1995/effibot-core/library/data"
)

type TypeConfigVersion struct {
	data.BaseConfigVersion
	Prompt   []string `json:"prompt"`
	MaxToken int      `json:"max_token"`
}

type TypeConfig struct {
	data.BaseConfig[*TypeConfigVersion]
}

type Type struct {
	data.BaseItem
}

func NewType() data.Item {
	return new(Type)
}

const (
	TypeDataPath = "style.json"
)
