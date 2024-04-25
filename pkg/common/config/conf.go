package config

type API struct {
	TCP struct {
		ListenIP string `mapstructure:"listen_ip"`
		Ports    []int  `mapstructure:"ports"`
	} `mapstructure:"tcp"`
	Prometheus struct {
		Enable     bool   `mapstructure:"enable"`
		Ports      []int  `mapstructure:"ports"`
		GrafanaURI string `mapstructure:"grafanaURI"`
	} `mapstructure:"prometheus"`
}

type Consul struct {
	Address string `mapstructure:"address"`
}

type Share struct {
	Secret          string          `mapstructure:"secret"`
	Env             string          `mapstructure:"env"`
	RpcRegisterName RpcRegisterName `mapstructure:"rpcRegisterName"`
	ImAdminUserID   []string        `mapstructure:"imAdminUserID"`
}

type RpcRegisterName struct {
	User           string `mapstructure:"user"`
	Friend         string `mapstructure:"friend"`
	Msg            string `mapstructure:"msg"`
	Push           string `mapstructure:"push"`
	MessageGateway string `mapstructure:"messageGateway"`
	Group          string `mapstructure:"group"`
	Auth           string `mapstructure:"auth"`
	Conversation   string `mapstructure:"conversation"`
	Third          string `mapstructure:"third"`
}

type Log struct {
	StorageLocation     string `mapstructure:"storageLocation"`
	RotationTime        uint   `mapstructure:"rotationTime"`
	RemainRotationCount uint   `mapstructure:"remainRotationCount"`
	RemainLogLevel      int    `mapstructure:"remainLogLevel"`
	IsStdout            bool   `mapstructure:"isStdout"`
	IsJson              bool   `mapstructure:"isJson"`
	WithStack           bool   `mapstructure:"withStack"`
}
