package core

import (
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"
)

type AiScript interface {
	data.Item
	Run(task *Task, param map[string]interface{}) error
}

func GetScriptExpr(typ string) string {
	return "BaseConfig.Type == " + typ
}

func ScriptInit(typ string, f data.InitFunc, config data.Config) {
	config.SetType(typ)
	err := data.SetupPool(base.ScriptDataName, f, config, &data.LoaderConfig{
		Type:      data.StorageWithFileInit,
		FilePath:  base.DataPath + ScriptDataPath,
		Expr:      GetScriptExpr(typ),
		TableName: base.ScriptTableName,
	})
	if err != nil {
		panic(err)
	}
}
