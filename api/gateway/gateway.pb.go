// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.18.0
// source: api/gateway/gateway.proto

package gateway

import (
	protocol "github.com/yanglunara/im/api/protocol"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PushMessageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys    []string        `protobuf:"bytes,1,rep,name=keys,proto3" json:"keys,omitempty"`
	Proto   *protocol.Proto `protobuf:"bytes,2,opt,name=proto,proto3" json:"proto,omitempty"`
	ProtoOp int32           `protobuf:"varint,3,opt,name=protoOp,proto3" json:"protoOp,omitempty"`
}

func (x *PushMessageReq) Reset() {
	*x = PushMessageReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_gateway_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMessageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMessageReq) ProtoMessage() {}

func (x *PushMessageReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_gateway_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMessageReq.ProtoReflect.Descriptor instead.
func (*PushMessageReq) Descriptor() ([]byte, []int) {
	return file_api_gateway_gateway_proto_rawDescGZIP(), []int{0}
}

func (x *PushMessageReq) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *PushMessageReq) GetProto() *protocol.Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

func (x *PushMessageReq) GetProtoOp() int32 {
	if x != nil {
		return x.ProtoOp
	}
	return 0
}

type PushMessageResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushMessageResp) Reset() {
	*x = PushMessageResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_gateway_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMessageResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMessageResp) ProtoMessage() {}

func (x *PushMessageResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_gateway_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMessageResp.ProtoReflect.Descriptor instead.
func (*PushMessageResp) Descriptor() ([]byte, []int) {
	return file_api_gateway_gateway_proto_rawDescGZIP(), []int{1}
}

type BroadcastReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProtoOp int32           `protobuf:"varint,1,opt,name=protoOp,proto3" json:"protoOp,omitempty"`
	Proto   *protocol.Proto `protobuf:"bytes,2,opt,name=proto,proto3" json:"proto,omitempty"`
	Speed   int32           `protobuf:"varint,3,opt,name=speed,proto3" json:"speed,omitempty"`
}

func (x *BroadcastReq) Reset() {
	*x = BroadcastReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_gateway_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastReq) ProtoMessage() {}

func (x *BroadcastReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_gateway_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastReq.ProtoReflect.Descriptor instead.
func (*BroadcastReq) Descriptor() ([]byte, []int) {
	return file_api_gateway_gateway_proto_rawDescGZIP(), []int{2}
}

func (x *BroadcastReq) GetProtoOp() int32 {
	if x != nil {
		return x.ProtoOp
	}
	return 0
}

func (x *BroadcastReq) GetProto() *protocol.Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

func (x *BroadcastReq) GetSpeed() int32 {
	if x != nil {
		return x.Speed
	}
	return 0
}

type BroadcastResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BroadcastResp) Reset() {
	*x = BroadcastResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_gateway_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastResp) ProtoMessage() {}

func (x *BroadcastResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_gateway_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastResp.ProtoReflect.Descriptor instead.
func (*BroadcastResp) Descriptor() ([]byte, []int) {
	return file_api_gateway_gateway_proto_rawDescGZIP(), []int{3}
}

type BroadcastRoomReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomID string          `protobuf:"bytes,1,opt,name=roomID,proto3" json:"roomID,omitempty"`
	Proto  *protocol.Proto `protobuf:"bytes,2,opt,name=proto,proto3" json:"proto,omitempty"`
}

func (x *BroadcastRoomReq) Reset() {
	*x = BroadcastRoomReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_gateway_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastRoomReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastRoomReq) ProtoMessage() {}

func (x *BroadcastRoomReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_gateway_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastRoomReq.ProtoReflect.Descriptor instead.
func (*BroadcastRoomReq) Descriptor() ([]byte, []int) {
	return file_api_gateway_gateway_proto_rawDescGZIP(), []int{4}
}

func (x *BroadcastRoomReq) GetRoomID() string {
	if x != nil {
		return x.RoomID
	}
	return ""
}

func (x *BroadcastRoomReq) GetProto() *protocol.Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

type BroadcastRoomResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BroadcastRoomResp) Reset() {
	*x = BroadcastRoomResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_gateway_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastRoomResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastRoomResp) ProtoMessage() {}

