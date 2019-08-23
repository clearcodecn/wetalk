package configs

// UniversalOptions information is required by UniversalClient to establish
// connections.
type RedisConfig struct {
	// 如果是哨兵，哨兵需要配置这个名字.
	MasterName string `json:"master_name" yaml:"master_name"`
	// 如果是集群配置，则配置多个地址
	// 单机配置为1个地址.
	Addrs    []string `json:"addrs" yaml:"addrs"`
	DB       int      `json:"db"`
	Password string   `json:"password" yaml:"password"`
	Timeout  int      `json:"timeout" yaml:"timeout"` // uint: 秒
}

type TCPConfig struct {
	ReadChannelSize  int    `json:"channel_size" yaml:"read_channel_size"`        // 读缓冲区大小
	WriteChannelSize int    `json:"write_channel_size" yaml:"write_channel_size"` // 写缓冲区大小
	Addr             string `json:"addr" yaml:"addr"`                             // 监听地址
	MaxConns         uint64 `json:"max_conns" yaml:"max_conns"`                   // 最大连接数量.
	Timeout          int    `json:"timeout" yaml:"timeout"`                       // 超时控制.
}

type WebsocketConfig struct {
	ReadChannelSize  int    `json:"channel_size" yaml:"read_channel_size"`        // 读缓冲区大小
	WriteChannelSize int    `json:"write_channel_size" yaml:"write_channel_size"` // 写缓冲区大小
	Addr             string `json:"addr" yaml:"addr"`                             // 监听地址
	MaxConns         int    `json:"max_conns" yaml:"max_conns"`                   // 最大连接数量.
	Timeout          int    `json:"timeout" yaml:"timeout"`                       // 超时控制.
}

type HttpConfig struct {
	Addr string `json:"addr"`
}
