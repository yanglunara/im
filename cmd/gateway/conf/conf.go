package conf

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
