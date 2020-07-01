package test

import (
	"context"
)

type Service interface {
	UpdateBalance(ctx context.Context, a int, b int) (sum int, err error)
}

type service struct{}

func (service) UpdateBalance(ctx context.Context, a int, b int) (sum int, err error) {
	// tx begin :)
	sum = a + b
	// tx commit
	return
}

func NewService() Service {
	return service{}
}
