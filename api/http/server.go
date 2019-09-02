package http

import (
	"github.com/clearcodecn/wetalk/configs"
	"github.com/clearcodecn/wetalk/pkg/fs"
	"github.com/clearcodecn/wetalk/pkg/mail"
	"github.com/clearcodecn/wetalk/pkg/sms"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine     *gin.Engine
	config     configs.WebConfig
	smsSender  sms.Sender
	mailSender mail.Sender
	uploader   fs.Uploader
}

type Option func(server *Server)

func WithSmsSender(sender sms.Sender) Option {
	return func(s *Server) {
		s.smsSender = sender
	}
}

func WithEmailSender(sender mail.Sender) Option {
	return func(server *Server) {
		server.mailSender = sender
	}
}

func WithUploader(uploader fs.Uploader) Option {
	return func(s *Server) {
		s.uploader = uploader
	}
}

func NewServer(config configs.WebConfig, options ...Option) *Server {
	s := new(Server)
	s.config = config
	s.engine = gin.New()
	s.registerRoutes()
	return s
}

func (s *Server) Run() error {
	if s.config.HttpConfig.Key != "" && s.config.HttpConfig.Cert != "" {
		return s.engine.RunTLS(s.config.HttpConfig.Addr, s.config.HttpConfig.Cert, s.config.HttpConfig.Key)
	}
	return s.engine.Run(s.config.HttpConfig.Addr)
}

func (s *Server) registerRoutes() {
	{
		g := s.engine.Group("/api/v1")
		g.POST("/login", s.login)
		g.POST("/register", s.register)
	}
	{
		g := s.engine.Group("/api/v1")
		g.PUT("/user", s.userUpdate)
		g.POST("/upload", s.Upload)
	}
}
