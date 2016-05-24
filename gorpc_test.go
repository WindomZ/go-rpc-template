package gorpc

import "testing"

func TestRpc(t *testing.T) {
	interrupt := make(chan error)
	if s, err := NewRpcServer(NewRpcServerConfig(10800), func(msg string, err error) {
		t.Logf("Server: %v - %v", msg, err)
	}); err != nil {
		t.Fatal(err)
	} else {
		s.Start()
	}
	if c, err := NewRpcClient(NewRpcClientConfig("", 10800, 0), func(msg string, err error) {
		t.Logf("Client: %v - %v", msg, err)
	}); err != nil {
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
	}
}
