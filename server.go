package wirenettransport

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/mediabuyerbot/go-wirenet"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport"
)

// Handler which should be called from the wirenet binding of the service
// implementation. The incoming request parameter, and returned response
// parameter, are both wirenet types, not user-domain.
type Handler interface {
	Serve(ctx context.Context, request interface{}) (context.Context, interface{}, error)
}

// Server wraps an endpoint and implements wirenet.Handler.
type Server struct {
	e            endpoint.Endpoint
	dec          DecodeRequestFunc
	enc          EncodeResponseFunc
	errorHandler transport.ErrorHandler
	before       []ServerRequestFunc
}

func NewServer(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	options ...ServerOption,
) *Server {
	s := &Server{
		e:            e,
		dec:          dec,
		enc:          enc,
		errorHandler: transport.NewLogErrorHandler(log.NewNopLogger()),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// ServerOption sets an optional parameter for servers.
type ServerOption func(*Server)

// ServerErrorHandler is used to handle non-terminal errors. By default,
//non-terminal errors are ignored. This is intended as a diagnostic measure.
func ServerErrorHandler(errorHandler transport.ErrorHandler) ServerOption {
	return func(s *Server) { s.errorHandler = errorHandler }
}

// ServerBefore functions are executed on the stream request object before the
// request is decoded.
func ServerBefore(before ...ServerRequestFunc) ServerOption {
	return func(s *Server) { s.before = append(s.before, before...) }
}

// Handle implements the Handler interface.
func (s Server) Handle(ctx context.Context, stream wirenet.Stream) {
	defer stream.Close()

	reader := stream.Reader()
	writer := stream.Writer()

	for _, f := range s.before {
		ctx = f(ctx)
	}

	request, err := s.dec(ctx, reader)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		return
	}

	response, err := s.e(ctx, request)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		return
	}

	if err := s.enc(ctx, response, writer); err != nil {
		s.errorHandler.Handle(ctx, err)
	}
}
