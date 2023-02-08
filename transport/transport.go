package transport

type Transport struct {
	clients map[string]*Client
	server  *Server
}

func NewTrasnport() *Transport {

}
