package gorpc

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type RpcPing int

func (t *RpcPing) Ping(s string, r *string) error {
	*r = s
	return nil
}

type RpcServer struct {
	config    RpcServerConfig
	logFunc   RpcLogFunc
	isRunning bool
	mutex     *sync.Mutex
}

func NewRpcServer(c RpcServerConfig, f RpcLogFunc) (*RpcServer, error) {
	if c == nil {
		return nil, ErrConfig
	}
	return &RpcServer{
		config:    c,
		logFunc:   f,
		isRunning: false,
		mutex:     &sync.Mutex{},
	}, nil
}

func (c *RpcServer) IsRunning() bool {
	return c.isRunning
}

func (c *RpcServer) GetLinkAddress() string {
	return fmt.Sprintf(":%v", c.config.Port)
}

func (c *RpcServer) log(str string, err error) {
	if c.logFunc != nil {
		c.logFunc(str, err)
	}
}

func (c *RpcServer) run() {
	c.isRunning = true
	c.log(fmt.Sprintf("gorpc: server running on: %v",
		c.config.Port), nil)
	if err := rpc.Register(new(RpcPing)); err != nil {
		c.log(fmt.Sprintf("gorpc: server register error: %v", err), err)
	} else if addr, err := net.ResolveTCPAddr("tcp", c.GetLinkAddress()); err != nil {
		c.log(fmt.Sprintf("gorpc: server resolve address error: %v", err), err)
	} else if listener, err := net.ListenTCP("tcp", addr); err != nil {
		c.log(fmt.Sprintf("gorpc: server listen error: %v", err), err)
	} else {
		for c.IsRunning() {
			if conn, err := listener.Accept(); err != nil {
				time.Sleep(time.Second)
				continue
			} else {
				rpc.ServeConn(conn)
				c.log(fmt.Sprintf("gorpc: server connected: %v",
					conn.RemoteAddr().String()), nil)
			}
		}
	}
	c.isRunning = false
	c.log("gorpc: server close", nil)
}

func (c *RpcServer) Start() {
	if c.IsRunning() {
		return
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !c.IsRunning() {
		go c.run()
	}
}

func (c *RpcServer) Stop() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isRunning = false
}
