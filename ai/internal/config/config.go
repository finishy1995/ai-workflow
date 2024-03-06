package config

import (
	"github.com/finishy1995/effibot-core/ai/internal/gpt"
	"github.com/finishy1995/go-library/storage"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Spec AiConfig
}

type AiConfig struct {
	GPT     gpt.Config
	Storage storage.Config
}
