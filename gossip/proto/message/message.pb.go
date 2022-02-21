// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: message.proto

package message

import (
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

type Memory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From      uint32   `protobuf:"varint,1,opt,name=from,proto3" json:"from,omitempty"`
	Target    uint32   `protobuf:"varint,2,opt,name=target,proto3" json:"target,omitempty"`
	Checksum  string   `protobuf:"bytes,3,opt,name=checksum,proto3" json:"checksum,omitempty"`
	Signature string   `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	Content   []uint32 `protobuf:"varint,5,rep,packed,name=content,proto3" json:"content,omitempty"`
}

func (x *Memory) Reset() {
	*x = Memory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Memory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Memory) ProtoMessage() {}

func (x *Memory) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Memory.ProtoReflect.Descriptor instead.
func (*Memory) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{0}
}

func (x *Memory) GetFrom() uint32 {
	if x != nil {
		return x.From
	}
	return 0
}

func (x *Memory) GetTarget() uint32 {
	if x != nil {
		return x.Target
	}
	return 0
}

func (x *Memory) GetChecksum() string {
	if x != nil {
		return x.Checksum
	}
	return ""
}

func (x *Memory) GetSignature() string {
	if x != nil {
		return x.Signature
	}
	return ""
}

func (x *Memory) GetContent() []uint32 {
	if x != nil {
		return x.Content
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MergedMemory  *Memory   `protobuf:"bytes,1,opt,name=mergedMemory,proto3" json:"mergedMemory,omitempty"`
	Parents       []*Memory `protobuf:"bytes,2,rep,name=parents,proto3" json:"parents,omitempty"`
	PeeringProofs []string  `protobuf:"bytes,3,rep,name=peeringProofs,proto3" json:"peeringProofs,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_message_proto_rawDescGZIP(), []int{1}
}

func (x *Message) GetMergedMemory() *Memory {
	if x != nil {
		return x.MergedMemory
	}
	return nil
}

func (x *Message) GetParents() []*Memory {
	if x != nil {
		return x.Parents
	}
	return nil
}

func (x *Message) GetPeeringProofs() []string {
	if x != nil {
		return x.PeeringProofs
	}
	return nil
}

var File_message_proto protoreflect.FileDescriptor

var file_message_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x88, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x6d,
	0x6f, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75, 0x6d, 0x12, 0x1c, 0x0a, 0x09, 0x73,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x22, 0x8f, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x33, 0x0a, 0x0c, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x64, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x52, 0x0c, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x64, 0x4d, 0x65,
	0x6d, 0x6f, 0x72, 0x79, 0x12, 0x29, 0x0a, 0x07, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x52, 0x07, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x73, 0x12,
	0x24, 0x0a, 0x0d, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x50,
	0x72, 0x6f, 0x6f, 0x66, 0x73, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x4c, 0x65, 0x6b, 0x73, 0x73, 0x61, 0x79, 0x73, 0x2f, 0x41, 0x44, 0x65,
	0x4c, 0x65, 0x2f, 0x67, 0x6f, 0x73, 0x73, 0x69, 0x70, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_message_proto_rawDescOnce sync.Once
	file_message_proto_rawDescData = file_message_proto_rawDesc
)

func file_message_proto_rawDescGZIP() []byte {
	file_message_proto_rawDescOnce.Do(func() {
		file_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_message_proto_rawDescData)
	})
	return file_message_proto_rawDescData
}

var file_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_message_proto_goTypes = []interface{}{
	(*Memory)(nil),  // 0: message.Memory
	(*Message)(nil), // 1: message.Message
}
var file_message_proto_depIdxs = []int32{
	0, // 0: message.Message.mergedMemory:type_name -> message.Memory
	0, // 1: message.Message.parents:type_name -> message.Memory
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_message_proto_init() }
func file_message_proto_init() {
	if File_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Memory); i {
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
		file_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
			RawDescriptor: file_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_message_proto_goTypes,
		DependencyIndexes: file_message_proto_depIdxs,
		MessageInfos:      file_message_proto_msgTypes,
	}.Build()
	File_message_proto = out.File
	file_message_proto_rawDesc = nil
	file_message_proto_goTypes = nil
	file_message_proto_depIdxs = nil
}