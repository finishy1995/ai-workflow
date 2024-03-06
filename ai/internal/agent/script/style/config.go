package style

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/common"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
	"github.com/finishy1995/go-library/storage"
)

type ConfigVersion struct {
	common.ConfigVersion
	EnvConfig
	StyleType    string `json:"style_type"`
	styleTypeMap map[string]*Type
}

type Config struct {
	data.BaseConfig[*ConfigVersion]
}

func (c *Config) Init() bool {
	c.SetEnv(defaultEnvSettings)
	return c.BaseConfig.Init()
}

func (cv *ConfigVersion) Init() bool {
	cv.BaseConfigVersion.Init()

	types := data.GetAllTypeItems(base.TypeDataName, cv.StyleType)
	cv.styleTypeMap = make(map[string]*Type, len(types))
	for _, typ := range types {
		typReal, ok := typ.(*Type)
		if !ok {
			return false
		}
		cv.styleTypeMap[typReal.GetName()] = typReal
	}

	if cv.PromptPrefix == nil {
		cv.PromptPrefix = make([]string, 0)
	}
	if cv.PromptSuffix == nil {
		cv.PromptSuffix = make([]string, 0)
	}
	return true
}

const (
	ScriptType = "style"
)

func Init(st storage.Storage) {
	err := st.CreateTable(TypeConfig{}, base.StyleTableName)
	if err != nil {
		panic(err)
	}
	err = data.SetupPool(base.TypeDataName, NewType, &TypeConfig{}, &data.LoaderConfig{
		Type:      data.StorageWithFileInit,
		FilePath:  base.DataPath + TypeDataPath,
		TableName: base.StyleTableName,
	})
	if err != nil {
		panic(err)
	}

	core.ScriptInit(ScriptType, NewScript, &Config{})
}
