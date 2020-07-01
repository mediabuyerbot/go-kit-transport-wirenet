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

func MakeUploadFileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UploadFileRequest)
		err = s.UploadFile(ctx, req.Filepath, req.Size, req.Name)
		return UploadFileResponse{Err: err}, nil
	}
}

type Set struct {
	UpdateBalanceEndpoint endpoint.Endpoint
	UploadFileEndpoint    endpoint.Endpoint
}

func NewEndpointSet(s Service) Set {
	return Set{
		UpdateBalanceEndpoint: MakeUpdateBalanceEndpoint(s),
		UploadFileEndpoint:    MakeUploadFileEndpoint(s),
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

func (s Set) UploadFile(ctx context.Context, filepath string, size int64, name string) error {
	resp, err := s.UploadFileEndpoint(ctx, UploadFileRequest{
		Filepath: filepath,
		Size:     size,
		Name:     name,
	})
	if err != nil {
		return err
	}
	response := resp.(UploadFileResponse)
	return response.Err
}
