package mail


type Sender interface {
	// init the config
	Init([]byte) error
	// send a email to param to
	Send(to string, title string, content string) error
}
