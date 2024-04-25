package api

import (
	"github.com/yanglunara/im/pkg/common/config"
)

type Conf struct {
	Rpc    config.API
	Consul config.Consul
	Share  config.Share
	Log    config.Log
}
