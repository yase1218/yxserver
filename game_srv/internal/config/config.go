package config

import "time"

var (
	Conf *Config

	// gate conf
	PendingWriteNum        = 2000
	MaxMsgLen       uint32 = 16 * 1024
	HTTPTimeout            = 10 * time.Second
	LenMsgLen              = 4
	LittleEndian           = true

	// skeleton conf
	GoLen              = 10000
	TimerDispatcherLen = 10000
	AsynCallLen        = 10000
	ChanRPCLen         = 10000
)

func init() {
	Conf = new(Config)
}

type Config struct {
	ServerId    uint32 `yaml:"server_id" mapstructure:"server_id"`
	ServerName  string `yaml:"server_name" mapstructure:"server_name"`
	CdkServerId string `yaml:"cdk_server_id" mapstructure:"cdk_server_id"`
	Sandbox     bool   `yaml:"sandbox" mapstructure:"sandbox"`
	Gm          bool   `yaml:"gm" mapstructure:"gm"`
	Whitelist   bool   `yaml:"whitelist" mapstructure:"whitelist"`
	GmTime      bool   `yaml:"gm_time" mapstructure:"gm_time"`
	Env         string `yaml:"env" mapstructure:"env"`
	Prof        bool   `yaml:"prof" mapstructure:"prof"`
	Debug       bool   `yaml:"debug" mapstructure:"debug"`
	Monitor     bool   `yaml:"monitor" mapstructure:"monitor"`

	MaxConnNum     int    `yaml:"max_conn_num" mapstructure:"max_conn_num"`
	MaxRegisterNum uint32 `yaml:"max_register_num" mapstructure:"max_register_num"`

	GoPoolSize uint32 `yaml:"go_pool_size" mapstructure:"go_pool_size"`
	CsvPath    string `yaml:"csv_path" mapstructure:"csv_path"`

	Tcp  TcpConfig  `yaml:"tcp"`
	Ws   WsConfig   `yaml:"ws"`
	Grpc GrpcConfig `yaml:"grpc"`

	Etcd  EtcdConfig  `yaml:"etcd"`
	Redis RedisConfig `yaml:"redis"`
	Nats  NatsConfig  `yaml:"nats"`

	GlobalMongo  MongoConfig `yaml:"global_mongo"  mapstructure:"global_mongo"`
	LocalMongo   MongoConfig `yaml:"local_mongo"  mapstructure:"local_mongo"`
	StaticsMongo MongoConfig `yaml:"statics_mongo"  mapstructure:"statics_mongo"`

	Rbi     RbiConfig  `yaml:"rbi"`
	Leiting LeitingSdk `yaml:"leiting_sdk"  mapstructure:"leiting_sdk"`
	Tap     Tapping    `yaml:"tapping"  mapstructure:"tapping"`
}

type (
	TcpConfig struct {
		Addr     string `yaml:"addr" mapstructure:"addr"`
		LinkAddr string `yaml:"link_addr" mapstructure:"link_addr"`
	}
	WsConfig struct {
		Addr     string `yaml:"addr" mapstructure:"addr"`
		LinkAddr string `yaml:"link_addr" mapstructure:"link_addr"`
		CertFile string `yaml:"cert_file" mapstructure:"cert_file"`
		KeyFile  string `yaml:"key_file" mapstructure:"key_file"`
	}

	GrpcConfig struct {
		Addr string `yaml:"addr" mapstructure:"addr"`
	}

	EtcdConfig struct {
		Addrs []string `yaml:"addrs" mapstructure:"addrs"`
	}

	RedisConfig struct {
		Addr string `yaml:"addr" mapstructure:"addr"`
		Pass string `yaml:"pass" mapstructure:"pass"`
	}

	NatsConfig struct {
		Addr string `yaml:"addr" mapstructure:"addr"`
	}

	MongoConfig struct {
		Url          string `yaml:"url" mapstructure:"url"`
		DB           string `yaml:"db" mapstructure:"db"`
		MaxConn      int    `yaml:"max_conn" mapstructure:"max_conn"`
		WorkCount    int    `yaml:"work_count" mapstructure:"work_count"`
		QueueSize    int    `yaml:"queue_size" mapstructure:"queue_size"`
		BatchSize    int    `yaml:"batch_size" mapstructure:"batch_size"`
		FlushSeconds int    `yaml:"flush_timeout" mapstructure:"flush_timeout"`
	}

	RbiConfig struct {
		Open bool   `yaml:"open" mapstructure:"open"`
		Url  string `yaml:"url" mapstructure:"url"`
		Port string `yaml:"port" mapstructure:"port"`
		Test bool   `yaml:"test" mapstructure:"test"`
	}

	LeitingSdk struct {
		Key          string `yaml:"key" mapstructure:"key"`
		Game         string `yaml:"game" mapstructure:"game"`
		Domain       string `yaml:"domain" mapstructure:"domain"`                 // 异常订单上报地址
		NotifyUrl    string `yaml:"notify_url" mapstructure:"notify_url"`         // sdk验单回调地址
		TextCheckUrl string `yaml:"text_check_url" mapstructure:"text_check_url"` // 文本检测地址
	}

	Tapping struct {
		Switch int32 `yaml:"switch" mapstructure:"switch"` // 1只有sdk打点，2只有本地打点，3sdk和本地都打点
	}
)

func (c *Config) GetLocalDB() string {
	return c.LocalMongo.DB
}
func (c *Config) GetGlobalDB() string {
	return c.GlobalMongo.DB
}
func (c *Config) GetStaticsDB() string {
	return c.StaticsMongo.DB
}

func (c *Config) IsDebug() bool {
	return c.Debug
}
