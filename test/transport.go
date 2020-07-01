package test

import (
	"context"
	"encoding/json"
	"io"

	"github.com/mediabuyerbot/go-wirenet"
	wirenettransport "github.com/mediabuyerbot/go-wirenet-gokit"
)

func MakeWirenetHandlers(wire wirenet.Wire, endpoints Set) {
	options := make([]wirenettransport.ServerOption, 0)
	wire.Stream("updateBalance", wirenettransport.NewServer(
		endpoints.UpdateBalanceEndpoint,
		decodeWirenetUpdateBalanceRequest,
		encodeWirenetUpdateBalanceResponse,
		options...,
	).Handle)
}

func MakeWirenetClient(wire wirenet.Wire) Service {
	options := make([]wirenettransport.ClientOption, 0)
	return &Set{
		UpdateBalanceEndpoint: wirenettransport.NewClient(
			wire,
			"updateBalance",
			encodeWirenetUpdateBalanceRequest,
			decodeWirenetUpdateBalanceResponse,
			options...,
		).Endpoint(),
	}
}

func decodeWirenetUpdateBalanceRequest(_ context.Context, r io.ReadCloser) (request interface{}, err error) {
	defer r.Close()
	var req UpdateBalanceRequest
	err = json.NewDecoder(r).Decode(&req)
	return req, err
}

func encodeWirenetUpdateBalanceRequest(_ context.Context, request interface{}, w io.WriteCloser) error {
	defer w.Close()
	return json.NewEncoder(w).Encode(&request)
}

func decodeWirenetUpdateBalanceResponse(_ context.Context, r io.ReadCloser) (response interface{}, err error) {
	defer r.Close()
	var resp UpdateBalanceResponse
	err = json.NewDecoder(r).Decode(&resp)
	return resp, err
}

func encodeWirenetUpdateBalanceResponse(_ context.Context, response interface{}, w io.WriteCloser) error {
	defer w.Close()
	return json.NewEncoder(w).Encode(response)
}
