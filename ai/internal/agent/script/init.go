package script

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/chat"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/common"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/lua"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/style"
	"github.com/finishy1995/effibot-core/ai/internal/agent/script/transformer"
	"github.com/finishy1995/go-library/storage"
)

func Init(st storage.Storage) {
	common.Init(st)
	chat.Init(st)
	lua.Init(st)
	style.Init(st)
	transformer.Init(st)
}
