package model

import (
	"context"

	"github.com/yanglunara/discovery/register"
)

type Replicant interface {
	LoadBalancerUpdate(ins []*register.ServiceInstance)
	OnlineProc(ctx context.Context)
	NodeAddrs(region, domain string, regionWeight float64) (domains, addrs []string)
}
