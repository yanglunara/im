package conf

var (
	Conf *Config
)

type Config struct {
	GrpcClient      *GrpcClient
	Bucket          *Bucket
	GRPCServer      *GRPCServer
	Consul          *Consul
	Protocol        *Protocol
	GracefulTimeout int `json:"gracefulTimeout" yaml:"gracefulTimeout" description:"优雅关闭超时时间" default:"10"`
}

type GrpcClient struct {
	Addr    string // gRPC 服务器的地址
	Timeout int64
}

type GRPCServer struct {
	Network           string `json:"network"`
	Address           string `json:"address"`
	Timeout           int    `json:"timeout"`
	MaxLiftTime       int    `json:"maxLifeTime"`
	IdleTimeout       int    `json:"idleTimeout"`
	ForceCloseWait    int    `json:"forceCloseWait"`
	KeepaliveInterval int    `json:"keepaliveInterval"`
	KeepaliveTimeout  int    `json:"keepaliveTimeout"`
}

type Bucket struct {
	Size          int
	Channel       int
	Room          int
	RoutineAmount uint64
	RoutineSize   int
}

type Consul struct {
	Address string
}

type Protocol struct {
	SvrProto int
	CliProto int
}
