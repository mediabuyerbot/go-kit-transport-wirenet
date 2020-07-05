package wirenettransport_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	wirenettransport "github.com/mediabuyerbot/go-kit-transport-wirenet"

	"github.com/stretchr/testify/assert"

	test "github.com/mediabuyerbot/go-kit-transport-wirenet/test"
	"github.com/mediabuyerbot/go-wirenet"
)

// {client1, client2} <-> server
func TestServer_Endpoint(t *testing.T) {
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
	wg.Add(2)

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
		for i := 0; i < 10; i++ {
			sum, err := client.UpdateBalance(ctx, 1, 5)
			assert.Nil(t, err)
			assert.Equal(t, 6, sum)
		}
		wg.Done()
	}()

	// client2
	go func() {
		cid := wirenet.Identification("client2")
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
		for i := 0; i < 10; i++ {
			sum, err := client.UpdateBalance(ctx, 1, 5)
			assert.Nil(t, err)
			assert.Equal(t, 6, sum)
		}
		wg.Done()
	}()
	wg.Wait()

	assert.Nil(t, server.Close())
}

// server <-> {client1, client2}
func TestClient_Endpoint(t *testing.T) {
	addr := randomAddr(t)
	initServer := make(chan struct{})
	wait := func() {
		time.Sleep(5 * time.Second)
	}
	// server side
	server, err := wirenet.Mount(addr,
		wirenet.WithConnectHook(func(closer io.Closer) {
			close(initServer)
		}),
	)
	assert.Nil(t, err)
	go func() {
		assert.Nil(t, server.Connect())
	}()
	<-initServer

	// request from server to {client1, client2}
	go func() {
		time.Sleep(3 * time.Second)
		client := test.MakeWirenetClient(server)
		ctx := context.Background()

		for i := 0; i < 2; i++ {
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

func randomAddr(t *testing.T) string {
	if t == nil {
		t = new(testing.T)
	}
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	assert.Nil(t, err)
	listener, err := net.ListenTCP("tcp", addr)
	assert.Nil(t, err)
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	return fmt.Sprintf(":%d", port)
}
