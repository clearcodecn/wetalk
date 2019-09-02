package sms

type Sender interface {
	// init the config
	Init([]byte) error
	// send a sms to the mobile
	Send(to string, content string) error
}
