package wirenettransport_test

import (
	"context"
	"io"
	"sync"
	"testing"

	wirenettransport "github.com/mediabuyerbot/go-kit-transport-wirenet"
	"github.com/mediabuyerbot/go-kit-transport-wirenet/test"
	"github.com/mediabuyerbot/go-wirenet"
	"github.com/stretchr/testify/assert"
)

func TestStreamServer_Endpoint(t *testing.T) {
	addr := randomAddr(t)
	initServer := make(chan struct{})

	// server side
	server, err := wirenet.Mount(addr,
		wirenet.WithConnectHook(func(closer io.Closer) {
			close(initServer)
		}),
	)
	assert.Nil(t, err)
	// go-kit
	svc := test.NewService()
	endpoints := test.NewEndpointSet(svc)
	test.MakeWirenetHandlers(server, endpoints)
	go func() {
		assert.Nil(t, server.Connect())
	}()
	<-initServer

	var wg sync.WaitGroup
	wg.Add(1)

	// client1
	go func() {
		cid := wirenet.Identification("client1")
		sessCh := make(chan wirenet.Session)
		wire, err := wirenet.Join(addr,
			wirenet.WithSessionOpenHook(func(session wirenet.Session) {
				sessCh <- session
			}),
			wirenet.WithIdentification(cid, nil))
		assert.Nil(t, err)
		go func() {
			assert.Nil(t, wire.Connect())
		}()

		sess := <-sessCh
		client := test.MakeWirenetClient(wire)
		ctx := wirenettransport.InjectSessionID(sess.ID(), context.Background())
		err = client.UploadFile(ctx, "./test/testdata/data.db", 1024, "data.db")
		assert.Nil(t, err)
		wg.Done()
	}()

	wg.Wait()

	assert.Nil(t, server.Close())
}
