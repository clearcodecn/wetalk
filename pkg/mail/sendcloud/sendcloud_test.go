package sendcloud

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var configTemplate = `
enable: true
api_user: 'mrjnamei_test_KGqwD4'
api_key: 'z8RKSCFt4mUPevvx'
from: 'admin@wetalk.com'
from_name: WeTalk
`

func TestEmail_Init(t *testing.T) {
	var e = new(Email)
	err := e.Init([]byte(configTemplate))
	require.Nil(t, err)

	err = e.Send("735416909@qq.com", "test-email", "<h1>this is a test email</h1>")
	require.Nil(t, err)
}
