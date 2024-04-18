package conf

import (
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	once     sync.Once
	Conf     *Config
	path     string
	region   string
	zone     string
	env      string
	hostName string
	weight   int64
	name     string
)

func init() {
	// var (
	// 	hostName, _ = os.Hostname()                                 // 获取主机名
	// 	weight, _   = strconv.ParseInt(os.Getenv("WEIGHT"), 10, 32) // 从环境变量中获取权重，并将其转换为 int64
	// )
	// // 使用 flag 包来解析命令行参数。这些参数可以用来覆盖默认的配置值。
	// // 配置文件路径
	// flag.StringVar(&path, "conf", "config.toml", "config file path")
	// // 区域名称
	// flag.StringVar(&region, "region", "local", "region name")
	// // 区域名称
	// flag.StringVar(&zone, "zone", "local", "zone name")
	// // 环境名称
	// flag.StringVar(&env, "env", "dev", "env name")
	// // 主机名
	// flag.StringVar(&hostName, "host", hostName, "host name")
	// // 区域权重
	// flag.Int64Var(&weight, "weight", weight, "region weight")
	// // 解析命令行参数
	// flag.Parse()
}

func Init() (err error) {
	once.Do(func() {
		Conf = Default()
		_, err = toml.DecodeFile(path, &Conf)
	})
	return
}

func Default() *Config {
	defaultAddr := ":9000"
	return &Config{
		Env: &Env{
			Region:    region,
			Zone:      zone,
			DeployEnv: env,
			Host:      hostName,
			Weight:    weight,
		},
		Grpc: &GRPCServer{
			Network:           "0.0.0.0",
			Address:           defaultAddr,
			Timeout:           60 * time.Second,
			MaxLiftTime:       5 * time.Second,
			IdleTimeout:       2 * time.Hour,
			ForceCloseWait:    20 * time.Second,
			KeepaliveInterval: 60 * time.Second,
			KeepaliveTimeout:  20 * time.Second,
		},
		// 服务发现与注册
		Discovery: &Discovery{
			ID:        hostName,
			Name:      name,
			Metadata:  make(map[string]string),
			Endpoints: []string{defaultAddr},
		},
		Consul: &Consul{
			Address: "127.0.0.1:8500",
		},
		Backoff: &Backoff{
			MaxDelay:  300,
			BaseDelay: 3,
			Factor:    1.8,
			Jitter:    1.3,
		},
	}
}

type Consul struct {
	Address string
}

type Discovery struct {
	ID        string
	Name      string
	Version   string
	Metadata  map[string]string
	Endpoints []string
}

type GRPCServer struct {
	Network           string        `json:"network"`
	Address           string        `json:"address"`
	Timeout           time.Duration `json:"timeout"`
	MaxLiftTime       time.Duration `json:"maxLifeTime"`
	IdleTimeout       time.Duration `json:"idleTimeout"`
	ForceCloseWait    time.Duration `json:"forceCloseWait"`
	KeepaliveInterval time.Duration `json:"keepaliveInterval"`
	KeepaliveTimeout  time.Duration `json:"keepaliveTimeout"`
}

type Env struct {
	Region    string
	Zone      string
	DeployEnv string
	Host      string
	Weight    int64
}
type Config struct {
	Kafka      *Kafka
	Redis      *Redis
	Node       *Node
	Env        *Env
	Grpc       *GRPCServer
	Consul     *Consul
	Discovery  *Discovery
	Backoff    *Backoff
	GrpcClient *GrpcClient
}

type Kafka struct {
	Topic   string
	Group   string
	Brokers []string
}

type Redis struct {
	Network      string
	Addr         string
	Auth         string
	Active       int
	Idle         int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Expire       time.Duration
}

type Node struct {
	DefaultName  string
	HostName     string
	TcpPort      int
	WsPort       int
	HeartbeatMax int
	Heartbeat    int
	RegionWeight float64
}

type Backoff struct {
	MaxDelay  int32
	BaseDelay int32
	Factor    float64
	Jitter    float64
}

type GrpcClient struct {
	Addr    string // gRPC 服务器的地址
	Timeout int64
}
