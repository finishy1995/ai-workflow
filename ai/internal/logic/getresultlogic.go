package logic

import (
	"context"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/finishy1995/effibot-core/ai/internal/agent/core"
	"github.com/finishy1995/effibot-core/ai/internal/svc"
	"github.com/finishy1995/effibot-core/ai/internal/workflow"
	"github.com/finishy1995/effibot-core/ai/pb/ai"
	"github.com/finishy1995/effibot-core/base"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResultLogic {
	return &GetResultLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetResultLogic) GetResult(in *ai.GetResultRequest) (resp *ai.GetResultResponse, errResp error) {
	resp = &ai.GetResultResponse{}

	if strutil.IsBlank(in.GetSessionId()) || strutil.IsBlank(in.GetNodeId()) {
		resp.Code = base.CodeInvalidParam
		return
	}

	var node core.Task
	if err := l.svcCtx.DB.First(&node, "", in.GetNodeId()); err != nil {
		logx.Errorf("l.svcCtx.DB.First error: %s, sessionId: %s, nodeId: %s", err.Error(), in.GetSessionId(), in.GetNodeId())
		resp.Code = base.CodeDBError
		return
	}
	result := ""
	if node.Output != nil && len(node.Output) > 0 {
		result = node.Output[0]
	}

	resp.Code = base.CodeOk
	resp.Node = &ai.Node{
		Id:     node.TaskId,
		Type:   workflow.GetType(node.WorkflowInstanceId),
		Status: ai.NodeStatus(node.Status),
		Params: core.Map2KeyValue(node.Param),
		Result: result,
	}
	return
}
