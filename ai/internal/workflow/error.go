package workflow

import "errors"

var (
	ErrInvalidSessionId  = errors.New("invalid sessionId format")
	ErrInvalidWorkflowId = errors.New("node type unsupported")
)
