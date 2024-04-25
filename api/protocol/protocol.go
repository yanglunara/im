package protocol

import (
	"errors"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

const (
	MaxBodySize = int32(1 << 13)

	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = 1 << 4
	_maxPackSize   = MaxBodySize + int32(16)
	// offset
	_packOffset   = 0
	_headerOffset = 1 << 2
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
)

var (
	ErrProtoPackLen   = errors.New("default server codec pack length error")
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

func (p *Proto) ReadWebsocket(ws *websocket.Conn) (err error) {
	var (
		buf []byte
	)
	if _, buf, err = ws.ReadMessage(); err != nil {
		return
	}
	if err = proto.Unmarshal(buf, p); err != nil {
		return
	}
	return
}
