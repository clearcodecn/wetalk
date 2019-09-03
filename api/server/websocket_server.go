package server

import (
	"github.com/clearcodecn/log"
	"github.com/clearcodecn/wetalk/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

// websocketConn is a wrap of websocket.Conn,
// it impl net.Conn interface
type websocketConn struct {
	conn    *websocket.Conn
	msgType int
}

func (w *websocketConn) Read(b []byte) (n int, err error) {
	_, r, err := w.conn.NextReader()
	if err != nil {
		return 0, err
	}
	return r.Read(b)
}

func (w *websocketConn) Write(b []byte) (n int, err error) {
	wt, err := w.conn.NextWriter(w.msgType)
	if err != nil {
		return 0, nil
	}
	return wt.Write(b)
}

func (w *websocketConn) Close() error {
	return w.conn.Close()
}

func (w *websocketConn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *websocketConn) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *websocketConn) SetDeadline(t time.Time) error {
	if err := w.conn.SetWriteDeadline(t); err != nil {
		return err
	}
	if err := w.conn.SetReadDeadline(t); err != nil {
		return err
	}
	return nil
}

func (w *websocketConn) SetReadDeadline(t time.Time) error {
	return w.conn.SetReadDeadline(t)
}

func (w *websocketConn) SetWriteDeadline(t time.Time) error {
	return w.conn.SetWriteDeadline(t)
}

func newWebsocketConn(conn *websocket.Conn, msgType int) net.Conn {
	return &websocketConn{
		conn:    conn,
		msgType: msgType,
	}
}

type WebsocketServer struct {
	upgrader websocket.Upgrader
	timeout  time.Duration
	msgType  int
}

func NewWebsocketServer(timeout time.Duration) *WebsocketServer {
	s := new(WebsocketServer)
	s.upgrader = websocket.Upgrader{
		HandshakeTimeout: timeout,
		ReadBufferSize:   defaultReadSize,
		WriteBufferSize:  defaultReadSize,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func (s *WebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.writeBadRequest(w)
		return
	}
	wsConn := newWebsocketConn(conn, s.msgType)
	c := NewConn(wsConn, s.timeout, defaultChannelSize)
	// TODO auth
	msg := &pb.Message{}
	err = c.ReadMessage(msg)
	if err != nil {
		log.Error("read auth message failed", zap.Error(err))
		return
	}
	return
}

func (s *WebsocketServer) writeBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("Bad Request"))
}
