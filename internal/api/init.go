package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/yanglunara/discovery/watcher/consul"
	"github.com/yanglunara/im/pkg/common/prometheus/ginprometheus"
	"github.com/yanglunara/im/pkg/common/prometheus/metrics"
	"github.com/yunbaifan/pkg/utils/errs"
)

func Start(ctx context.Context, index int, conf *Conf) error {
	fun := func(array []int, index int) (int, error) {
		if index < 0 || index >= len(array) {
			return 0, fmt.Errorf("index out of range")
		}
		return array[index], nil
	}
	port, err := fun(conf.Rpc.TCP.Ports, index)
	if err != nil {
		return err
	}
	pPort, err := fun(conf.Rpc.Prometheus.Ports, index)
	if err != nil {
		return err
	}
	// consul
	cli, err := api.NewClient(&api.Config{
		Address: conf.Consul.Address,
	})
	if err != nil {
		return err
	}
	var (
		Done = make(chan struct{}, 1)
		Err  error
	)
	router := startGinRouter(consul.NewRegistry(cli), conf)
	// 是否开启普罗米修斯
	if conf.Rpc.Prometheus.Enable {
		go func() {
			p := ginprometheus.NewPrometheus("app", metrics.GetImApiGinMetrics())
			p.SetListenAddress(fmt.Sprintf(":%d", pPort))
			if err = p.Use(router); err != nil && err != http.ErrServerClosed {
				Err = errs.Warp(err, fmt.Sprintf("prometheus start err: %d", pPort))
				Done <- struct{}{}
			}
		}()
	}
	address := net.JoinHostPort(
		conf.Rpc.TCP.ListenIP,
		strconv.Itoa(port),
	)
	server := http.Server{
		Addr:    address,
		Handler: router,
	}
	go func() {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Err = errs.Warp(err, fmt.Sprintf("server start err: %s", address))
			Done <- struct{}{}
		}
	}()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	cxt, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	select {
	case <-signalChan:
		Exit()
		if err = server.Shutdown(cxt); err != nil {
			return errs.Warp(err, "server shutdown err")
		}
	case <-Done:
		close(Done)
		return Err
	}
	return nil
}

func startGinRouter(reg *consul.Registry, conf *Conf) *gin.Engine {

	return nil
}

func Exit() {
	progName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Warning %s receive process terminal SIGTERM exit 0\n", progName)
}
