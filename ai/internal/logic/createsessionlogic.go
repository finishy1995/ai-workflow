package logic

import (
	"context"
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/ai/internal/svc"
	"github.com/finishy1995/effibot-core/ai/internal/workflow"
	"github.com/finishy1995/effibot-core/ai/pb/ai"
	"github.com/finishy1995/effibot-core/base"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSessionLogic) CreateSession(in *ai.CreateSessionRequest) (resp *ai.CreateSessionResponse, errResp error) {
	resp = &ai.CreateSessionResponse{}

	params := core.KeyValue2Map(in.Params, in.Settings)
	if params == nil {
		resp.Code = base.CodeInvalidParam
		return
	}
	sessionId, nodeId, err := workflow.ProcessOld(in.Type, "", in.UserId, params)
	if err != nil {
		logx.Errorf("aiworkflow.Process error: %s", err.Error())
		resp.Code = base.CodeAiAgentProcessError
		return
	}

	resp.Code = base.CodeOk
	resp.SessionId = sessionId
	resp.RootNodeId = nodeId
	return
}
