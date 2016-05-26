package gorpc

import "errors"

var (
	ErrClient  error = errors.New("gorpc: Did not start the client")
	ErrConnect       = errors.New("gorpc: No connected rpc server")
	ErrConfig        = errors.New("gorpc: Invalid config")
	ErrService       = errors.New("gorpc: Invalid service")
)
