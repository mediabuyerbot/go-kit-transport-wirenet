package test

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeUpdateBalanceEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateBalanceRequest)
		sum, err := s.UpdateBalance(ctx, req.A, req.B)
		return UpdateBalanceResponse{Sum: sum, Err: err}, nil
	}
}

type Set struct {
	UpdateBalanceEndpoint endpoint.Endpoint
}

func NewEndpointSet(s Service) Set {
	return Set{
		UpdateBalanceEndpoint: MakeUpdateBalanceEndpoint(s),
	}
}

func (s Set) UpdateBalance(ctx context.Context, a, b int) (int, error) {
	resp, err := s.UpdateBalanceEndpoint(ctx, UpdateBalanceRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(UpdateBalanceResponse)
	return response.Sum, response.Err
}
