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

type AddNodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddNodeLogic {
	return &AddNodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddNodeLogic) AddNode(in *ai.AddNodeRequest) (resp *ai.AddNodeResponse, errResp error) {
	resp = &ai.AddNodeResponse{}

	params := core.KeyValue2Map(in.Params, in.Settings)
	if params == nil {
		resp.Code = base.CodeInvalidParam
		return
	}

	_, nodeId, err := workflow.ProcessOld(in.Type, in.SessionId, "", params)
	if err != nil {
		logx.Errorf("aiworkflow.Process error: %s", err.Error())
		resp.Code = base.CodeAiAgentProcessError
		return
	}

	resp.Code = base.CodeOk
	resp.NodeId = nodeId
	return
}
