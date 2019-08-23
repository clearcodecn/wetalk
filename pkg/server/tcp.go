package server

import (
	"github.com/clearcodecn/wetalk/configs"
	"net"
	"sync/atomic"
	"time"

	"github.com/clearcodecn/log"
	"github.com/clearcodecn/wetalk/proto"
	"go.uber.org/zap"
)

const (
	defaultTimeout = 30 * time.Second
)

type TCPServer struct {
	addr  *net.TCPAddr
	ln    net.Listener
	conns map[string]*Conn

	config  configs.TCPConfig
	timeout time.Duration

	OnAuthorization func()
	OnClose         func()

	connCount uint64
}

func NewTCPServer(config configs.TCPConfig) *TCPServer {
	ts := new(TCPServer)
	if addr, err := net.ResolveTCPAddr("tcp", config.Addr); err != nil {
		panic(err)
	} else {
		ts.addr = addr
	}

	var timeout = time.Duration(config.Timeout) * time.Second
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
		if s.timeout != 0 {
			_ = conn.(*net.TCPConn).SetKeepAlive(true)
			_ = conn.(*net.TCPConn).SetKeepAlivePeriod(s.timeout)
		}
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
	if s.config.MaxConns != 0 && atomic.LoadUint64(&s.connCount) > s.config.MaxConns {
		_ = conn.Close()
		return
	}
	atomic.AddUint64(&s.connCount, 1)
	c := NewConn(conn, s.timeout, defaultChannelSize)
	c.onClose = func() {
		atomic.AddUint64(&s.connCount, -1)
	}
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

func (s *TCPServer) Close() error {
	var (
		err error
	)

	defer func() {
		err = s.ln.Close()
	}()

	for _, conn := range s.conns {
		if err := conn.Close(); err != nil {
			return err
		}
	}

	return err
}
