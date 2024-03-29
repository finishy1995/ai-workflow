// Code generated by goctl. DO NOT EDIT.
// Source: ai.proto

package aiclient

import (
	"context"

	"github.com/finishy1995/effibot-core/ai/pb/ai"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ActiveVersionRequest  = ai.ActiveVersionRequest
	ActiveVersionResponse = ai.ActiveVersionResponse
	AddNodeRequest        = ai.AddNodeRequest
	AddNodeResponse       = ai.AddNodeResponse
	AddVersionRequest     = ai.AddVersionRequest
	AddVersionResponse    = ai.AddVersionResponse
	CreateSessionRequest  = ai.CreateSessionRequest
	CreateSessionResponse = ai.CreateSessionResponse
	DeleteSessionRequest  = ai.DeleteSessionRequest
	DeleteSessionResponse = ai.DeleteSessionResponse
	GetItemRequest        = ai.GetItemRequest
	GetItemResponse       = ai.GetItemResponse
	GetResultRequest      = ai.GetResultRequest
	GetResultResponse     = ai.GetResultResponse
	KeyValue              = ai.KeyValue
	Node                  = ai.Node

	Ai interface {
		GetResult(ctx context.Context, in *GetResultRequest, opts ...grpc.CallOption) (*GetResultResponse, error)
		CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*CreateSessionResponse, error)
		AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*AddNodeResponse, error)
		DeleteSession(ctx context.Context, in *DeleteSessionRequest, opts ...grpc.CallOption) (*DeleteSessionResponse, error)
		GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*GetItemResponse, error)
		AddVersion(ctx context.Context, in *AddVersionRequest, opts ...grpc.CallOption) (*AddVersionResponse, error)
		ActiveVersion(ctx context.Context, in *ActiveVersionRequest, opts ...grpc.CallOption) (*ActiveVersionResponse, error)
	}

	defaultAi struct {
		cli zrpc.Client
	}
)

func NewAi(cli zrpc.Client) Ai {
	return &defaultAi{
		cli: cli,
	}
}

func (m *defaultAi) GetResult(ctx context.Context, in *GetResultRequest, opts ...grpc.CallOption) (*GetResultResponse, error) {
	client := ai.NewAiClient(m.cli.Conn())
	return client.GetResult(ctx, in, opts...)
}

func (m *defaultAi) CreateSession(ctx context.Context, in *CreateSessionRequest, opts ...grpc.CallOption) (*CreateSessionResponse, error) {
	client := ai.NewAiClient(m.cli.Conn())
	return client.CreateSession(ctx, in, opts...)
}

func (m *defaultAi) AddNode(ctx context.Context, in *AddNodeRequest, opts ...grpc.CallOption) (*AddNodeResponse, error) {
	client := ai.NewAiClient(m.cli.Conn())
	return client.AddNode(ctx, in, opts...)
}

func (m *defaultAi) DeleteSession(ctx context.Context, in *DeleteSessionRequest, opts ...grpc.CallOption) (*DeleteSessionResponse, error) {
	client := ai.NewAiClient(m.cli.Conn())
	return client.DeleteSession(ctx, in, opts...)
}

func (m *defaultAi) GetItem(ctx context.Context, in *GetItemRequest, opts ...grpc.CallOption) (*GetItemResponse, error) {
	client := ai.NewAiClient(m.cli.Conn())
	return client.GetItem(ctx, in, opts...)
}

func (m *defaultAi) AddVersion(ctx context.Context, in *AddVersionRequest, opts ...grpc.CallOption) (*AddVersionResponse, error) {
	client := ai.NewAiClient(m.cli.Conn())
	return client.AddVersion(ctx, in, opts...)
}

func (m *defaultAi) ActiveVersion(ctx context.Context, in *ActiveVersionRequest, opts ...grpc.CallOption) (*ActiveVersionResponse, error) {
	client := ai.NewAiClient(m.cli.Conn())
	return client.ActiveVersion(ctx, in, opts...)
}
