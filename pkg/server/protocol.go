package server

import (
	goio "io"
	"net"
	"time"

	"github.com/clearcodecn/log"
	bool2 "github.com/clearcodecn/wetalk/pkg/bool"
	"github.com/clearcodecn/wetalk/pkg/io"
	"github.com/clearcodecn/wetalk/proto"
	gogoproto "github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	defaultChannelSize = 1024
)

var (
	ErrConnectionClosed = errors.New("connect is closed")
)

type Conn struct {
	timeout   time.Duration
	conn      interface{} // it is a net.Conn or websocket.Conn
	readChan  chan *proto.Message
	writeChan chan *proto.Message
	closed    *bool2.Bool

	rw *io.ReadWriter
}

func (c *Conn) ReadMessage(message *proto.Message) error {
	if c.closed.True() {
		return ErrConnectionClosed
	}
	c.setReadDeadLine()
	b, err := c.rw.ReadPacket()
	if err != nil {
		return errors.Wrap(err, "failed to read message")
	}
	return gogoproto.Unmarshal(b, message)
}

func (c *Conn) WriteMessage(message *proto.Message) error {
	if c.closed.True() {
		return ErrConnectionClosed
	}
	if c.writeChan != nil {
		c.writeChan <- message
		return nil
	}

	panic("bug")
}

func NewConn(conn interface{}, timeout time.Duration, channelSize int) *Conn {
	if channelSize <= 0 {
		channelSize = defaultChannelSize
	}

	c := &Conn{
		timeout:   timeout,
		conn:      conn,
		readChan:  make(chan *proto.Message, channelSize),
		writeChan: make(chan *proto.Message, channelSize),
		rw:        io.NewReadWriter(conn.(goio.ReadWriter), false),
	}

	go c.writePending()

	return c
}

func (c *Conn) writePending() {
	for m := range c.writeChan {
		data, err := gogoproto.Marshal(m)
		if err != nil {
			log.Error("marshal message failed", zap.Error(err))
			break
		}
		err = c.rw.WritePacket(data)
		if err != nil {
			log.Error("marshal message failed", zap.Error(err))
			return
		}
		err = c.rw.Flush()
		if err != nil {
			log.Error("marshal message failed", zap.Error(err))
			return
		}
	}
}

func (c *Conn) setReadDeadLine() {
	if conn, ok := c.conn.(net.Conn); ok {
		_ = conn.SetReadDeadline(time.Now().Add(c.timeout))
	}
}

// nolint
func (c *Conn) setWriteDeadLine() {
	if conn, ok := c.conn.(net.Conn); ok {
		_ = conn.SetWriteDeadline(time.Now().Add(c.timeout))
	}
}

func (c *Conn) Close() error {
	if c.closed.True() {
		return errors.New("connection already closed")
	}
	c.closed.SetValue(bool2.True)
	close(c.readChan)
	close(c.writeChan)

	if conn, ok := c.conn.(net.Conn); ok {
		return conn.Close()
	}

	return nil
}
