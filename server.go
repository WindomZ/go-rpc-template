package gorpc

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type IRpcServer interface {
	Ping(string, *string) error
}

type RpcServer struct {
	config    RpcServerConfig
	logFunc   RpcLogFunc
	isRunning bool
	mutex     *sync.Mutex
}

func NewRpcServer(c RpcServerConfig, f RpcLogFunc) (*RpcServer, error) {
	if !c.Valid() {
		return nil, ErrConfig
	}
	return &RpcServer{
		config:    c,
		logFunc:   f,
		isRunning: false,
		mutex:     &sync.Mutex{},
	}, nil
}

func (s *RpcServer) IsRunning() bool {
	return s.isRunning
}

func (s *RpcServer) GetLinkAddress() string {
	return fmt.Sprintf(":%v", s.config.Port)
}

func (s *RpcServer) log(msg string, err error) {
	if s.logFunc != nil {
		s.logFunc(msg, err)
	}
}

func (s *RpcServer) run() {
	s.isRunning = true
	s.log(fmt.Sprintf("gorpc: server running on: %v",
		s.config.Port), nil)
	if err := rpc.Register(s.config.Service); err != nil {
		s.log(fmt.Sprintf("gorpc: server register error: %v", err), err)
	} else if addr, err := net.ResolveTCPAddr("tcp", s.GetLinkAddress()); err != nil {
		s.log(fmt.Sprintf("gorpc: server resolve address error: %v", err), err)
	} else if listener, err := net.ListenTCP("tcp", addr); err != nil {
		s.log(fmt.Sprintf("gorpc: server listen error: %v", err), err)
	} else {
		for s.IsRunning() {
			if conn, err := listener.Accept(); err != nil {
				time.Sleep(time.Second)
				continue
			} else {
				rpc.ServeConn(conn)
				s.log(fmt.Sprintf("gorpc: server connected: %v",
					conn.RemoteAddr().String()), nil)
			}
		}
	}
	s.isRunning = false
	s.log("gorpc: server close", nil)
}

func (s *RpcServer) Start() {
	if s.IsRunning() {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.IsRunning() {
		go s.run()
	}
}

func (s *RpcServer) Stop() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.isRunning = false
}
