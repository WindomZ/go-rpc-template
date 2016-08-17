package gorpc

import (
	"fmt"
	"net/rpc"
	"sync"
	"time"
)

const PING string = "PING"

type RpcClient struct {
	client    *rpc.Client
	config    RpcClientConfig
	logFunc   RpcLogFunc
	tryCnt    int
	isRunning bool
	mutex     *sync.Mutex
}

func NewRpcClient(c RpcClientConfig, f RpcLogFunc) (*RpcClient, error) {
	if !c.Valid() {
		return nil, ErrConfig
	}
	return &RpcClient{
		config:    c,
		logFunc:   f,
		isRunning: false,
		mutex:     &sync.Mutex{},
	}, nil
}

func (c *RpcClient) IsRunning() bool {
	return c.isRunning
}

func (c *RpcClient) GetLinkAddress() string {
	return fmt.Sprintf("%v:%v", c.config.Host, c.config.Port)
}

func (c *RpcClient) tryAgain() bool {
	if !c.IsRunning() {
		return false
	} else if c.config.TryCount <= 0 {
		return true
	}
	defer func() {
		c.tryCnt++
	}()
	return c.tryCnt <= c.config.TryCount
}

func (c *RpcClient) log(msg string, err error) {
	if c.logFunc != nil {
		c.logFunc(msg, err)
	}
}

func (c *RpcClient) run() {
	c.isRunning = true
	c.tryCnt = 0
	for c.tryAgain() {
		if c.client == nil {
			var err error
			if c.client, err = rpc.Dial(
				"tcp", c.GetLinkAddress(),
			); err != nil {
				if c.tryCnt%15 == 0 {
					c.log(fmt.Sprintf("gorpc: client dial error: %v", err), err)
				}
			} else {
				c.log(fmt.Sprintf("gorpc: client running on: %v",
					c.GetLinkAddress()), nil)
			}
		} else if r, err := c.Ping(); err != nil {
			c.log(fmt.Sprintf("gorpc: client ping error: %v", err), err)
		} else if len(r) != 0 {
			c.tryCnt = 0
			time.Sleep(time.Second)
			continue
		}
		c.tryCnt++
		time.Sleep(time.Second)
	}
	c.isRunning = false
	c.log("gorpc: client close", nil)
}

func (c *RpcClient) Connect() {
	if c.IsRunning() || c.IsConnected() {
		return
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !c.IsRunning() {
		go c.run()
	}
}

func (c *RpcClient) IsConnected() bool {
	if c.client == nil {
		return false
	} else if _, err := c.Ping(); err == nil {
		return true
	}
	return false
}

func (c *RpcClient) Ping(str ...string) (r string, err error) {
	if len(str) == 0 {
		str = append(str, PING)
	}
	err = c.Call(c.config.Service+".Ping", str[0], &r)
	return
}

func (c *RpcClient) Disconnect() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isRunning = false
}

func (c *RpcClient) Call(serviceMethod string, args interface{}, reply interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.client == nil {
		return ErrClient
	} else if err := c.client.Call(serviceMethod, args, reply); err != nil {
		c.client.Close()
		c.client = nil
	}
	return nil
}
