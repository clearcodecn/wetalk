package http

import (
	"github.com/clearcodecn/wetalk/configs"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	config configs.HttpConfig
}

func NewServer(config configs.HttpConfig) *Server {
	s := new(Server)
	s.config = config
	s.engine = gin.New()
	s.registerRoutes()
	return s
}

func (s *Server) Run() error {
	return s.engine.Run(s.config.Addr)
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
	}
}
