package logic

import (
	"context"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"

	"github.com/finishy1995/effibot-core/ai/internal/svc"
	"github.com/finishy1995/effibot-core/ai/pb/ai"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActiveVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActiveVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActiveVersionLogic {
	return &ActiveVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ActiveVersionLogic) ActiveVersion(in *ai.ActiveVersionRequest) (*ai.ActiveVersionResponse, error) {
	resp := &ai.ActiveVersionResponse{
		Code: base.CodeOk,
	}
	item, ok := data.GetOneItem(in.GetType(), in.GetId())
	if !ok {
		resp.Code = base.CodeNotFound
		return resp, nil
	}

	ok = item.GetConfig().ActiveVersion(data.VersionType(in.GetVersion()))
	if !ok {
		resp.Code = base.CodeDBError
	}

	return resp, nil
}
