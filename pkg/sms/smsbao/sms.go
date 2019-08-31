package smsbao

import (
	"fmt"
	"github.com/clearcodecn/wetalk/pkg/util"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

const (
	sendURL = `http://api.smsbao.com/sms?u=%s&p=%s&m=%s&c=%s`
)

var (
	errCodeMap = map[string]string{
		"30": "错误密码",
		"40": "账号不存在",
		"41": "余额不足",
		"43": "IP地址限制",
		"50": "内容含有敏感词",
		"51": "手机号码不正确",
	}
)

type Config struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type Sms struct {
	cfg *Config
}

func (s *Sms) Init(b []byte) error {
	cfg := new(Config)
	err := yaml.Unmarshal(b, cfg)
	if err != nil {
		return err
	}
	if cfg.Username == "" {
		return errors.New(`invalid username`)
	}
	if cfg.Password == "" {
		return errors.New("invalid password")
	}
	cfg.Password = util.Md5(cfg.Password)
	s.cfg = cfg
	return nil
}

func (s *Sms) Send(to string, content string) error {
	url := fmt.Sprintf(sendURL, s.cfg.Username, s.cfg.Password, to, content)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(data) == "0" {
		return nil
	}
	if e, ok := errCodeMap[string(data)]; ok {
		return errors.New(e)
	}
	return fmt.Errorf("unknown error: %s", data)
}

