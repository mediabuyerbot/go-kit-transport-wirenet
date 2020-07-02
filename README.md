# go-wirenet-gokit
Transport layer for go-kit

See examples in tests files at the links:

[StreamEndpoint](https://github.com/mediabuyerbot/go-wirenet-gokit/blob/master/client_stream_test.go),
[Endpoint](https://github.com/mediabuyerbot/go-wirenet-gokit/blob/master/client_test.go),
[TestService](https://github.com/mediabuyerbot/go-wirenet-gokit/tree/master/test)

## Examples (Key points)
+ [StreamEndpoint](#stream-endpoint)
+ [Endpoint](#endpoint)

### Stream Endpoint
transport/wirenet.go
```go
// server side
func MakeWirenetHandlers(wire wirenet.Wire, endpoints Set) {
	wire.Stream("uploadFile", wirenettransport.NewStreamServer(
		endpoints.UploadFileEndpoint,
		uploadFileServerSideCodec,
		[]wirenettransport.StreamServerOption{}...,
	).Handle)
}
```

transport/wirenet.go
```go
// client side
func MakeWirenetClient(wire wirenet.Wire) Service {
	return &Set{
		UploadFileEndpoint: wirenettransport.NewStreamClient(
			wire,
			"uploadFile",
			uploadFileClientSideCodec,
			[]wirenettransport.StreamClientOption{}...,
		).Endpoint(),
	}
}
```

transport/wirenet_codec.go
```go
func uploadFileServerSideCodec(_ context.Context, s wirenet.Stream) (interface{}, error) {
	defer s.Close()

	w := s.Writer()
	r := s.Reader()

	// read fileInfo
	var req UploadFileRequest
	if err := json.NewDecoder(r).Decode(&req); err != nil {
		return nil, err
	}
	r.Close()

	// read data
	fp := filepath.Join(os.TempDir(), req.Name)
	file, err := os.Create(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	_, err = s.WriteTo(file)

	resp := &UploadFileResponse{
		Err: err,
	}

	// write data
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return nil, err
	}
	return resp, nil
}
```

transport/wirenet_codec.go
```go
func uploadFileClientSideCodec(_ context.Context, request interface{}, s wirenet.Stream) (interface{}, error) {
	req := request.(UploadFileRequest)
	file, err := os.Open(req.Filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	defer s.Close()

	w := s.Writer()
	r := s.Reader()

	// write fileInfo
	if err := json.NewEncoder(w).Encode(&req); err != nil {
		return nil, err
	}
	w.Close()

	// write data
	if _, err = s.ReadFrom(file); err != nil {
		return nil, err
	}

	// read data
	var resp UploadFileResponse
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return err, nil
	}
	r.Close()

	return resp, nil
}
```

cmd/client.go
```go
client := MakeWirenetClient(wire)
...
sess := sessionHub.findSessionByUser("123")
...
// Sets the current wirenet session id
ctxWithSess := wirenettransport.InjectSessionID(sess.ID(), ctx)
err = client.UploadFile(ctxWithSess, "./test/testdata/data.db", 1024, "data.db")
```

cmd/server.go
```go
svc := NewService()
endpoints := NewEndpointSet(svc)
MakeWirenetHandlers(wire, endpoints)
```

## Endpoint
transport/wirenet.go
```go
// server side
func MakeWirenetHandlers(wire wirenet.Wire, endpoints Set) {
	options := make([]wirenettransport.ServerOption, 0)
	wire.Stream("updateBalance", wirenettransport.NewServer(
		endpoints.UpdateBalanceEndpoint,
		decodeWirenetUpdateBalanceRequest,
		encodeWirenetUpdateBalanceResponse,
		options...,
	).Handle)
}
```

transport/wirenet.go
```go
// client side
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
```

transport/wirenet_encode_decode.go
```go
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
```

cmd/server.go
```go
// go-kit
svc := NewService()
endpoints := NewEndpointSet(svc)
MakeWirenetHandlers(wire, endpoints)
```

cmd/client.go
```go
client := MakeWirenetClient(wire)
...
sess := sessionHub.findSessionByUser("123")
...
// Sets the current wirenet session id
ctxWithSess := wirenettransport.InjectSessionID(sess.ID(), ctx)
sum, err := client.UpdateBalance(ctxWithSess, 1, 4)
```