func (x *BroadcastRoomResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_gateway_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastRoomResp.ProtoReflect.Descriptor instead.
func (*BroadcastRoomResp) Descriptor() ([]byte, []int) {
	return file_api_gateway_gateway_proto_rawDescGZIP(), []int{5}
}

type RoomsReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RoomsReq) Reset() {
	*x = RoomsReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_gateway_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoomsReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoomsReq) ProtoMessage() {}

func (x *RoomsReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_gateway_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoomsReq.ProtoReflect.Descriptor instead.
func (*RoomsReq) Descriptor() ([]byte, []int) {
	return file_api_gateway_gateway_proto_rawDescGZIP(), []int{6}
}

type RoomsResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rooms map[string]bool `protobuf:"bytes,1,rep,name=rooms,proto3" json:"rooms,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *RoomsResp) Reset() {
	*x = RoomsResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_gateway_gateway_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoomsResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoomsResp) ProtoMessage() {}

func (x *RoomsResp) ProtoReflect() protoreflect.Message {
	mi := &file_api_gateway_gateway_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoomsResp.ProtoReflect.Descriptor instead.
func (*RoomsResp) Descriptor() ([]byte, []int) {
	return file_api_gateway_gateway_proto_rawDescGZIP(), []int{7}
}

func (x *RoomsResp) GetRooms() map[string]bool {
	if x != nil {
		return x.Rooms
	}
	return nil
}

var File_api_gateway_gateway_proto protoreflect.FileDescriptor

var file_api_gateway_gateway_proto_rawDesc = []byte{
	0x0a, 0x19, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x69, 0x6d, 0x2e,
	0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x1a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x68, 0x0a, 0x0e, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x28, 0x0a, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x69, 0x6d, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x22, 0x11,
	0x0a, 0x0f, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x22, 0x68, 0x0a, 0x0c, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x4f, 0x70, 0x12, 0x28, 0x0a, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x69, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x73, 0x70, 0x65, 0x65, 0x64, 0x22, 0x0f, 0x0a, 0x0d, 0x42,
	0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x22, 0x54, 0x0a, 0x10,
	0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x71,
	0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x44, 0x12, 0x28, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x69, 0x6d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x13, 0x0a, 0x11, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52,
	0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x22, 0x0a, 0x0a, 0x08, 0x52, 0x6f, 0x6f, 0x6d, 0x73,
	0x52, 0x65, 0x71, 0x22, 0x7d, 0x0a, 0x09, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x36, 0x0a, 0x05, 0x72, 0x6f, 0x6f, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x69, 0x6d, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x52, 0x6f, 0x6f,
	0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x05, 0x72, 0x6f, 0x6f, 0x6d, 0x73, 0x1a, 0x38, 0x0a, 0x0a, 0x52, 0x6f, 0x6f, 0x6d,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x32, 0x97, 0x02, 0x0a, 0x07, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x12, 0x40,
	0x0a, 0x09, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x18, 0x2e, 0x69, 0x6d,
	0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x19, 0x2e, 0x69, 0x6d, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x12, 0x46, 0x0a, 0x0b, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x1a, 0x2e, 0x69, 0x6d, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x50, 0x75, 0x73,
	0x68, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x1b, 0x2e, 0x69, 0x6d,
	0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x4c, 0x0a, 0x0d, 0x42, 0x72, 0x6f, 0x61,
	0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x1c, 0x2e, 0x69, 0x6d, 0x2e, 0x67,
	0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x1d, 0x2e, 0x69, 0x6d, 0x2e, 0x67, 0x61, 0x74,
	0x65, 0x77, 0x61, 0x79, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x6f,
	0x6f, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x12, 0x34, 0x0a, 0x05, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x12,
	0x14, 0x2e, 0x69, 0x6d, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2e, 0x52, 0x6f, 0x6f,
	0x6d, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x15, 0x2e, 0x69, 0x6d, 0x2e, 0x67, 0x61, 0x74, 0x65, 0x77,
	0x61, 0x79, 0x2e, 0x52, 0x6f, 0x6f, 0x6d, 0x73, 0x52, 0x65, 0x73, 0x70, 0x42, 0x17, 0x5a, 0x15,
	0x2e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x3b, 0x67, 0x61,
	0x74, 0x65, 0x77, 0x61, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_gateway_gateway_proto_rawDescOnce sync.Once
	file_api_gateway_gateway_proto_rawDescData = file_api_gateway_gateway_proto_rawDesc
)

func file_api_gateway_gateway_proto_rawDescGZIP() []byte {
	file_api_gateway_gateway_proto_rawDescOnce.Do(func() {
		file_api_gateway_gateway_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_gateway_gateway_proto_rawDescData)
	})
	return file_api_gateway_gateway_proto_rawDescData
}

var file_api_gateway_gateway_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_api_gateway_gateway_proto_goTypes = []interface{}{
	(*PushMessageReq)(nil),    // 0: im.gateway.PushMessageReq
	(*PushMessageResp)(nil),   // 1: im.gateway.PushMessageResp
	(*BroadcastReq)(nil),      // 2: im.gateway.BroadcastReq
	(*BroadcastResp)(nil),     // 3: im.gateway.BroadcastResp
	(*BroadcastRoomReq)(nil),  // 4: im.gateway.BroadcastRoomReq
	(*BroadcastRoomResp)(nil), // 5: im.gateway.BroadcastRoomResp
	(*RoomsReq)(nil),          // 6: im.gateway.RoomsReq
	(*RoomsResp)(nil),         // 7: im.gateway.RoomsResp
	nil,                       // 8: im.gateway.RoomsResp.RoomsEntry
	(*protocol.Proto)(nil),    // 9: im.protocol.Proto
}
var file_api_gateway_gateway_proto_depIdxs = []int32{
	9, // 0: im.gateway.PushMessageReq.proto:type_name -> im.protocol.Proto
	9, // 1: im.gateway.BroadcastReq.proto:type_name -> im.protocol.Proto
	9, // 2: im.gateway.BroadcastRoomReq.proto:type_name -> im.protocol.Proto
	8, // 3: im.gateway.RoomsResp.rooms:type_name -> im.gateway.RoomsResp.RoomsEntry
	2, // 4: im.gateway.Gateway.Broadcast:input_type -> im.gateway.BroadcastReq
	0, // 5: im.gateway.Gateway.PushMessage:input_type -> im.gateway.PushMessageReq
	4, // 6: im.gateway.Gateway.BroadcastRoom:input_type -> im.gateway.BroadcastRoomReq
	6, // 7: im.gateway.Gateway.Rooms:input_type -> im.gateway.RoomsReq
	3, // 8: im.gateway.Gateway.Broadcast:output_type -> im.gateway.BroadcastResp
	1, // 9: im.gateway.Gateway.PushMessage:output_type -> im.gateway.PushMessageResp
	5, // 10: im.gateway.Gateway.BroadcastRoom:output_type -> im.gateway.BroadcastRoomResp
	7, // 11: im.gateway.Gateway.Rooms:output_type -> im.gateway.RoomsResp
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_gateway_gateway_proto_init() }
func file_api_gateway_gateway_proto_init() {
	if File_api_gateway_gateway_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_gateway_gateway_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMessageReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_gateway_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMessageResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_gateway_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_gateway_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_gateway_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastRoomReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_gateway_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastRoomResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_gateway_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoomsReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_gateway_gateway_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoomsResp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_gateway_gateway_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_gateway_gateway_proto_goTypes,
		DependencyIndexes: file_api_gateway_gateway_proto_depIdxs,
		MessageInfos:      file_api_gateway_gateway_proto_msgTypes,
	}.Build()
	File_api_gateway_gateway_proto = out.File
	file_api_gateway_gateway_proto_rawDesc = nil
	file_api_gateway_gateway_proto_goTypes = nil
	file_api_gateway_gateway_proto_depIdxs = nil
}
