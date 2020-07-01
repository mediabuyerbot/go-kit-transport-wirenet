package wirenettransport

import (
	"context"

	"github.com/google/uuid"

	"github.com/go-kit/kit/endpoint"
	"github.com/mediabuyerbot/go-wirenet"
)

// StreamClient wraps a wirenet connection and provides a method
// that implements endpoint.Endpoint.
type StreamClient struct {
	streamName string
	wire       wirenet.Wire
	before     []ClientRequestFunc
	codec      ClientCodec
}

// NewStreamClient constructs a usable StreamClient for a single remote endpoint.
func NewStreamClient(
	wire wirenet.Wire,
	streamName string,
	codec ClientCodec,
	options ...StreamClientOption,
) *StreamClient {
	c := &StreamClient{
		streamName: streamName,
		wire:       wire,
		codec:      codec,
		before:     []ClientRequestFunc{},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// StreamClientOption sets an optional parameter for clients.
type StreamClientOption func(*StreamClient)

// StreamClientBefore sets the RequestFuncs that are applied to the outgoing request
// before it's invoked.
func StreamClientBefore(before ...ClientRequestFunc) StreamClientOption {
	return func(c *StreamClient) { c.before = append(c.before, before...) }
}

// Endpoint returns a usable endpoint that will invoke the wirenet specified
// by the client.
func (c StreamClient) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		for _, f := range c.before {
			ctx = f(ctx)
		}

		sessionID, ok := ctx.Value(ContextKeySessionID{}).(uuid.UUID)
		if !ok {
			return nil, ErrSessionIDNotDefined
		}

		session, err := c.wire.Session(sessionID)
		if err != nil {
			return nil, err
		}

		stream, err := session.OpenStream(c.streamName)
		if err != nil {
			return nil, err
		}
		defer stream.Close()

		return c.codec(ctx, request, stream)
	}
}
