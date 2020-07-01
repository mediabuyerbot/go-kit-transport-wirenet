package wirenettransport

import (
	"context"
	"io"
)

type EncodeRequestFunc func(context.Context, interface{}, io.WriteCloser) error
type DecodeResponseFunc func(context.Context, io.ReadCloser) (response interface{}, err error)
type DecodeRequestFunc func(context.Context, io.ReadCloser) (request interface{}, err error)
type EncodeResponseFunc func(context.Context, interface{}, io.WriteCloser) error
