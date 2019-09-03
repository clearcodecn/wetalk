package server

import (
	"bufio"
	"github.com/clearcodecn/log"
	pb "github.com/clearcodecn/wetalk/proto"
	gogoproto "github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net"
	"sync/atomic"
	"time"
)

const (
	defaultReadSize     = 2048
	defaultWriteSize    = 2048
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
)

var (
	connectionId uint64
)

// bufReadConn wrap net.Conn by bufio.Reader
type bufReadConn struct {
	net.Conn
	r *bufio.Reader
}

// newBufReadConn returns the bufReadConn.
func newBufReadConn(conn net.Conn) *bufReadConn {
	return &bufReadConn{
		Conn: conn,
		r:    bufio.NewReaderSize(conn, defaultReadSize),
	}
}

type clientConn struct {
	id          uint64 // the connection ID.
	bufReadConn *bufReadConn
	pkt         *packet
	server      *Server
	remoteAddr  net.Addr
	ctx         *ConnectionContext
	timeout     time.Duration
	writeChan   chan *pb.Message
}

type ClientConnOptions struct {
	WriteChannelSize int
}

func newClientConn(conn net.Conn, server *Server, options ClientConnOptions) *clientConn {
	client := new(clientConn)
	client.server = server
	client.id = atomic.AddUint64(&connectionId, 1)
	client.bufReadConn = newBufReadConn(conn)
	client.pkt = NewPacket(client.bufReadConn)
	if options.WriteChannelSize <= 0 {
		options.WriteChannelSize = defaultChannelSize
	}
	client.writeChan = make(chan *pb.Message, options.WriteChannelSize)
	go client.doWrite()
	return client
}

func (c *clientConn) ReadMessage(message *pb.Message) error {
	c.bufReadConn.SetReadDeadline(time.Now().Add(c.timeout))
	b, err := c.pkt.ReadPacket()
	if err != nil {
		return errors.Wrap(err, "failed to read message")
	}
	return gogoproto.Unmarshal(b, message)
}

func (c *clientConn) WriteMessage(message *pb.Message) {
	if c.writeChan != nil {
		c.writeChan <- message
	}
}

func (c *clientConn) doWrite() {
	for msg := range c.writeChan {
		data, err := gogoproto.Marshal(msg)
		if err != nil {
			log.Error("write message failed", zap.Any("msg", msg), zap.String("err", err.Error()))
			continue
		}
		err = c.pkt.WritePacket(data)
		if err != nil {
			log.Error("write message failed", zap.Any("msg", msg), zap.String("err", err.Error()))
			return
		}

		err = c.pkt.Flush()
		if err != nil {
			log.Error("write message failed", zap.Any("msg", msg), zap.String("err", err.Error()))
			return
		}
	}
}

func (c *clientConn) Close() error {
	var err error
	close(c.writeChan)

	if c.ctx != nil {
		er := c.ctx.Close()
		if er != nil {
			err = er
		}
	}

	if c.bufReadConn != nil {
		er := c.bufReadConn.Close()
		if er != nil {
			err = er
		}
	}

	return err
}
