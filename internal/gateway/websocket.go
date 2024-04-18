package gateway

import (
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	pb "github.com/yanglunara/im/api/protocol"
	conf "github.com/yanglunara/im/internal/conf/gateway"
)

type WebSocketInter interface {
	Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, string, error)
}

type imSocket struct {
	keyGUID  []byte
	conf     *conf.Config
	serverID string
	gs       *GatewayServer
}

func NewWebSocket(conf *conf.Config, gs *GatewayServer, serverID string) WebSocketInter {
	return &imSocket{
		keyGUID:  []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B22"),
		conf:     conf,
		serverID: serverID,
		gs:       gs,
	}
}

func (i *imSocket) acceptKey(challengeKey string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(challengeKey))
	_, _ = h.Write(i.keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (i *imSocket) Upgrade(w http.ResponseWriter, r *http.Request) (conn *websocket.Conn, k string, err error) {
	ws := new(websocket.Upgrader)
	ws.CheckOrigin = func(req *http.Request) bool {
		check := true
		switch {
		case req.Method != "GET":
			check = false
		case req.Header.Get("Sec-Websocket-Version") != "13":
			check = false
		case !strings.Contains(strings.ToLower(req.Header.Get("Connection")), "upgrade"):
			check = false
		}
		if challengeKey := req.Header.Get("Sec-Websocket-Key"); challengeKey == "" {
			check = false
		} else {
			req.Header.Set("Sec-Websocket-Key", i.acceptKey(challengeKey))
		}
		return check
	}
	conn = new(websocket.Conn)
	if conn, err = ws.Upgrade(w, r, nil); err != nil {
		return nil, "", err
	}
	var (
		p *pb.Proto
	)
	ch := NewChannel(
		i.conf.Protocol.CliProto,
		i.conf.Protocol.SvrProto,
	)
	var (
		a *AuthWebSocket
	)
	// 注册链接
	if p, err = ch.CliProto.Put(); err == nil {
		if a, err = NewClientManager(i.serverID).Start(r.Context(), p, r.Header.Get("Cookie"), conn); err == nil {
			ch.Watch(a.Accepts...)
			if b := i.gs.Bucket(ch.Key); b != nil {
				err = b.Put(a.RID, ch)
			}
		}
	}
	if err != nil {
		conn.Close()
		return nil, "", err
	}
	return conn, r.Header.Get("Cookie"), err
}
