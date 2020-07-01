package wirenettransport

import (
	"context"
	"io"

	"github.com/mediabuyerbot/go-wirenet"
)

// EncodeRequestFunc encodes the passed request object into the wirenet request object.
type EncodeRequestFunc func(context.Context, interface{}, io.WriteCloser) error

// DecodeResponseFunc extracts a user-domain response object from a wirenet response object.
type DecodeResponseFunc func(context.Context, io.ReadCloser) (response interface{}, err error)

// DecodeRequestFunc extracts a user-domain request object from a wirenet request.
type DecodeRequestFunc func(context.Context, io.ReadCloser) (request interface{}, err error)

// EncodeResponseFunc encodes the passed response object to the wirenet response message.
type EncodeResponseFunc func(context.Context, interface{}, io.WriteCloser) error

// ClientCodec encodes and decodes the byte stream in the user-domain.
type ClientCodec func(context.Context, interface{}, wirenet.Stream) (interface{}, error)

// ServerCodec encodes and decodes the byte stream in the user-domain.
type ServerCodec func(context.Context, wirenet.Stream) (interface{}, error)
