package mail

import "github.com/clearcodecn/wetalk/pkg/mail/sendcloud"

var (
	Senders = map[string]func() Sender{
		"sendcloud": func() Sender { return new(sendcloud.Email) },
	}
)

type Sender interface {
	// init the config
	Init([]byte) error
	// send a email to param to
	Send(to string, title string, content string) error
}
