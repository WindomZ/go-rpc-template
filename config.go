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
	TryCount int
}

func NewRpcClientConfig(host string, port, tryCount int) RpcClientConfig {
	if len(host) == 0 {
		host = "127.0.0.1"
	}
	return RpcClientConfig{
		rpcConfig: rpcConfig{
			Host: host,
			Port: port,
		},
		TryCount: tryCount,
	}
}

type RpcServerConfig struct {
	rpcConfig
}

func NewRpcServerConfig(port int) RpcServerConfig {
	return RpcServerConfig{
		rpcConfig: rpcConfig{
			Host: "127.0.0.1",
			Port: port,
		},
	}
}
