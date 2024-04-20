package logic

import (
	"context"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/yanglunara/discovery/builder"
	"github.com/yanglunara/discovery/register"
	"github.com/yanglunara/im/internal/model"
)

var (
	_ model.Replicant = (*replicant)(nil)
)

type replicant struct {
	onLineTick  time.Duration
	event       chan []*register.ServiceInstance
	discovery   register.Discovery
	cancel      context.CancelFunc
	serviceName string
	allIns      []*register.ServiceInstance
	totalConns  int64
	totalIPs    int64
	nodesMutex  sync.Mutex
	nodes       map[string]*weighted
	totalWeight int64
}

func newReplicant(endpoint, serviceName string) model.Replicant {
	r := &replicant{
		onLineTick: time.Second * 5,
		discovery:  builder.NewConsulDiscovery(endpoint),
	}
	// 初始化
	r.event = make(chan []*register.ServiceInstance, 1)

	_ = r.initNodes(context.Background(), serviceName)
	r.serviceName = serviceName

	return r
}

func (r *replicant) initNodes(_ context.Context, serviceName string) register.Watcher {
	var (
		disErr error
		w      register.Watcher
	)
	nctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	done := make(chan struct{}, 1)
	go func() {
		w, disErr = r.discovery.Watch(nctx, serviceName)
		var (
			ins []*register.ServiceInstance
		)
		if ins, disErr = w.Next(); disErr == nil {
			r.event <- ins
		}
		close(done)
	}()
	var (
		err error
	)
	select {
	case <-done:
		err = disErr
	case <-time.After(10 * time.Second):
		err = errors.New("discovery create watcher overtime")
	}
	if err != nil {
		cancel()
		r.discovery.Close()
		return w
	}
	// 启动监控
	go r.OnlineProc(nctx)
	go func() {
		timeTick := time.NewTicker(5 * time.Second)
		defer timeTick.Stop()
		for {
			select {
			case <-nctx.Done():
				return
			case <-timeTick.C:
				// 监控服务最新状态 随时进行负载均衡
				if w, err := r.discovery.Watch(nctx, serviceName); err == nil {
					if ins, err := w.Next(); err == nil && len(ins) > 0 {
						r.event <- ins
					}
				}
			}
		}
	}()
	return w
}

func (r *replicant) LoadBalancerUpdate(ins []*register.ServiceInstance) {
	var (
		totalConns  int64
		totalWeight int64
		nodes       = make(map[string]*weighted, len(ins))
	)
	if len(ins) == 0 || float32(len(ins))/float32(len(r.allIns)) < 0.5 {
		return
	}
	r.nodesMutex.Lock()
	defer r.nodesMutex.Unlock()

	for _, in := range ins {
		if old, ok := r.nodes[in.Name]; ok && old.updated == in.LastTs {
			nodes[in.Name] = old
			totalConns += old.currentConns
			totalWeight += old.currentWeight
		} else {
			meta := in.Metadata
			weight, err := strconv.ParseInt(meta["weight"], 10, 32)
			if err != nil {
				continue
			}
			conns, err := strconv.ParseInt(meta["connCount"], 10, 32)
			if err != nil {
				continue
			}
			var (
				addrs string
			)
			if len(in.Endpoints) > 0 {
				for _, addr := range in.Endpoints {
					uri, _ := url.Parse(addr)
					if uri.Scheme == "grpc" {
						addrs = uri.Host
						break
					}
				}
			}
			nodes[in.Name] = &weighted{
				updated:      in.LastTs,
				fixedWeight:  weight,
				currentConns: conns,
				addres:       addrs,
				hostName:     in.Name,
			}
			totalConns += conns
			totalWeight += weight
		}
	}
	r.nodes = nodes
	r.totalConns = totalConns
	r.totalWeight = totalWeight
}

func (r *replicant) OnlineProc(ctx context.Context) {
	go func() {
		for si := range r.event {
			var (
				totalConns int64
				totalIPs   int64
			)
			r.allIns = make([]*register.ServiceInstance, 0, len(si))
			// 遍历
			for _, instance := range si {
				if len(instance.Metadata) == 0 {
					continue
				}
				conns, err := strconv.ParseInt(instance.Metadata["connCount"], 10, 32)
				if err != nil || conns <= 0 {
					continue
				}
				ips, err := strconv.ParseInt(instance.Metadata["ipCount"], 10, 32)
				if err != nil || ips <= 0 {
					continue
				}
				totalConns += conns
				totalIPs += ips
				r.allIns = append(r.allIns, instance)
			}
			r.totalConns = totalConns
			r.totalIPs = totalIPs
			r.LoadBalancerUpdate(r.allIns)
		}
	}()

}

func (r *replicant) weightedNodes(region string, regionWeight float64) (nodes []*weighted) {
	nodes = make([]*weighted, 0, len(r.nodes))
	for _, n := range r.nodes {
		var (
			gainWeight = float64(1.0)
		)
		if n.region == region {
			gainWeight *= regionWeight
		}
		n.calculateWeight(r.totalWeight, r.totalConns, gainWeight)
		nodes = append(nodes, n)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].currentWeight > nodes[j].currentWeight
	})
	if len(nodes) > 0 {
		nodes[0].chosen()
		r.totalConns++
	}
	return
}

func (r *replicant) NodeAddrs(region, domain string, regionWeight float64) (domains, addrs []string) {
	r.nodesMutex.Lock()
	defer r.nodesMutex.Unlock()
	node := r.weightedNodes(region, regionWeight)
	for i, n := range node {
		if i == 5 {
			break
		}
		domains = append(domains, n.addres)
		addrs = append(addrs, n.addres)
	}
	return
}
