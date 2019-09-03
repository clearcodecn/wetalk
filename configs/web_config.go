package configs

import (
	"github.com/clearcodecn/wetalk/pkg/fs"
	"github.com/clearcodecn/wetalk/pkg/fs/qiniu"
	"github.com/clearcodecn/wetalk/pkg/mail"
	"github.com/clearcodecn/wetalk/pkg/mail/sendcloud"
	"github.com/clearcodecn/wetalk/pkg/sms"
	"github.com/clearcodecn/wetalk/pkg/sms/smsbao"
	"gopkg.in/yaml.v2"
)

type WebConfig struct {
	HttpConfig   HttpConfig   `yaml:"http_config" json:"http_config"`
	RedisConfig  RedisConfig  `yaml:"redis_config" json:"redis_config"`
	EmailConfig  EmailConfig  `yaml:"email_config" json:"email_config"`
	SmsConfig    SmsConfig    `yaml:"sms_config" json:"sms_config"`
	DbConfig     DbConfig     `yaml:"db_config" json:"db_config"`
	PushServer   PushServer   `yaml:"push_server" json:"push_server"`
	UploadConfig UploadConfig `yaml:"upload_config" json:"upload_config"`
}

type HttpConfig struct {
	Addr           string `json:"addr" yaml:"addr"`
	Key            string `json:"key" yaml:"key"`
	Cert           string `json:"cert" yaml:"cert"`
	EnableRegister bool   `json:"enable_register" yaml:"enable_register"`
	EnableVerify   bool   `json:"enable_verify" yaml:"enable_verify"`
}

type EmailConfig struct {
	Enable  bool        `yaml:"enable" json:"enable"`
	Driver  string      `yaml:"driver" json:"driver"`
	Content interface{} `yaml:"content" json:"content"`
}

type SmsConfig struct {
	Enable  bool        `yaml:"enable" json:"enable"`
	Driver  string      `yaml:"driver" json:"driver"`
	Content interface{} `yaml:"content" json:"content"`
}

type PushServer struct {
	Addr []string `json:"addr" yaml:"addr"`
}

type DbConfig struct {
	Driver string `json:"driver"`
	Dsn    string `json:"dsn"`
}

type UploadConfig struct {
	Driver  string      `json:"driver"`
	Content interface{} `yaml:"content" json:"content"`
}

func ParseWebConfig(data []byte) (WebConfig, error) {
	config := WebConfig{}
	err := yaml.Unmarshal(data, &config)
	return config, err
}

var (
	Uploaders = map[string]func() fs.Uploader{
		"qiniu": func() fs.Uploader {
			return new(qiniu.Uploader)
		},
	}

	MailSenders = map[string]func() mail.Sender{
		"sendcloud": func() mail.Sender { return new(sendcloud.Email) },
	}

	SmsSenders = map[string]func() sms.Sender{
		"smsbao": func() sms.Sender { return new(smsbao.Sms) },
	}
)
