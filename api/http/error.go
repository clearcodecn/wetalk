package http

import "fmt"

type ErrMsg struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e ErrMsg) Error() string {
	return e.Message
}

func (e ErrMsg) String() string {
	return fmt.Sprintf("code: %d, message:%s", e.Code, e.Message)
}

const (
	Ok int = iota
	Fail
)

func newError(code int, message string) error {
	return &ErrMsg{Code: code, Message: message}
}

func success(message string) error {
	return &ErrMsg{Code: Ok, Message: message}
}

func fail(message string, err error) error {
	if err != nil {
		return &ErrMsg{Code: Fail, Message: fmt.Sprintf("%s : %s", message, err)}
	}
	return &ErrMsg{Code: Fail, Message: message}
}
