package smsbao

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	cfgTemplate = `
username: '***'
password: '***'
`
)

const (
	mobile  = `***`
	content = `测试短信内容`
)

func TestSms_Send(t *testing.T) {
	s := new(Sms)
	err := s.Init([]byte(cfgTemplate))
	require.Nil(t, err)

	err = s.Send(mobile, content)
	require.Nil(t, err)
}
