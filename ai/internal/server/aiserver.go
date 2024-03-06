// Code generated by goctl. DO NOT EDIT.
// Source: ai.proto

package server

import (
	"context"

	"github.com/finishy1995/effibot-core/ai/internal/logic"
	"github.com/finishy1995/effibot-core/ai/internal/svc"
	"github.com/finishy1995/effibot-core/ai/pb/ai"
)

type AiServer struct {
	svcCtx *svc.ServiceContext
	ai.UnimplementedAiServer
}

func NewAiServer(svcCtx *svc.ServiceContext) *AiServer {
	return &AiServer{
		svcCtx: svcCtx,
	}
}

func (s *AiServer) GetResult(ctx context.Context, in *ai.GetResultRequest) (*ai.GetResultResponse, error) {
	l := logic.NewGetResultLogic(ctx, s.svcCtx)
	return l.GetResult(in)
}

func (s *AiServer) CreateSession(ctx context.Context, in *ai.CreateSessionRequest) (*ai.CreateSessionResponse, error) {
	l := logic.NewCreateSessionLogic(ctx, s.svcCtx)
	return l.CreateSession(in)
}

func (s *AiServer) AddNode(ctx context.Context, in *ai.AddNodeRequest) (*ai.AddNodeResponse, error) {
	l := logic.NewAddNodeLogic(ctx, s.svcCtx)
	return l.AddNode(in)
}

func (s *AiServer) DeleteSession(ctx context.Context, in *ai.DeleteSessionRequest) (*ai.DeleteSessionResponse, error) {
	l := logic.NewDeleteSessionLogic(ctx, s.svcCtx)
	return l.DeleteSession(in)
}

func (s *AiServer) GetItem(ctx context.Context, in *ai.GetItemRequest) (*ai.GetItemResponse, error) {
	l := logic.NewGetItemLogic(ctx, s.svcCtx)
	return l.GetItem(in)
}

func (s *AiServer) AddVersion(ctx context.Context, in *ai.AddVersionRequest) (*ai.AddVersionResponse, error) {
	l := logic.NewAddVersionLogic(ctx, s.svcCtx)
	return l.AddVersion(in)
}

func (s *AiServer) ActiveVersion(ctx context.Context, in *ai.ActiveVersionRequest) (*ai.ActiveVersionResponse, error) {
	l := logic.NewActiveVersionLogic(ctx, s.svcCtx)
	return l.ActiveVersion(in)
}
