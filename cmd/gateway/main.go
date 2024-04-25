package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/yanglunara/discovery/register"
	bgrpc "github.com/yanglunara/discovery/transport/grpc"
	"github.com/yanglunara/discovery/watcher/consul"
	conf "github.com/yanglunara/im/internal/conf/gateway"
	"github.com/yanglunara/im/internal/gateway"
	"github.com/yunbaifan/pkg/logger"
)

var (
	// Version is the version of the application
	Version     string = "1.0.0"
	ServerName  string = "im.gateway"
	flagconf    string
	ServerID, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", fmt.Sprintf("config path, eg: -conf %s.yaml", ServerName))
}

func main() {
	flag.Parse()
	ServerID = uuid.New().String()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := conf.Init(flagconf); err != nil {
		panic(err)
	}
	// init logger
	_ = logger.InitLogger(conf.Conf.Logger)
	ctx := context.Background()
	// 退出流程
	gracefulShutdown(
		ctx,
		gateway.NewGatewayGrpcServer(
			&conf.Conf.GrpcServer,
			gateway.NewGatewayService(),
		),
		gateway.NewServer(ctx, conf.Conf, ServerID),
	)
}

func gracefulShutdown(ctx context.Context, rpcSrv *bgrpc.Service, ws *gateway.Websocket) {
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
	var (
		localIP string
	)
	if len(uri.Host) > 0 {
		localIP = strings.Split(uri.Host, ":")[0]
	}
	addrs := []string{localIP}
	instance := &register.ServiceInstance{
		ID:      ServerID,
		Name:    ServerName,
		Version: Version,
		Metadata: map[string]string{
			"weight": strconv.Itoa(int(conf.Conf.GlobalEnv.Weight)),
			"addrs":  strings.Join(addrs, ","),
		},
		Endpoints: []string{
			fmt.Sprintf("%s://%s", uri.Scheme, uri.Host),
		},
	}
	// 服务注册与发现
	discovery := consul.NewRegistry(cli)
	//http 服务
	gateway.StartHttpSever(ctx)
	// 服务注册
	go gateway.RegisterServer(ctx, discovery, instance)
	//websocket
	go gateway.StartWebsocket(ws)
	// rpc 服务
	go gateway.StartGrpcServer(ctx, rpcSrv)
	// 定时上报服务
	go gateway.ReportIpAndConnCount(ctx, discovery, ws, ServerName, ServerID)

	<-signalCtx.Done()
	stop()
	//设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(),
		5*time.Second,
	)
	defer func() {
		_ = rpcSrv.Stop(ctx)
		//关闭服务
		_ = discovery.Close()
		// 注销当前在consul 保存服务数据
		_ = discovery.Deregister(ctx, instance)
		// 关闭channel
		_ = ws.Close()
		cancel()
		stop()
	}()
	logger.Logger.Info("Shutdown Server ...")

}
