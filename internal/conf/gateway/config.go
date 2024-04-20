package conf

import (
	"time"

	conf "github.com/yanglunara/im/internal/conf"
	"github.com/yunbaifan/pkg/logger"
)

var (
	Conf *Config
)

// Init 初始化配置
func Init(filePath string) error {
	if err := conf.Binding(&Conf, filePath); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Logger     logger.LogConfig `json:"logger" yaml:"logger"`
	GrpcClient GrpcClient       `yaml:"grpcClient" description:"grpc客户端"`
	Bucket     Bucket           `yaml:"bucket" description:"桶配置"`
	GrpcServer GrpcServer       `yaml:"grpcServer" description:"grpc服务"`
	HTTPServer HTTPServer       `yaml:"httpServer" description:"http服务"`
	Consul     Consul           `yaml:"consul" description:"consul配置"`
	Protocol   Protocol         `yaml:"protocol" description:"协议配置"`
	GlobalEnv  GlobalEnv        `yaml:"globalEnv" description:"全局环境"`
}

type GrpcClient struct {
	Addr    string        `json:"addr" yaml:"addr" description:"服务地址" default:"im.logic"`
	Timeout time.Duration `json:"timeout" yaml:"timeout" description:"超时时间" default:"1s"`
}

type GrpcServer struct {
	Network string        `json:"network" yaml:"network" description:"网络" default:"grpc"`
	Addr    string        `json:"addr" yaml:"addr" description:"服务地址" default:"0.0.0.0:9002"`
	Timeout time.Duration `json:"timeout" yaml:"timeout" description:"超时时间" default:"10s"`
}

type HTTPServer struct {
	Name            string        `json:"name" yaml:"name" description:"服务名称"`
	Host            string        `json:"host" yaml:"host" description:"服务地址" default:"0.0.0.0"`
	Port            int           `json:"port" yaml:"port" description:"服务端口" default:"8888"`
	JwtSecret       string        `json:"jwtSecret" yaml:"jwtSecret" description:"jwt密钥"`
	Expire          time.Duration `json:"expire" yaml:"expire" description:"jwt过期时间" default:"7200"`
	GracefulTimeout time.Duration `json:"gracefulTimeout" yaml:"gracefulTimeout" description:"优雅关闭超时时间" default:"10"`
}

type Bucket struct {
	Size          int    `json:"size" yaml:"size" description:"桶大小" default:"32"`
	Channel       int    `json:"channel" yaml:"channel" description:"通道" default:"1024"`
	Room          int    `json:"room" yaml:"room" description:"房间" default:"1024"`
	RoutineAmount uint64 `json:"routineAmount" yaml:"routineAmount" description:"协程数量" default:"32"`
	RoutineSize   int    `json:"routineSize" yaml:"routineSize" description:"协程大小" default:"1024"`
}

type Consul struct {
	Address string `json:"address" yaml:"address" description:"consul地址"`
}

type Protocol struct {
	SvrProto int `json:"svrProto" yaml:"svrProto" description:"服务协议" default:"10"`
	CliProto int `json:"cliProto" yaml:"cliProto" description:"客户端协议" default:"5"`
}

type GlobalEnv struct {
	Region    string   `yaml:"region" description:"区域" default:"sh"`
	Zone      string   `yaml:"zone" description:"可用区" default:"sh001"`
	DeployEnv string   `yaml:"deployEnv" description:"部署环境" default:"dev"`
	Host      string   `yaml:"host" description:"主机名" default:"localhost"`
	Weight    int64    `yaml:"weight" description:"权重" default:"10"`
	Addrs     []string `yaml:"addrs" description:"服务地址" default:"localhost"`
}
