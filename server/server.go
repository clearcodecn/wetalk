package server

type Server struct {
	*TCPServer
	*WebsocketServer
}

func newServer() *Server {
	s := new(Server)
	return s
}
