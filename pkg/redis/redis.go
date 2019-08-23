package redis

import (
	"github.com/clearcodecn/wetalk/configs"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"time"
)

type Client struct {
	redis.UniversalClient
	config configs.RedisConfig
}

func NewClient(config configs.RedisConfig) (*Client, error) {
	cli := new(Client)
	if config.Timeout == 0 {
		config.Timeout = 30
	}
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = time.Second * 30
	}
	cli.UniversalClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        config.Addrs,
		DB:           config.DB,
		Password:     config.Password,
		DialTimeout:  time.Duration(config.Timeout) * time.Second,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		PoolTimeout:  timeout,
		IdleTimeout:  timeout,
		MasterName:   config.MasterName,
		MaxRetries: 2,
	})
	if _, err := cli.UniversalClient.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "ping redis failed")
	}
	return cli, nil
}
