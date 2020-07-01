package test

import (
	"context"
)

type Service interface {
	UpdateBalance(ctx context.Context, a int, b int) (sum int, err error)
	UploadFile(ctx context.Context, filepath string, size int64, name string) error
}

type service struct{}

func (service) UpdateBalance(ctx context.Context, a int, b int) (sum int, err error) {
	// tx begin :)
	sum = a + b
	// tx commit
	return
}

func (service) UploadFile(ctx context.Context, filepath string, size int64, name string) error {
	return nil
}

func NewService() Service {
	return service{}
}
