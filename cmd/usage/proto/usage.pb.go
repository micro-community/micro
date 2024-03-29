// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: cmd/usage/proto/usage.proto

package usage

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

type Usage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// service name
	Service string `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	// version of service
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// unique service id
	Id string `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`
	// unix timestamp of report
	Timestamp uint64 `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// window of report in seconds
	Window uint64 `protobuf:"varint,5,opt,name=window,proto3" json:"window,omitempty"`
	// usage metrics
	Metrics *Metrics `protobuf:"bytes,6,opt,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *Usage) Reset() {
	*x = Usage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_usage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Usage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Usage) ProtoMessage() {}

func (x *Usage) ProtoReflect() protoreflect.Message {
	mi := &file_usage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Usage.ProtoReflect.Descriptor instead.
func (*Usage) Descriptor() ([]byte, []int) {
	return file_usage_proto_rawDescGZIP(), []int{0}
}

func (x *Usage) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *Usage) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *Usage) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Usage) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Usage) GetWindow() uint64 {
	if x != nil {
		return x.Window
	}
	return 0
}

func (x *Usage) GetMetrics() *Metrics {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type Metrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// counts such as requests, services, etc
	Count map[string]uint64 `protobuf:"bytes,1,rep,name=count,proto3" json:"count,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *Metrics) Reset() {
	*x = Metrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_usage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metrics) ProtoMessage() {}

func (x *Metrics) ProtoReflect() protoreflect.Message {
	mi := &file_usage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metrics.ProtoReflect.Descriptor instead.
func (*Metrics) Descriptor() ([]byte, []int) {
	return file_usage_proto_rawDescGZIP(), []int{1}
}

func (x *Metrics) GetCount() map[string]uint64 {
	if x != nil {
		return x.Count
	}
	return nil
}

var File_usage_proto protoreflect.FileDescriptor

var file_usage_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x75,
	0x73, 0x61, 0x67, 0x65, 0x22, 0xab, 0x01, 0x0a, 0x05, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x16, 0x0a, 0x06, 0x77, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x06, 0x77, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x12, 0x28, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x22, 0x74, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x2f, 0x0a,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x75,
	0x73, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x1a, 0x38,
	0x0a, 0x0a, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x75,
	0x73, 0x61, 0x67, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_usage_proto_rawDescOnce sync.Once
	file_usage_proto_rawDescData = file_usage_proto_rawDesc
)

func file_usage_proto_rawDescGZIP() []byte {
	file_usage_proto_rawDescOnce.Do(func() {
		file_usage_proto_rawDescData = protoimpl.X.CompressGZIP(file_usage_proto_rawDescData)
	})
	return file_usage_proto_rawDescData
}

var file_usage_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_usage_proto_goTypes = []interface{}{
	(*Usage)(nil),   // 0: usage.Usage
	(*Metrics)(nil), // 1: usage.Metrics
	nil,             // 2: usage.Metrics.CountEntry
}
var file_usage_proto_depIdxs = []int32{
	1, // 0: usage.Usage.metrics:type_name -> usage.Metrics
	2, // 1: usage.Metrics.count:type_name -> usage.Metrics.CountEntry
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_usage_proto_init() }
func file_usage_proto_init() {
	if File_usage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_usage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Usage); i {
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
		file_usage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metrics); i {
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
			RawDescriptor: file_usage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_usage_proto_goTypes,
		DependencyIndexes: file_usage_proto_depIdxs,
		MessageInfos:      file_usage_proto_msgTypes,
	}.Build()
	File_usage_proto = out.File
	file_usage_proto_rawDesc = nil
	file_usage_proto_goTypes = nil
	file_usage_proto_depIdxs = nil
}
