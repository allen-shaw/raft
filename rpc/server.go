package rpc

type Server interface {
	init(opts interface{}) bool
	shutdown()
}
