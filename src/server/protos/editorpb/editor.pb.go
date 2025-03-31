// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v6.30.0
// source: protos/editor.proto

package editorpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ShareReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DocName       string                 `protobuf:"bytes,1,opt,name=doc_name,json=docName,proto3" json:"doc_name,omitempty"`
	Doc           []byte                 `protobuf:"bytes,2,opt,name=doc,proto3" json:"doc,omitempty"`
	UserId        string                 `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShareReq) Reset() {
	*x = ShareReq{}
	mi := &file_protos_editor_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShareReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShareReq) ProtoMessage() {}

func (x *ShareReq) ProtoReflect() protoreflect.Message {
	mi := &file_protos_editor_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShareReq.ProtoReflect.Descriptor instead.
func (*ShareReq) Descriptor() ([]byte, []int) {
	return file_protos_editor_proto_rawDescGZIP(), []int{0}
}

func (x *ShareReq) GetDocName() string {
	if x != nil {
		return x.DocName
	}
	return ""
}

func (x *ShareReq) GetDoc() []byte {
	if x != nil {
		return x.Doc
	}
	return nil
}

func (x *ShareReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ShareReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DocId         string                 `protobuf:"bytes,1,opt,name=doc_id,json=docId,proto3" json:"doc_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ShareReply) Reset() {
	*x = ShareReply{}
	mi := &file_protos_editor_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ShareReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShareReply) ProtoMessage() {}

func (x *ShareReply) ProtoReflect() protoreflect.Message {
	mi := &file_protos_editor_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShareReply.ProtoReflect.Descriptor instead.
func (*ShareReply) Descriptor() ([]byte, []int) {
	return file_protos_editor_proto_rawDescGZIP(), []int{1}
}

func (x *ShareReply) GetDocId() string {
	if x != nil {
		return x.DocId
	}
	return ""
}

type DeleteReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DocId         string                 `protobuf:"bytes,1,opt,name=doc_id,json=docId,proto3" json:"doc_id,omitempty"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteReq) Reset() {
	*x = DeleteReq{}
	mi := &file_protos_editor_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteReq) ProtoMessage() {}

func (x *DeleteReq) ProtoReflect() protoreflect.Message {
	mi := &file_protos_editor_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteReq.ProtoReflect.Descriptor instead.
func (*DeleteReq) Descriptor() ([]byte, []int) {
	return file_protos_editor_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteReq) GetDocId() string {
	if x != nil {
		return x.DocId
	}
	return ""
}

func (x *DeleteReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type DeleteReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteReply) Reset() {
	*x = DeleteReply{}
	mi := &file_protos_editor_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteReply) ProtoMessage() {}

func (x *DeleteReply) ProtoReflect() protoreflect.Message {
	mi := &file_protos_editor_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteReply.ProtoReflect.Descriptor instead.
func (*DeleteReply) Descriptor() ([]byte, []int) {
	return file_protos_editor_proto_rawDescGZIP(), []int{3}
}

type EditReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DocId         string                 `protobuf:"bytes,1,opt,name=doc_id,json=docId,proto3" json:"doc_id,omitempty"`
	Rev           int32                  `protobuf:"varint,2,opt,name=rev,proto3" json:"rev,omitempty"`
	Ops           []*Op                  `protobuf:"bytes,3,rep,name=ops,proto3" json:"ops,omitempty"`
	UserId        string                 `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Title         string                 `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EditReq) Reset() {
	*x = EditReq{}
	mi := &file_protos_editor_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EditReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EditReq) ProtoMessage() {}

func (x *EditReq) ProtoReflect() protoreflect.Message {
	mi := &file_protos_editor_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EditReq.ProtoReflect.Descriptor instead.
func (*EditReq) Descriptor() ([]byte, []int) {
	return file_protos_editor_proto_rawDescGZIP(), []int{4}
}

func (x *EditReq) GetDocId() string {
	if x != nil {
		return x.DocId
	}
	return ""
}

func (x *EditReq) GetRev() int32 {
	if x != nil {
		return x.Rev
	}
	return 0
}

func (x *EditReq) GetOps() []*Op {
	if x != nil {
		return x.Ops
	}
	return nil
}

func (x *EditReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *EditReq) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

type Ack struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Ack) Reset() {
	*x = Ack{}
	mi := &file_protos_editor_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Ack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ack) ProtoMessage() {}

func (x *Ack) ProtoReflect() protoreflect.Message {
	mi := &file_protos_editor_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ack.ProtoReflect.Descriptor instead.
func (*Ack) Descriptor() ([]byte, []int) {
	return file_protos_editor_proto_rawDescGZIP(), []int{5}
}

type Op struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	N             int32                  `protobuf:"varint,1,opt,name=n,proto3" json:"n,omitempty"`
	S             string                 `protobuf:"bytes,2,opt,name=s,proto3" json:"s,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Op) Reset() {
	*x = Op{}
	mi := &file_protos_editor_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Op) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Op) ProtoMessage() {}

