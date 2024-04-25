package gateway

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yanglunara/discovery/register"
	bgrpc "github.com/yanglunara/discovery/transport/grpc"
	"github.com/yanglunara/discovery/watcher/consul"
	"github.com/yunbaifan/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

func StartWebsocket(ws *Websocket) {
	for i := 0; i < runtime.NumCPU(); i++ {
		go ws.Start()
	}
}

func ReportIpAndConnCount(ctx context.Context, discovery *consul.Registry, ws *Websocket, serverName, serverID string) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			var (
				err   error
				conns int
				ips   = make(map[string]struct{})
			)
			for _, bucket := range ws.Buckets() {
				for ip := range bucket.IpCount() {
					ips[ip] = struct{}{}
				}
				conns += bucket.CountChannel()
			}
			// 定时上报 先get 服务的信息
			olds, err := discovery.GetService(ctx, serverName)
			if err != nil {
				logger.Logger.Error("get service info fail", zap.Error(err))
				continue
			}
			for _, old := range olds {
				if old.ID == serverID {
					old.Metadata["connCount"] = strconv.Itoa(conns)
					old.Metadata["ipCount"] = strconv.Itoa(len(ips))
					old.LastTs = time.Now().Unix()
					_ = discovery.Register(ctx, old)
					break
				}
			}
		}
	}
}

func StartGrpcServer(ctx context.Context, srv *bgrpc.Service) {
	logger.Logger.Info("start Grpc server ...")
	if err := srv.Start(ctx); err != nil {
		logger.Logger.Fatal("rpc server start :%s \n", zap.Error(err))
	}
}

func RegisterServer(ctx context.Context, discovery *consul.Registry, instance *register.ServiceInstance) {
	if err := discovery.Register(ctx, instance); err != nil {
		logger.Logger.Fatal("Register: %s\n", zap.Error(err))
	}
	logger.Logger.Info("Register Service Success")
}

func StartHttpSever(_ context.Context) {
	ginHttp := gin.Default()

	ginRouter(ginHttp)

	go func() {
		if err := (&http.Server{
			Addr:    fmt.Sprintf(":%d", 8080),
			Handler: ginHttp,
		}).ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("listen: %s\n", zap.Error(err))
		}
		logger.Logger.Info("Http Service Success")
	}()
}

func ginRouter(engine *gin.Engine) {
	engine.GET("/ws", wsHandler)
}

func wsHandler(c *gin.Context) {
	if err := WebSocket.Upgrade(c.Writer, c.Request); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "websocket upgrade error",
			"data":    nil,
		})
		return
	}
}
