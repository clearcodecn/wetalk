package server

import (
	"github.com/clearcodecn/log"
	"github.com/clearcodecn/wetalk/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type WebsocketServer struct {
	upgrader websocket.Upgrader
	timeout  time.Duration
}

func NewWebsocketServer(timeout time.Duration) *WebsocketServer {
	s := new(WebsocketServer)
	s.upgrader = websocket.Upgrader{
		HandshakeTimeout:  timeout,
		ReadBufferSize:    0,
		WriteBufferSize:   0,
		WriteBufferPool:   nil,
		Subprotocols:      nil,
		Error:             nil,
		CheckOrigin:       nil,
		EnableCompression: false,
	}
}

func (s *WebsocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.writeBadRequest(w)
		return
	}
	c := NewConn(conn, s.timeout, defaultChannelSize)
	// TODO auth
	msg := &proto.Message{}
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
