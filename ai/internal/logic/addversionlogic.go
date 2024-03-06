package logic

import (
	"context"
	"github.com/finishy1995/effibot-core/base"
	"github.com/finishy1995/effibot-core/library/data"

	"github.com/finishy1995/effibot-core/ai/internal/svc"
	"github.com/finishy1995/effibot-core/ai/pb/ai"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddVersionLogic {
	return &AddVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddVersionLogic) AddVersion(in *ai.AddVersionRequest) (*ai.AddVersionResponse, error) {
	resp := &ai.AddVersionResponse{
		Code: base.CodeOk,
	}
	item, ok := data.GetOneItem(in.GetType(), in.GetId())
	if !ok {
		resp.Code = base.CodeNotFound
		return resp, nil
	}

	version, err := item.GetConfig().NewVersion([]byte(in.GetChange()))
	if err != nil {
		l.Errorf("AddVersion failed by error: %s", err.Error())
		resp.Code = base.CodeDBError
	} else {
		item.GetConfig().ActiveVersion(version)
		resp.Version = uint32(version)
	}

	return resp, nil
}
