package conf

import (
	"time"

	conf "github.com/yanglunara/im/internal/conf"
	"github.com/yunbaifan/pkg/logger"
)

var (
	Conf *Config
)

func Init(filePath string) error {
	if err := conf.Binding(&Conf, filePath); err != nil {
		return err
	}
	return nil
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

type GrpcServer struct {
	Network string        `json:"network" yaml:"network" description:"网络" default:"grpc"`
	Addr    string        `json:"addr" yaml:"addr" description:"服务地址" default:"0.0.0.0:9002"`
	Timeout time.Duration `json:"timeout" yaml:"timeout" description:"超时时间" default:"10s"`
}

type GlobalEnv struct {
	Region    string `yaml:"region" description:"区域" default:"sh"`
	Zone      string `yaml:"zone" description:"可用区" default:"sh001"`
	DeployEnv string `yaml:"deployEnv" description:"部署环境" default:"dev"`
	Host      string `yaml:"host" description:"主机名" default:"localhost"`
	Weight    int64  `yaml:"weight" description:"权重" default:"10"`
}
type Config struct {
	Logger     logger.LogConfig    `yaml:"logger"`
	GrpcServer GrpcServer          `yaml:"grpcServer" description:"grpc服务"`
	Consul     Consul              `yaml:"consul" description:"consul配置"`
	Kafka      Kafka               `yaml:"kafka" description:"kafka配置"`
	Redis      Redis               `yaml:"redis" description:"redis配置"`
	GlobalEnv  GlobalEnv           `yaml:"globalEnv" description:"全局环境"`
	Regions    map[string][]string `yaml:"regions" description:"区域配置"`
	Node       Node                `yaml:"nodes" description:"节点配置"`
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
