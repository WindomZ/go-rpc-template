package gorpc

type rpcConfig struct {
	Host string
	Port int
}

func (c *rpcConfig) Valid() bool {
	return (c != nil && len(c.Host) != 0 && c.Port > 0)
}

type RpcClientConfig struct {
	rpcConfig
	Service  string
	TryCount int
}

func NewRpcClientConfig(service string, host string, port int, tryCount int) RpcClientConfig {
	if len(service) == 0 {
		panic(ErrService)
	} else if len(host) == 0 {
		host = "127.0.0.1"
	}
	return RpcClientConfig{
		rpcConfig: rpcConfig{
			Host: host,
			Port: port,
		},
		Service:  service,
		TryCount: tryCount,
	}
}

type RpcServerConfig struct {
	rpcConfig
	Service IRpcServer
}

func NewRpcServerConfig(service IRpcServer, port int) RpcServerConfig {
	if service == nil {
		panic(ErrService)
	}
	return RpcServerConfig{
		rpcConfig: rpcConfig{
			Host: "127.0.0.1",
			Port: port,
		},
		Service: service,
	}
}
