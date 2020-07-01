package wirenettransport

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/mediabuyerbot/go-wirenet"
)

var (
	// ErrSessionIDNotAssigned returned when no session is defined in context.
	ErrSessionIDNotDefined = errors.New("session id not defined")
)

// Client wraps a wirenet connection and provides a method
// that implements endpoint.Endpoint.
type Client struct {
	streamName string
	wire       wirenet.Wire
	enc        EncodeRequestFunc
	dec        DecodeResponseFunc
	before     []ClientRequestFunc
}

// NewClient constructs a usable Client for a single remote endpoint.
func NewClient(
	wire wirenet.Wire,
	streamName string,
	enc EncodeRequestFunc,
	dec DecodeResponseFunc,
	options ...ClientOption,
) *Client {
	c := &Client{
		streamName: streamName,
		wire:       wire,
		enc:        enc,
		dec:        dec,
		before:     []ClientRequestFunc{},
	}
	for _, option := range options {
		option(c)
	}
	return c
}

// ClientOption sets an optional parameter for clients.
type ClientOption func(*Client)

// ClientBefore sets the RequestFuncs that are applied to the outgoing request
// before it's invoked.
func ClientBefore(before ...ClientRequestFunc) ClientOption {
	return func(c *Client) { c.before = append(c.before, before...) }
}

// Endpoint returns a usable endpoint that will invoke the wirenet specified
// by the client.
func (c Client) Endpoint() endpoint.Endpoint {
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

		reader := stream.Reader()
		writer := stream.Writer()

		if err := c.enc(ctx, request, writer); err != nil {
			return nil, err
		}
		return c.dec(ctx, reader)
	}
}
