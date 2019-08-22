package server

import (
	"net"
	"time"

	"github.com/clearcodecn/log"
	"github.com/clearcodecn/wetalk/proto"
	"go.uber.org/zap"
)

const (
	defaultTimeout = 30 * time.Second
)

type TCPServer struct {
	addr    *net.TCPAddr
	timeout time.Duration
	ln      net.Listener
	conns   map[string]*Conn
}

func NewTCPServer(addr string, timeout time.Duration) *TCPServer {
	ts := new(TCPServer)
	if addr, err := net.ResolveTCPAddr("tcp", addr); err != nil {
		panic(err)
	} else {
		ts.addr = addr
	}
	if timeout == 0 {
		timeout = defaultTimeout
	}
	ts.timeout = timeout
	return ts
}

func (s *TCPServer) Run() error {
	var (
		ln  net.Listener
		err error
	)
	if ln, err = net.ListenTCP("tcp", s.addr); err != nil {
		return err
	}
	s.ln = ln
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, e := s.ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return e
		}

		go s.handleConnection(conn)
	}
}

// nolint
func (s *TCPServer) handleConnection(conn net.Conn) {
	c := NewConn(conn, s.timeout, defaultChannelSize)
	// TODO auth
	msg := &proto.Message{}
	err := c.ReadMessage(msg)
	if err != nil {
		log.Error("read auth message failed", zap.Error(err))
		return
	}
	s.conns["random"] = c
	return
}
