package svc

import (
	"github.com/finishy1995/effibot-core/ai/internal/agent"
	"github.com/finishy1995/effibot-core/ai/internal/config"
	"github.com/finishy1995/effibot-core/ai/internal/gpt"
	"github.com/finishy1995/effibot-core/ai/internal/workflow"
	"github.com/finishy1995/go-library/storage"
)

type ServiceContext struct {
	Config config.Config
	DB     storage.Storage
}

func NewServiceContext(c config.Config) *ServiceContext {
	st := storage.NewStorage(&c.Spec.Storage)
	gpt.Setup(c.Spec.GPT)
	agent.Init(st)
	workflow.Init(st)

	return &ServiceContext{
		Config: c,
		DB:     st,
	}
}
