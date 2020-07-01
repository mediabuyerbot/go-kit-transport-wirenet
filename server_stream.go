package wirenettransport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/mediabuyerbot/go-wirenet"
)

// StreamServer wraps an endpoint and implements wirenet.Handler.
type StreamServer struct {
	e            endpoint.Endpoint
	codec        ServerCodec
	errorHandler transport.ErrorHandler
	before       []ServerRequestFunc
}

// NewStreamServer constructs a new server, which implements wraps the provided
// endpoint and implements the Handler interface.
func NewStreamServer(
	e endpoint.Endpoint,
	codec ServerCodec,
	options ...StreamServerOption,
) *StreamServer {
	s := &StreamServer{
		e:            e,
		codec:        codec,
		errorHandler: transport.NewLogErrorHandler(log.NewNopLogger()),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// StreamServerOption sets an optional parameter for servers.
type StreamServerOption func(*StreamServer)

// StreamServerErrorHandler is used to handle non-terminal errors. By default,
//non-terminal errors are ignored. This is intended as a diagnostic measure.
func StreamServerErrorHandler(errorHandler transport.ErrorHandler) StreamServerOption {
	return func(s *StreamServer) { s.errorHandler = errorHandler }
}

// StreamServerBefore functions are executed on the stream request object before the
// request is decoded.
func StreamServerBefore(before ...ServerRequestFunc) StreamServerOption {
	return func(s *StreamServer) { s.before = append(s.before, before...) }
}

// Handle implements the Handler interface.
func (s StreamServer) Handle(ctx context.Context, stream wirenet.Stream) {
	defer stream.Close()

	for _, f := range s.before {
		ctx = f(ctx)
	}

	request, err := s.codec(ctx, stream)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		return
	}

	_, err = s.e(ctx, request)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		return
	}
}
