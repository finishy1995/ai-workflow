package logic

import (
	"context"
	"encoding/json"
	"github.com/finishy1995/effibot-core/ai/internal/svc"
	"github.com/finishy1995/effibot-core/ai/pb/ai"
	"github.com/finishy1995/effibot-core/library/data"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetItemLogic {
	return &GetItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetItemLogic) GetItem(in *ai.GetItemRequest) (*ai.GetItemResponse, error) {
	resp := &ai.GetItemResponse{}
	result := data.GetAllItems(in.GetType())
	if result == nil || len(result) == 0 {
		return resp, nil
	}
	handled := make([]data.Config, 0, len(result))
	for _, value := range result {
		handled = append(handled, value.GetConfig())
	}

	b, err := json.Marshal(handled)
	if err != nil {
		l.Errorf("cannot marshal data, error: %s", err.Error())
	}
	resp.Result = string(b)

	return resp, nil
}
