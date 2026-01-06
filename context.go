package socketiominiservice

import (
	"context"

	"github.com/google/uuid"
)

// RequestIdKey is the key for storing request ID in context
const RequestIdKey string = "requestId"

// ClientContextKey is the key for storing client information in context
const ClientContextKey string = "clientInfo"

// SetRequestIdToOrigCtx sets a new request ID to the context
func SetRequestIdToOrigCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestIdKey, uuid.NewString())
}

// GetRequestIdFromCtx retrieves the request ID from context
func GetRequestIdFromCtx(ctx context.Context) string {
	var requestId string
	var ok bool
	newReqId := func() string {
		return "##" + uuid.NewString()
	}
	if requestId, ok = ctx.Value(RequestIdKey).(string); !ok || requestId == "" {
		requestId = newReqId()
	}
	return requestId
}

// SetClientToCtx sets client information to the context
func SetClientToCtx[T SocketIOConstraint](ctx context.Context, clientInfo ClientInterface[T]) context.Context {
	return context.WithValue(ctx, ClientContextKey, clientInfo)
}

// GetClientFromContext retrieves client information from the context
func GetClientFromContext[T SocketIOConstraint](ctx context.Context) ClientInterface[T] {
	if clientInfo, ok := ctx.Value(ClientContextKey).(ClientInterface[T]); ok {
		return clientInfo
	} else {
		return nil
	}
}
