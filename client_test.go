package wirenettransport_test

import (
	"context"
	"io"
	"sync"
	"testing"
	"time"

	wirenettransport "github.com/mediabuyerbot/go-wirenet-gokit"

	"github.com/stretchr/testify/assert"

	"github.com/mediabuyerbot/go-wirenet"
	test "github.com/mediabuyerbot/go-wirenet-gokit/test"
)

// server <-> {client1, client2}
func TestClient_Endpoint(t *testing.T) {
	addr := ":8989"
	initServer := make(chan struct{})
	wait := func() {
		time.Sleep(5 * time.Second)
	}
	// server side
	server, err := wirenet.Mount(addr,
		wirenet.WithConnectHook(func(closer io.Closer) {
			close(initServer)
		}),
		wirenet.WithSessionOpenHook(func(session wirenet.Session) {
			t.Logf("open session id %s", session.ID())
		}),
		wirenet.WithSessionCloseHook(func(session wirenet.Session) {
			t.Logf("close session id %s", session.ID())
		}),
	)
	assert.Nil(t, err)
	go func() {
		assert.Nil(t, server.Connect())
	}()
	<-initServer

	// request from server to {client1, client2}
	go func() {
		client := test.MakeWirenetClient(server)
		ctx := context.Background()

		for {
			time.Sleep(500 * time.Millisecond)
			for uuid, sess := range server.Sessions() {
				t.Logf("request to %s, uuid %s", sess.Identification(), sess.ID())
				ctxWithCurrentSession := wirenettransport.InjectSessionID(uuid, ctx)
				sum, err := client.UpdateBalance(ctxWithCurrentSession, 1, 4)
				assert.Nil(t, err)
				assert.Equal(t, 5, sum)
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	// client1
	go func() {
		// transport
		cid := wirenet.Identification("client1")
		client1, err := wirenet.Join(addr, wirenet.WithIdentification(cid, nil))
		assert.Nil(t, err)
		// go-kit
		svc := test.NewService()
		endpoints := test.NewEndpointSet(svc)
		test.MakeWirenetHandlers(client1, endpoints)
		go func() {
			defer wg.Done()
			wait()
			assert.Nil(t, client1.Close())
		}()
		assert.Nil(t, client1.Connect())
	}()

	// client2
	go func() {
		// transport
		cid := wirenet.Identification("client2")
		client2, err := wirenet.Join(addr, wirenet.WithIdentification(cid, nil))
		assert.Nil(t, err)
		// go-kit
		svc := test.NewService()
		endpoints := test.NewEndpointSet(svc)
		test.MakeWirenetHandlers(client2, endpoints)
		go func() {
			defer wg.Done()
			wait()
			assert.Nil(t, client2.Close())
		}()
		assert.Nil(t, client2.Connect())
	}()

	wg.Wait()
	assert.Nil(t, server.Close())
	assert.Len(t, server.Sessions(), 0)
}
