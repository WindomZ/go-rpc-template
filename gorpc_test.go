package gorpc

import (
	"testing"
	"time"
)

type TestRpcPing int

func (t *TestRpcPing) Ping(s string, r *string) error {
	*r = s
	return nil
}

func TestRpc(t *testing.T) {
	interrupt := make(chan error)
	service := new(TestRpcPing)
	s, err := NewRpcServer(NewRpcServerConfig(service, 10800), func(msg string, err error) {
		t.Logf("Server: %v - %v", msg, err)
		if err != nil {
			interrupt <- err
		}
	})
	if err != nil {
		t.Fatal(err)
	} else {
		s.Start()
	}
	c, err := NewRpcClient(NewRpcClientConfig("TestRpcPing", "", 10800, 0), func(msg string, err error) {
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
