// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: enum.proto

package raft

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

type EntryType int32

const (
	EntryType_ENTRY_TYPE_UNKNOWN       EntryType = 0
	EntryType_ENTRY_TYPE_NO_OP         EntryType = 1
	EntryType_ENTRY_TYPE_DATA          EntryType = 2
	EntryType_ENTRY_TYPE_CONFIGURATION EntryType = 3
)

// Enum value maps for EntryType.
var (
	EntryType_name = map[int32]string{
		0: "ENTRY_TYPE_UNKNOWN",
		1: "ENTRY_TYPE_NO_OP",
		2: "ENTRY_TYPE_DATA",
		3: "ENTRY_TYPE_CONFIGURATION",
	}
	EntryType_value = map[string]int32{
		"ENTRY_TYPE_UNKNOWN":       0,
		"ENTRY_TYPE_NO_OP":         1,
		"ENTRY_TYPE_DATA":          2,
		"ENTRY_TYPE_CONFIGURATION": 3,
	}
)

func (x EntryType) Enum() *EntryType {
	p := new(EntryType)
	*p = x
	return p
}

func (x EntryType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EntryType) Descriptor() protoreflect.EnumDescriptor {
	return file_enum_proto_enumTypes[0].Descriptor()
}

func (EntryType) Type() protoreflect.EnumType {
	return &file_enum_proto_enumTypes[0]
}

func (x EntryType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EntryType.Descriptor instead.
func (EntryType) EnumDescriptor() ([]byte, []int) {
	return file_enum_proto_rawDescGZIP(), []int{0}
}

type ErrorType int32

const (
	ErrorType_ERROR_TYPE_NONE          ErrorType = 0
	ErrorType_ERROR_TYPE_LOG           ErrorType = 1
	ErrorType_ERROR_TYPE_STABLE        ErrorType = 2
	ErrorType_ERROR_TYPE_SNAPSHOT      ErrorType = 3
	ErrorType_ERROR_TYPE_STATE_MACHINE ErrorType = 4
)

// Enum value maps for ErrorType.
var (
	ErrorType_name = map[int32]string{
		0: "ERROR_TYPE_NONE",
		1: "ERROR_TYPE_LOG",
		2: "ERROR_TYPE_STABLE",
		3: "ERROR_TYPE_SNAPSHOT",
		4: "ERROR_TYPE_STATE_MACHINE",
	}
	ErrorType_value = map[string]int32{
		"ERROR_TYPE_NONE":          0,
		"ERROR_TYPE_LOG":           1,
		"ERROR_TYPE_STABLE":        2,
		"ERROR_TYPE_SNAPSHOT":      3,
		"ERROR_TYPE_STATE_MACHINE": 4,
	}
)

func (x ErrorType) Enum() *ErrorType {
	p := new(ErrorType)
	*p = x
	return p
}

func (x ErrorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorType) Descriptor() protoreflect.EnumDescriptor {
	return file_enum_proto_enumTypes[1].Descriptor()
}

func (ErrorType) Type() protoreflect.EnumType {
	return &file_enum_proto_enumTypes[1]
}

func (x ErrorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorType.Descriptor instead.
func (ErrorType) EnumDescriptor() ([]byte, []int) {
	return file_enum_proto_rawDescGZIP(), []int{1}
}

var File_enum_proto protoreflect.FileDescriptor

var file_enum_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x65, 0x6e, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x72, 0x61,
	0x66, 0x74, 0x2a, 0x6c, 0x0a, 0x09, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x16, 0x0a, 0x12, 0x45, 0x4e, 0x54, 0x52, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x45, 0x4e, 0x54, 0x52, 0x59,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x5f, 0x4f, 0x50, 0x10, 0x01, 0x12, 0x13, 0x0a,
	0x0f, 0x45, 0x4e, 0x54, 0x52, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x41, 0x54, 0x41,
	0x10, 0x02, 0x12, 0x1c, 0x0a, 0x18, 0x45, 0x4e, 0x54, 0x52, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x49, 0x47, 0x55, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x03,
	0x2a, 0x82, 0x01, 0x0a, 0x09, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x13,
	0x0a, 0x0f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x4e,
	0x45, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x4c, 0x4f, 0x47, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x02, 0x12, 0x17,
	0x0a, 0x13, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x4e, 0x41,
	0x50, 0x53, 0x48, 0x4f, 0x54, 0x10, 0x03, 0x12, 0x1c, 0x0a, 0x18, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x4d, 0x41, 0x43, 0x48,
	0x49, 0x4e, 0x45, 0x10, 0x04, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x6c, 0x6c, 0x65, 0x6e, 0x53, 0x68, 0x61, 0x77, 0x31, 0x39, 0x2f,
	0x72, 0x61, 0x66, 0x74, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x3b, 0x72, 0x61, 0x66, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_enum_proto_rawDescOnce sync.Once
	file_enum_proto_rawDescData = file_enum_proto_rawDesc
)

func file_enum_proto_rawDescGZIP() []byte {
	file_enum_proto_rawDescOnce.Do(func() {
		file_enum_proto_rawDescData = protoimpl.X.CompressGZIP(file_enum_proto_rawDescData)
	})
	return file_enum_proto_rawDescData
}

var file_enum_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_enum_proto_goTypes = []interface{}{
	(EntryType)(0), // 0: raft.EntryType
	(ErrorType)(0), // 1: raft.ErrorType
}
var file_enum_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_enum_proto_init() }
func file_enum_proto_init() {
	if File_enum_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_enum_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_enum_proto_goTypes,
		DependencyIndexes: file_enum_proto_depIdxs,
		EnumInfos:         file_enum_proto_enumTypes,
	}.Build()
	File_enum_proto = out.File
	file_enum_proto_rawDesc = nil
	file_enum_proto_goTypes = nil
	file_enum_proto_depIdxs = nil
}