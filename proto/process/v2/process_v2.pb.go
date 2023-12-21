// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: process/v2/process_v2.proto

package processv2

import (
	v1 "github.com/conduitio/conduit-connector-protocol/proto/opencdc/v1"
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

type ProcessRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Records []*v1.Record `protobuf:"bytes,1,rep,name=records,proto3" json:"records,omitempty"`
}

func (x *ProcessRequest) Reset() {
	*x = ProcessRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_v2_process_v2_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessRequest) ProtoMessage() {}

func (x *ProcessRequest) ProtoReflect() protoreflect.Message {
	mi := &file_process_v2_process_v2_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessRequest.ProtoReflect.Descriptor instead.
func (*ProcessRequest) Descriptor() ([]byte, []int) {
	return file_process_v2_process_v2_proto_rawDescGZIP(), []int{0}
}

func (x *ProcessRequest) GetRecords() []*v1.Record {
	if x != nil {
		return x.Records
	}
	return nil
}

type ProcessResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Records []*v1.Record `protobuf:"bytes,1,rep,name=records,proto3" json:"records,omitempty"`
}

func (x *ProcessResponse) Reset() {
	*x = ProcessResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_v2_process_v2_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessResponse) ProtoMessage() {}

func (x *ProcessResponse) ProtoReflect() protoreflect.Message {
	mi := &file_process_v2_process_v2_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessResponse.ProtoReflect.Descriptor instead.
func (*ProcessResponse) Descriptor() ([]byte, []int) {
	return file_process_v2_process_v2_proto_rawDescGZIP(), []int{1}
}

func (x *ProcessResponse) GetRecords() []*v1.Record {
	if x != nil {
		return x.Records
	}
	return nil
}

var File_process_v2_process_v2_proto protoreflect.FileDescriptor

var file_process_v2_process_v2_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2f, 0x76, 0x32, 0x2f, 0x70, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x5f, 0x76, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x70,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x32, 0x1a, 0x18, 0x6f, 0x70, 0x65, 0x6e, 0x63,
	0x64, 0x63, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x63, 0x64, 0x63, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x63, 0x64, 0x63,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x72, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x73, 0x22, 0x3f, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6f, 0x70, 0x65, 0x6e, 0x63, 0x64,
	0x63, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x73, 0x32, 0x58, 0x0a, 0x10, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f,
	0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x44, 0x0a, 0x07, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x32,
	0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x76, 0x32, 0x2e, 0x50, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x9f,
	0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x76,
	0x32, 0x42, 0x0e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x56, 0x32, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x48, 0x02, 0x50, 0x01, 0x5a, 0x32, 0x62, 0x75, 0x66, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64,
	0x2f, 0x6d, 0x65, 0x72, 0x6f, 0x78, 0x61, 0x2f, 0x74, 0x75, 0x72, 0x62, 0x69, 0x6e, 0x65, 0x2d,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2f, 0x76, 0x32, 0x3b,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x76, 0x32, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa,
	0x02, 0x0a, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x56, 0x32, 0xca, 0x02, 0x0a, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5c, 0x56, 0x32, 0xe2, 0x02, 0x16, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x5c, 0x56, 0x32, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x0b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x3a, 0x3a, 0x56, 0x32,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_process_v2_process_v2_proto_rawDescOnce sync.Once
	file_process_v2_process_v2_proto_rawDescData = file_process_v2_process_v2_proto_rawDesc
)

func file_process_v2_process_v2_proto_rawDescGZIP() []byte {
	file_process_v2_process_v2_proto_rawDescOnce.Do(func() {
		file_process_v2_process_v2_proto_rawDescData = protoimpl.X.CompressGZIP(file_process_v2_process_v2_proto_rawDescData)
	})
	return file_process_v2_process_v2_proto_rawDescData
}

var file_process_v2_process_v2_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_process_v2_process_v2_proto_goTypes = []interface{}{
	(*ProcessRequest)(nil),  // 0: process.v2.ProcessRequest
	(*ProcessResponse)(nil), // 1: process.v2.ProcessResponse
	(*v1.Record)(nil),       // 2: opencdc.v1.Record
}
var file_process_v2_process_v2_proto_depIdxs = []int32{
	2, // 0: process.v2.ProcessRequest.records:type_name -> opencdc.v1.Record
	2, // 1: process.v2.ProcessResponse.records:type_name -> opencdc.v1.Record
	0, // 2: process.v2.ProcessorService.Process:input_type -> process.v2.ProcessRequest
	1, // 3: process.v2.ProcessorService.Process:output_type -> process.v2.ProcessResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_process_v2_process_v2_proto_init() }
func file_process_v2_process_v2_proto_init() {
	if File_process_v2_process_v2_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_process_v2_process_v2_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessRequest); i {
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
		file_process_v2_process_v2_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessResponse); i {
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
			RawDescriptor: file_process_v2_process_v2_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_process_v2_process_v2_proto_goTypes,
		DependencyIndexes: file_process_v2_process_v2_proto_depIdxs,
		MessageInfos:      file_process_v2_process_v2_proto_msgTypes,
	}.Build()
	File_process_v2_process_v2_proto = out.File
	file_process_v2_process_v2_proto_rawDesc = nil
	file_process_v2_process_v2_proto_goTypes = nil
	file_process_v2_process_v2_proto_depIdxs = nil
}
