package http

import (
	"github.com/clearcodecn/log"
	"github.com/clearcodecn/wetalk/api/model"
	"github.com/clearcodecn/wetalk/configs"
	"github.com/clearcodecn/wetalk/pkg/fs"
	"github.com/clearcodecn/wetalk/pkg/mail"
	"github.com/clearcodecn/wetalk/pkg/sms"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Server struct {
	engine     *gin.Engine
	smsSender  sms.Sender
	mailSender mail.Sender
	mailChan   chan *MailInfo

	uploader fs.Uploader
	model    *model.Model

	mux    sync.Mutex
	config configs.WebConfig

	configReloader func() chan configs.WebConfig
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

func WithConfigReloader(f func() chan configs.WebConfig) Option {
	return func(s *Server) {
		s.configReloader = f
	}
}

func NewServer(config configs.WebConfig, options ...Option) *Server {
	s := new(Server)
	for _, opt := range options {
		opt(s)
	}

	if s.configReloader != nil {
		go s.watchConfig()
	}

	var err error
	s.model, err = model.NewModel(config.DbConfig.Driver, config.DbConfig.Dsn)
	if err != nil {
		panic(err)
	}

	s.config = config
	s.engine = gin.New()
	s.registerRoutes()
	return s
}

func (s *Server) Run() error {

	if s.mailSender != nil {
		s.mailChan = make(chan *MailInfo, 100)
		go s.startUpMailer()
	}

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
		g.POST("/email", s.sendEmailVerifyCode)
	}
	{
		g := s.engine.Group("/api/v1")
		g.PUT("/user", s.userUpdate)
		g.POST("/upload", s.Upload)
	}
}

func (s *Server) watchConfig() {
	if s.configReloader == nil {
		return
	}

	_ = s.configReloader()

	// TODO:: reload config.
}

type MailInfo struct {
	*model.VerifyCode
	Title   string
	Content string
}

func (s *Server) startUpMailer() {
	for vc := range s.mailChan {
		log.Info("send a new mail", zap.Any("mailinfo", vc))
		if err := s.mailSender.Send(vc.User, vc.Title, vc.Content); err != nil {
			log.Error("send email failed", zap.Error(err))
			continue
		}
		vc.CreateTime = time.Now()
		_ = s.model.CreateVerifyCode(vc.VerifyCode)
	}
}
