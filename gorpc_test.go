package gorpc

import (
	"testing"
	"time"
)

func TestRpc(t *testing.T) {
	interrupt := make(chan error)
	s, err := NewRpcServer(NewRpcServerConfig(10800), func(msg string, err error) {
		t.Logf("Server: %v - %v", msg, err)
		if err != nil {
			interrupt <- err
		}
	})
	if err != nil {
		t.Fatal(err)
	} else {
		s.Start(new(RpcPing))
	}
	c, err := NewRpcClient(NewRpcClientConfig("", 10800, 0), func(msg string, err error) {
		t.Logf("Client: %v - %v", msg, err)
		if err != nil {
			interrupt <- err
		}
	})
	if err != nil {
		t.Fatal(err)
	} else {
		c.Connect()
		defer func() {
			c.Disconnect()
		}()
	}
	select {
	case err := <-interrupt:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(30 * time.Second):
		c.Disconnect()
		s.Stop()
		close(interrupt)
	}
}
