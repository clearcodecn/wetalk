package sendcloud

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"net/http"
	"net/url"
	"strings"
)

const (
	sendURL = `http://api.sendcloud.net/apiv2/mail/send`
)

type EmailConfig struct {
	Fake     bool   `json:"fake" yaml:"fake"` // 如果为true. 将直接返回，而不是真正的发送.
	ApiUser  string `json:"api_user" yaml:"api_user"`
	ApiKey   string `json:"api_key" yaml:"api_key"`
	From     string `json:"from" yaml:"from"`
	FromName string `json:"from_name" yaml:"from_name"`
}

type Email struct {
	cfg *EmailConfig
}

func (e *Email) Init(b []byte) error {
	var cfg = new(EmailConfig)
	err := yaml.Unmarshal(b, cfg)
	if err != nil {
		return err
	}
	if cfg.ApiKey == "" {
		return errors.New(`invalid email config: api_key`)
	}
	if cfg.ApiUser == "" {
		return errors.New(`invalid email config: api_user`)
	}
	if cfg.From == "" {
		return errors.New(`invalid email config: from`)
	}
	e.cfg = cfg
	return nil
}

type SendReply struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Result     bool   `json:"result"`
}

func (e *Email) Send(email string, title string, content string) error {
	if e.cfg.Fake {
		return nil
	}
	data := e.makeParams(email, title, content)
	req, err := http.NewRequest(http.MethodPost, sendURL, strings.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	sp := new(SendReply)
	err = json.NewDecoder(resp.Body).Decode(sp)
	if err != nil {
		return err
	}
	if sp.StatusCode != 200 || !sp.Result {
		return fmt.Errorf("send failed: code: %d, message: %s", sp.StatusCode, sp.Message)
	}
	return nil
}

func (e *Email) makeParams(email string, title, content string) string {
	v := &url.Values{}
	v.Add("apiUser", e.cfg.ApiUser)
	v.Add("apiKey", e.cfg.ApiKey)
	v.Add("from", e.cfg.From)
	v.Add("fromName", e.cfg.FromName)
	v.Add("subject", title)
	v.Add("html", content)
	v.Add("to", email)
	return v.Encode()
}
