package wirenettransport

import (
	"context"

	"github.com/google/uuid"
)

// ContextKeySessionID used to inject into session context.
type ContextKeySessionID struct{}

// ClientRequestFunc can take information from the context and use it to build the stream.
type ClientRequestFunc func(context.Context) context.Context

// ServerRequestFunc can take information from the context and use it to the handle stream.
type ServerRequestFunc func(context.Context) context.Context

// SetSessionID returns a ClientRequestFunc that sets the session id.
func SetSessionID(sid uuid.UUID) ClientRequestFunc {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ContextKeySessionID{}, sid)
	}
}

// InjectSessionID returns a new context with session id.
func InjectSessionID(sid uuid.UUID, ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeySessionID{}, sid)
}
