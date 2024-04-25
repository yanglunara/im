package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/yanglunara/discovery/register"
	bgrpc "github.com/yanglunara/discovery/transport/grpc"
	"github.com/yanglunara/discovery/watcher/consul"
	conf "github.com/yanglunara/im/internal/conf/logic"
	"github.com/yanglunara/im/internal/logic"
	"github.com/yanglunara/im/internal/model"
	"github.com/yunbaifan/pkg/logger"
	"go.uber.org/zap"
)

var (
	// Version is the version of the application
	Version     string = "1.0.0"
	ServerName  string = "im.logic"
	flagconf    string
	ServerID, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", fmt.Sprintf("config path, eg: -conf %s.yaml", ServerName))
}

func main() {
	flag.Parse()
	//ServerID = uuid.New().String()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := conf.Init(flagconf); err != nil {
		panic(err)
	}
	// init logger
	_ = logger.InitLogger(conf.Conf.Logger)

	ctx := context.Background()
	logicRPC := logic.NewLogic(conf.Conf)

	gracefulShutdown(
		ctx,
		logic.NewLogicService(
			&conf.Conf.GrpcServer,
			logicRPC,
		), logicRPC)
}

func gracefulShutdown(ctx context.Context, rpcSrv *bgrpc.Service, logicRPC model.LogicConnService) {
	signalCtx, stop := signal.NotifyContext(
		ctx,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	//注册服务
	cli, err := api.NewClient(&api.Config{
		Address: conf.Conf.Consul.Address,
	})
	if err != nil {
		panic(err)
	}
	uri, _ := rpcSrv.Endpoint()
	//服务注册
	instance := &register.ServiceInstance{
		ID:       ServerID,
		Name:     ServerName,
		Version:  Version,
		LastTs:   time.Now().Unix(),
		Metadata: make(map[string]string),
		Endpoints: []string{
			fmt.Sprintf("%s://%s", uri.Scheme, uri.Host),
		},
	}
	// 服务注册与发现
	discovery := consul.NewRegistry(cli)
	go func() {
		// 注册服务
		if err := discovery.Register(ctx, instance); err != nil {
			logger.Logger.Fatal("Register: %s\n", zap.Error(err))
		}
		logger.Logger.Info("Register Service Success")
	}()

	go func() {
		logger.Logger.Info("Grpc Service Success")
		if err := rpcSrv.Start(ctx); err != nil {
			logger.Logger.Fatal("rpc server start :%s \n", zap.Error(err))
		}
	}()
	go func() {
		logicRPC.SetReplicant(ctx, "im.gateway")
	}()
	<-signalCtx.Done()
	stop()
	//设置超时时间

	ctx, cancel := context.WithTimeout(context.Background(),
		5*time.Second,
	)
	defer func() {
		_ = rpcSrv.Stop(ctx)
		cancel()
		stop()
		//关闭服务
		_ = discovery.Close()
		// 注销当前在consul 保存服务数据
		_ = discovery.Deregister(ctx, instance)
	}()

	logger.Logger.Info("Shutdown Server ...")

}
