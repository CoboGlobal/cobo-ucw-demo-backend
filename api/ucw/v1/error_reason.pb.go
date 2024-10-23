// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: ucw/v1/error_reason.proto

package v1

import (
	_ "github.com/go-kratos/kratos/v2/errors"
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

type ErrorReason int32

const (
	ErrorReason_GREETER_UNSPECIFIED          ErrorReason = 0
	ErrorReason_USER_NOT_FOUND               ErrorReason = 1
	ErrorReason_COBO_NODE_INVALID            ErrorReason = 2
	ErrorReason_KEY_GEN_DIFFERENT_USERS      ErrorReason = 3
	ErrorReason_UNSUPPORTED_TRANSACTION_TYPE ErrorReason = 4
	ErrorReason_VAULT_NOT_FOUND              ErrorReason = 5
	ErrorReason_UNAUTHORIZED                 ErrorReason = 6
	ErrorReason_INVALID_REQUEST_PARAMS       ErrorReason = 7
)

// Enum value maps for ErrorReason.
var (
	ErrorReason_name = map[int32]string{
		0: "GREETER_UNSPECIFIED",
		1: "USER_NOT_FOUND",
		2: "COBO_NODE_INVALID",
		3: "KEY_GEN_DIFFERENT_USERS",
		4: "UNSUPPORTED_TRANSACTION_TYPE",
		5: "VAULT_NOT_FOUND",
		6: "UNAUTHORIZED",
		7: "INVALID_REQUEST_PARAMS",
	}
	ErrorReason_value = map[string]int32{
		"GREETER_UNSPECIFIED":          0,
		"USER_NOT_FOUND":               1,
		"COBO_NODE_INVALID":            2,
		"KEY_GEN_DIFFERENT_USERS":      3,
		"UNSUPPORTED_TRANSACTION_TYPE": 4,
		"VAULT_NOT_FOUND":              5,
		"UNAUTHORIZED":                 6,
		"INVALID_REQUEST_PARAMS":       7,
	}
)

func (x ErrorReason) Enum() *ErrorReason {
	p := new(ErrorReason)
	*p = x
	return p
}

func (x ErrorReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ErrorReason) Descriptor() protoreflect.EnumDescriptor {
	return file_ucw_v1_error_reason_proto_enumTypes[0].Descriptor()
}

func (ErrorReason) Type() protoreflect.EnumType {
	return &file_ucw_v1_error_reason_proto_enumTypes[0]
}

func (x ErrorReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ErrorReason.Descriptor instead.
func (ErrorReason) EnumDescriptor() ([]byte, []int) {
	return file_ucw_v1_error_reason_proto_rawDescGZIP(), []int{0}
}

var File_ucw_v1_error_reason_proto protoreflect.FileDescriptor

var file_ucw_v1_error_reason_proto_rawDesc = []byte{
	0x0a, 0x19, 0x75, 0x63, 0x77, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x68, 0x65, 0x6c,
	0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x73, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a,
	0xeb, 0x01, 0x0a, 0x0b, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12,
	0x17, 0x0a, 0x13, 0x47, 0x52, 0x45, 0x45, 0x54, 0x45, 0x52, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45,
	0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x0e, 0x55, 0x53, 0x45, 0x52,
	0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x1a, 0x04, 0xa8, 0x45,
	0x94, 0x03, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x4f, 0x42, 0x4f, 0x5f, 0x4e, 0x4f, 0x44, 0x45, 0x5f,
	0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x4b, 0x45, 0x59,
	0x5f, 0x47, 0x45, 0x4e, 0x5f, 0x44, 0x49, 0x46, 0x46, 0x45, 0x52, 0x45, 0x4e, 0x54, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x53, 0x10, 0x03, 0x12, 0x20, 0x0a, 0x1c, 0x55, 0x4e, 0x53, 0x55, 0x50, 0x50,
	0x4f, 0x52, 0x54, 0x45, 0x44, 0x5f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41, 0x43, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x10, 0x04, 0x12, 0x19, 0x0a, 0x0f, 0x56, 0x41, 0x55, 0x4c,
	0x54, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x05, 0x1a, 0x04, 0xa8,
	0x45, 0x94, 0x03, 0x12, 0x16, 0x0a, 0x0c, 0x55, 0x4e, 0x41, 0x55, 0x54, 0x48, 0x4f, 0x52, 0x49,
	0x5a, 0x45, 0x44, 0x10, 0x06, 0x1a, 0x04, 0xa8, 0x45, 0x93, 0x03, 0x12, 0x1a, 0x0a, 0x16, 0x49,
	0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x5f, 0x50,
	0x41, 0x52, 0x41, 0x4d, 0x53, 0x10, 0x07, 0x1a, 0x04, 0xa0, 0x45, 0x90, 0x03, 0x42, 0x4a, 0x0a,
	0x0d, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x2e, 0x76, 0x31, 0x50, 0x01,
	0x5a, 0x25, 0x63, 0x6f, 0x62, 0x6f, 0x2d, 0x75, 0x63, 0x77, 0x2d, 0x62, 0x61, 0x63, 0x6b, 0x65,
	0x6e, 0x64, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c,
	0x64, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0xa2, 0x02, 0x0f, 0x41, 0x50, 0x49, 0x48, 0x65, 0x6c,
	0x6c, 0x6f, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_ucw_v1_error_reason_proto_rawDescOnce sync.Once
	file_ucw_v1_error_reason_proto_rawDescData = file_ucw_v1_error_reason_proto_rawDesc
)

func file_ucw_v1_error_reason_proto_rawDescGZIP() []byte {
	file_ucw_v1_error_reason_proto_rawDescOnce.Do(func() {
		file_ucw_v1_error_reason_proto_rawDescData = protoimpl.X.CompressGZIP(file_ucw_v1_error_reason_proto_rawDescData)
	})
	return file_ucw_v1_error_reason_proto_rawDescData
}

var file_ucw_v1_error_reason_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_ucw_v1_error_reason_proto_goTypes = []any{
	(ErrorReason)(0), // 0: helloworld.v1.ErrorReason
}
var file_ucw_v1_error_reason_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ucw_v1_error_reason_proto_init() }
func file_ucw_v1_error_reason_proto_init() {
	if File_ucw_v1_error_reason_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ucw_v1_error_reason_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ucw_v1_error_reason_proto_goTypes,
		DependencyIndexes: file_ucw_v1_error_reason_proto_depIdxs,
		EnumInfos:         file_ucw_v1_error_reason_proto_enumTypes,
	}.Build()
	File_ucw_v1_error_reason_proto = out.File
	file_ucw_v1_error_reason_proto_rawDesc = nil
	file_ucw_v1_error_reason_proto_goTypes = nil
	file_ucw_v1_error_reason_proto_depIdxs = nil
}