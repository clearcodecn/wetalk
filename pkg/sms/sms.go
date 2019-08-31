package sms

import "github.com/clearcodecn/wetalk/pkg/sms/smsbao"

var (
	Senders = map[string]func() Sender{
		"smsbao": func() Sender { return new(smsbao.Sms) },
	}
)

type Sender interface {
	// init the config
	Init([]byte) error
	// send a sms to the mobile
	Send(to string, content string) error
}
