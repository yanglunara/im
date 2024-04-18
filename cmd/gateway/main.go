package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"github.com/yanglunara/discovery/register"
	"github.com/yanglunara/discovery/watcher/consul"
	conf "github.com/yanglunara/im/internal/conf/gateway"
	"github.com/yanglunara/im/internal/gateway"
	"github.com/yunbaifan/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// Version is the version of the application
	Version            = "0.0.1"
	ServerName  string = "im.gateway"
	flagconf    string
	ServerID, _ = os.Hostname()
	port        = 8080
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", fmt.Sprintf("config path, eg: -conf %s.yaml", ServerName))
}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	//注册服务
	cli, err := api.NewClient(&api.Config{
		Address: conf.Conf.Consul.Address,
	})
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	//服务注册
	instance := &register.ServiceInstance{
		ID:        ServerID,
		Name:      ServerName,
		LastTs:    time.Now().Unix(),
		Metadata:  make(map[string]string),
		Endpoints: []string{},
	}
	ginhttp := gin.Default()
	var (
		discovery = consul.NewRegistry(cli)
	)
	// 服务注册
	_ = discovery.Register(ctx, instance)
	// gin http websocket 服务
	gracefulShutdown(
		ctx,
		gateway.NewPushServer(
			conf.Conf.GRPCServer,
			gateway.NewPushService(),
		), ginhttp, discovery, instance)
}

func gracefulShutdown(ctx context.Context, rpcSrv *grpc.Server, ginHttp *gin.Engine, discovery *consul.Registry, instance *register.ServiceInstance) {
	signalCtx, stop := signal.NotifyContext(
		ctx,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	src := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: ginHttp,
	}
	go func() {
		if err := src.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("listen: %s\n", zap.Error(err))
		}
	}()
	<-signalCtx.Done()
	stop()
	//设置超时时间
	ctx, cancel := context.WithTimeout(signalCtx,
		time.Duration(conf.Conf.GracefulTimeout)*time.Second,
	)
	defer func() {
		rpcSrv.GracefulStop()
		cancel()
		stop()
		//关闭服务
		discovery.Close()
		// 注销当前在consul 保存服务数据
		_ = discovery.Deregister(ctx, instance)

	}()

}
