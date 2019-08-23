package redis

import (
	"github.com/clearcodecn/wetalk/configs"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	RedisHost     string
	RedisPassword string
	MasterName    string
)

func init() {
	RedisHost = os.Getenv("REDIS_HOST")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	MasterName = os.Getenv("REDIS_MASTER")
}

func TestNewClient(t *testing.T) {
	if RedisHost == "" {
		t.Skipf("redis_host not configed, skipped")
	}
	cfg := configs.RedisConfig{
		MasterName: MasterName,
		Addrs:      []string{RedisHost},
		DB:         0,
		Password:   RedisPassword,
		Timeout:    0,
	}

	cli, err := NewClient(cfg)
	require.Nil(t, err)
	require.NotNil(t, cli)

	_, err = cli.Set("foo", "bar", 0).Result()
	require.Nil(t, err)

	bar, err := cli.Get("foo").Result()
	require.Nil(t, err)
	require.Equal(t, bar, "bar")
}