func (x *Op) ProtoReflect() protoreflect.Message {
	mi := &file_protos_editor_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Op.ProtoReflect.Descriptor instead.
func (*Op) Descriptor() ([]byte, []int) {
	return file_protos_editor_proto_rawDescGZIP(), []int{6}
}

func (x *Op) GetN() int32 {
	if x != nil {
		return x.N
	}
	return 0
}

func (x *Op) GetS() string {
	if x != nil {
		return x.S
	}
	return ""
}

type WatchReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	DocId         string                 `protobuf:"bytes,1,opt,name=doc_id,json=docId,proto3" json:"doc_id,omitempty"`
	UserId        string                 `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *WatchReq) Reset() {
	*x = WatchReq{}
	mi := &file_protos_editor_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WatchReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WatchReq) ProtoMessage() {}

func (x *WatchReq) ProtoReflect() protoreflect.Message {
	mi := &file_protos_editor_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WatchReq.ProtoReflect.Descriptor instead.
func (*WatchReq) Descriptor() ([]byte, []int) {
	return file_protos_editor_proto_rawDescGZIP(), []int{7}
}

func (x *WatchReq) GetDocId() string {
	if x != nil {
		return x.DocId
	}
	return ""
}

func (x *WatchReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type Update struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ops           []*Op                  `protobuf:"bytes,1,rep,name=ops,proto3" json:"ops,omitempty"`
	Doc           []byte                 `protobuf:"bytes,2,opt,name=doc,proto3" json:"doc,omitempty"`
	Rev           int32                  `protobuf:"varint,3,opt,name=rev,proto3" json:"rev,omitempty"`
	Title         string                 `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Update) Reset() {
	*x = Update{}
	mi := &file_protos_editor_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Update) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Update) ProtoMessage() {}

