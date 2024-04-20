package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	pt "github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/yanglunara/discovery/register"
	bgrpc "github.com/yanglunara/discovery/transport/grpc"
	"github.com/yanglunara/discovery/watcher/consul"
	"github.com/yanglunara/im/api/logic"
	"github.com/yanglunara/im/api/protocol"
	conf "github.com/yanglunara/im/internal/conf/gateway"
	"github.com/yanglunara/im/internal/gateway"
	"github.com/yunbaifan/pkg/logger"
	"go.uber.org/zap"
)

var (
	// Version is the version of the application
	Version     string
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
	// gin http 服务
	ginhttp := gin.Default()

	// 注册服务
	registerSrv(ctx)
	// 启动 gateway 房间服务
	routers(ginhttp)
	// gin http websocket grpc 服务
	gracefulShutdown(
		ctx,
		int32(conf.Conf.HTTPServer.Port),
		gateway.NewGatewayGrpcServer(
			&conf.Conf.GrpcServer,
			gateway.NewGatewayService(),
		), ginhttp)
}

func registerSrv(ctx context.Context) {
	// websocket 服务
	socket := gateway.NewWebSocket(conf.Conf,
		gateway.NewServer(ctx,
			conf.Conf,
			ServerID,
		),
		ServerID)
	// grpc client
	_ = gateway.NewLogicGrpc(ctx, conf.Conf)
	// 启动管理者
	go socket.DispatchWebsocket()
}

func routers(engine *gin.Engine) {
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "ok",
			"data":    nil,
		})
	})
	engine.GET("/ws", wsHandler)
}

const (
	minServerHeartbeat = time.Minute * 10
	maxServerHeartbeat = time.Minute * 30
)

func wsHandler(c *gin.Context) {
	// 创建websocket 连接 进行auth 认证
	ch, _, err := gateway.WebSocket.Upgrade(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "websocket upgrade error",
			"data":    nil,
		})
		return
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				_, _ = gateway.LogicGrpcClient.Disconnect(c.Request.Context(), &logic.DisconnectReq{
					Mid:    ch.Mid,
					Key:    ch.Key,
					Server: ServerID,
				})
				ch.Close()
			}
		}()
		var (
			p               *gateway.ProtoRing
			serverHeartbeat = (minServerHeartbeat + time.Duration(rand.Int63n(int64(maxServerHeartbeat-minServerHeartbeat))))
		)
		for {
			if p, err = ch.CliProto.Put(); err != nil {
				break
			}
			var (
				message []byte
				err     error
			)
			if _, message, err = p.Conn.ReadMessage(); err != nil {
				logger.Logger.Error("read message fail ", zap.String("zddr", p.Conn.RemoteAddr().String()), zap.Error(err))
				break
			}
			if err = pt.Unmarshal(message, p.Proto); err != nil {
				break
			}
			if p.Op == int32(protocol.Operation_Heartbeat) {
				p.Op = int32(protocol.Operation_HeartbeatResp)
				p.Body = nil
				if now := time.Now(); now.Sub(ch.LastHB) > serverHeartbeat {
					if _, err := gateway.LogicGrpcClient.Heartbeat(context.Background(), &logic.HeartbeatReq{
						Mid:    ch.Mid,
						Key:    ch.Key,
						Server: ServerID,
					}); err == nil {
						ch.LastHB = now
					}
				}
			} else {
				// operate
			}
			ch.CliProto.IncrementWritePointer()
			ch.Signal()
		}
	}()

}

func gracefulShutdown(ctx context.Context, port int32, rpcSrv *bgrpc.Service, ginHttp *gin.Engine) {
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
	instance := &register.ServiceInstance{
		ID:       ServerID,
		Name:     ServerName,
		LastTs:   time.Now().Unix(),
		Metadata: map[string]string{},
		Endpoints: []string{
			fmt.Sprintf("%s://%s", uri.Scheme, uri.Host),
			fmt.Sprintf("%s://%s:%d", "http", localIP, conf.Conf.HTTPServer.Port),
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
	//http 服务
	go func() {
		if err := src.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("listen: %s\n", zap.Error(err))
		}
		logger.Logger.Info("Http Service Success")
	}()

	go func() {
		if err := rpcSrv.Start(ctx); err != nil {
			logger.Logger.Fatal("rpc server start :%s \n", zap.Error(err))
		}
		logger.Logger.Info("Grpc Service Success")
	}()
	<-signalCtx.Done()
	stop()
	//设置超时时间

	ctx, cancel := context.WithTimeout(context.Background(),
		5*time.Second,
	)
	defer func() {
		rpcSrv.Stop(ctx)
		cancel()
		stop()
		//关闭服务
		discovery.Close()
		// 注销当前在consul 保存服务数据
		_ = discovery.Deregister(ctx, instance)
	}()

	logger.Logger.Info("Shutdown Server ...")

}
