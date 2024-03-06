package logic

import (
	"context"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/finishy1995/effibot-core/ai/internal/svc"
	"github.com/finishy1995/effibot-core/ai/pb/ai"
	"github.com/finishy1995/effibot-core/base"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSessionLogic) DeleteSession(in *ai.DeleteSessionRequest) (*ai.DeleteSessionResponse, error) {

	if strutil.IsBlank(in.GetSessionId()) {
		return &ai.DeleteSessionResponse{Code: base.CodeInvalidParam}, nil
	}

	// todo: delete

	return &ai.DeleteSessionResponse{Code: base.CodeOk}, nil
}