func (x *Update) ProtoReflect() protoreflect.Message {
	mi := &file_protos_editor_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Update.ProtoReflect.Descriptor instead.
func (*Update) Descriptor() ([]byte, []int) {
	return file_protos_editor_proto_rawDescGZIP(), []int{8}
}

func (x *Update) GetOps() []*Op {
	if x != nil {
		return x.Ops
	}
	return nil
}

func (x *Update) GetDoc() []byte {
	if x != nil {
		return x.Doc
	}
	return nil
}

func (x *Update) GetRev() int32 {
	if x != nil {
		return x.Rev
	}
	return 0
}

func (x *Update) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

var File_protos_editor_proto protoreflect.FileDescriptor

var file_protos_editor_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x22, 0x50, 0x0a,
	0x08, 0x53, 0x68, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x12, 0x19, 0x0a, 0x08, 0x64, 0x6f, 0x63,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x6f, 0x63,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x6f, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x03, 0x64, 0x6f, 0x63, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22,
	0x23, 0x0a, 0x0a, 0x53, 0x68, 0x61, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x15, 0x0a,
	0x06, 0x64, 0x6f, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64,
	0x6f, 0x63, 0x49, 0x64, 0x22, 0x3b, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x12, 0x15, 0x0a, 0x06, 0x64, 0x6f, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x64, 0x6f, 0x63, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x22, 0x0d, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x7f, 0x0a, 0x07, 0x45, 0x64, 0x69, 0x74, 0x52, 0x65, 0x71, 0x12, 0x15, 0x0a, 0x06, 0x64,
	0x6f, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x6f, 0x63,
	0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x72, 0x65, 0x76, 0x12, 0x1c, 0x0a, 0x03, 0x6f, 0x70, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0a, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x70, 0x52, 0x03, 0x6f,
	0x70, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x69, 0x74, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c,
	0x65, 0x22, 0x05, 0x0a, 0x03, 0x41, 0x63, 0x6b, 0x22, 0x20, 0x0a, 0x02, 0x4f, 0x70, 0x12, 0x0c,
	0x0a, 0x01, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x6e, 0x12, 0x0c, 0x0a, 0x01,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x73, 0x22, 0x3a, 0x0a, 0x08, 0x57, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x12, 0x15, 0x0a, 0x06, 0x64, 0x6f, 0x63, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x6f, 0x63, 0x49, 0x64, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x60, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x12, 0x1c, 0x0a, 0x03, 0x6f, 0x70, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x4f, 0x70, 0x52, 0x03, 0x6f, 0x70, 0x73, 0x12, 0x10,
	0x0a, 0x03, 0x64, 0x6f, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x64, 0x6f, 0x63,
	0x12, 0x10, 0x0a, 0x03, 0x72, 0x65, 0x76, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x72,
	0x65, 0x76, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x32, 0xcc, 0x01, 0x0a, 0x04, 0x4e, 0x6f, 0x64,
	0x65, 0x12, 0x2f, 0x0a, 0x05, 0x53, 0x68, 0x61, 0x72, 0x65, 0x12, 0x10, 0x2e, 0x65, 0x64, 0x69,
	0x74, 0x6f, 0x72, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x65,
	0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x22, 0x00, 0x12, 0x32, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x65,
	0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x1a,
	0x13, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x26, 0x0a, 0x04, 0x45, 0x64, 0x69, 0x74, 0x12, 0x0f,
	0x2e, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x45, 0x64, 0x69, 0x74, 0x52, 0x65, 0x71, 0x1a,
	0x0b, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x41, 0x63, 0x6b, 0x22, 0x00, 0x12, 0x37,
	0x0a, 0x0d, 0x57, 0x61, 0x74, 0x63, 0x68, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x10, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65,
	0x71, 0x1a, 0x0e, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x20, 0x5a, 0x1e, 0x65, 0x64, 0x69, 0x74, 0x6f,
	0x72, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2f, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_protos_editor_proto_rawDescOnce sync.Once
	file_protos_editor_proto_rawDescData []byte
)

func file_protos_editor_proto_rawDescGZIP() []byte {
	file_protos_editor_proto_rawDescOnce.Do(func() {
		file_protos_editor_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_protos_editor_proto_rawDesc), len(file_protos_editor_proto_rawDesc)))
	})
	return file_protos_editor_proto_rawDescData
}

var file_protos_editor_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_protos_editor_proto_goTypes = []any{
	(*ShareReq)(nil),    // 0: editor.ShareReq
	(*ShareReply)(nil),  // 1: editor.ShareReply
	(*DeleteReq)(nil),   // 2: editor.DeleteReq
	(*DeleteReply)(nil), // 3: editor.DeleteReply
	(*EditReq)(nil),     // 4: editor.EditReq
	(*Ack)(nil),         // 5: editor.Ack
	(*Op)(nil),          // 6: editor.Op
	(*WatchReq)(nil),    // 7: editor.WatchReq
	(*Update)(nil),      // 8: editor.Update
}
var file_protos_editor_proto_depIdxs = []int32{
	6, // 0: editor.EditReq.ops:type_name -> editor.Op
	6, // 1: editor.Update.ops:type_name -> editor.Op
	0, // 2: editor.Node.Share:input_type -> editor.ShareReq
	2, // 3: editor.Node.Delete:input_type -> editor.DeleteReq
	4, // 4: editor.Node.Edit:input_type -> editor.EditReq
	7, // 5: editor.Node.WatchDocument:input_type -> editor.WatchReq
	1, // 6: editor.Node.Share:output_type -> editor.ShareReply
	3, // 7: editor.Node.Delete:output_type -> editor.DeleteReply
	5, // 8: editor.Node.Edit:output_type -> editor.Ack
	8, // 9: editor.Node.WatchDocument:output_type -> editor.Update
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_protos_editor_proto_init() }
func file_protos_editor_proto_init() {
	if File_protos_editor_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_protos_editor_proto_rawDesc), len(file_protos_editor_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_editor_proto_goTypes,
		DependencyIndexes: file_protos_editor_proto_depIdxs,
		MessageInfos:      file_protos_editor_proto_msgTypes,
	}.Build()
	File_protos_editor_proto = out.File
	file_protos_editor_proto_goTypes = nil
	file_protos_editor_proto_depIdxs = nil
}
