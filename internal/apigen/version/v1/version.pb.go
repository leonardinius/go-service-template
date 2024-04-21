// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: version/v1/version.proto

package versionv1

import (
	v1 "github.com/leonardinius/go-service-template/internal/apigen/shared/v1"
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

// VcsType is the type of version control system.
type VcsType int32

const (
	VcsType_VCS_TYPE_UNSPECIFIED VcsType = 0
	VcsType_VCS_TYPE_GIT         VcsType = 1
)

// Enum value maps for VcsType.
var (
	VcsType_name = map[int32]string{
		0: "VCS_TYPE_UNSPECIFIED",
		1: "VCS_TYPE_GIT",
	}
	VcsType_value = map[string]int32{
		"VCS_TYPE_UNSPECIFIED": 0,
		"VCS_TYPE_GIT":         1,
	}
)

func (x VcsType) Enum() *VcsType {
	p := new(VcsType)
	*p = x
	return p
}

func (x VcsType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VcsType) Descriptor() protoreflect.EnumDescriptor {
	return file_version_v1_version_proto_enumTypes[0].Descriptor()
}

func (VcsType) Type() protoreflect.EnumType {
	return &file_version_v1_version_proto_enumTypes[0]
}

func (x VcsType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VcsType.Descriptor instead.
func (VcsType) EnumDescriptor() ([]byte, []int) {
	return file_version_v1_version_proto_rawDescGZIP(), []int{0}
}

// Version is the version information of the service.
type Version struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Vcs         VcsType `protobuf:"varint,1,opt,name=vcs,proto3,enum=version.v1.VcsType" json:"vcs,omitempty"`
	ServiceName string  `protobuf:"bytes,2,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	RefName     string  `protobuf:"bytes,3,opt,name=ref_name,json=refName,proto3" json:"ref_name,omitempty"`
	Commit      string  `protobuf:"bytes,4,opt,name=commit,proto3" json:"commit,omitempty"`
	BuildTime   string  `protobuf:"bytes,5,opt,name=build_time,json=buildTime,proto3" json:"build_time,omitempty"`
	FullVersion string  `protobuf:"bytes,6,opt,name=full_version,json=fullVersion,proto3" json:"full_version,omitempty"`
}

func (x *Version) Reset() {
	*x = Version{}
	if protoimpl.UnsafeEnabled {
		mi := &file_version_v1_version_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Version) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Version) ProtoMessage() {}

func (x *Version) ProtoReflect() protoreflect.Message {
	mi := &file_version_v1_version_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Version.ProtoReflect.Descriptor instead.
func (*Version) Descriptor() ([]byte, []int) {
	return file_version_v1_version_proto_rawDescGZIP(), []int{0}
}

func (x *Version) GetVcs() VcsType {
	if x != nil {
		return x.Vcs
	}
	return VcsType_VCS_TYPE_UNSPECIFIED
}

func (x *Version) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *Version) GetRefName() string {
	if x != nil {
		return x.RefName
	}
	return ""
}

func (x *Version) GetCommit() string {
	if x != nil {
		return x.Commit
	}
	return ""
}

func (x *Version) GetBuildTime() string {
	if x != nil {
		return x.BuildTime
	}
	return ""
}

func (x *Version) GetFullVersion() string {
	if x != nil {
		return x.FullVersion
	}
	return ""
}

type GetVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetVersionRequest) Reset() {
	*x = GetVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_version_v1_version_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersionRequest) ProtoMessage() {}

func (x *GetVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_version_v1_version_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersionRequest.ProtoReflect.Descriptor instead.
func (*GetVersionRequest) Descriptor() ([]byte, []int) {
	return file_version_v1_version_proto_rawDescGZIP(), []int{1}
}

type GetVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error   *v1.Error `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
	Version *Version  `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *GetVersionResponse) Reset() {
	*x = GetVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_version_v1_version_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVersionResponse) ProtoMessage() {}

func (x *GetVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_version_v1_version_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVersionResponse.ProtoReflect.Descriptor instead.
func (*GetVersionResponse) Descriptor() ([]byte, []int) {
	return file_version_v1_version_proto_rawDescGZIP(), []int{2}
}

func (x *GetVersionResponse) GetError() *v1.Error {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *GetVersionResponse) GetVersion() *Version {
	if x != nil {
		return x.Version
	}
	return nil
}

var File_version_v1_version_proto protoreflect.FileDescriptor

var file_version_v1_version_proto_rawDesc = []byte{
	0x0a, 0x18, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x15, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x76,
	0x31, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc8, 0x01,
	0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x03, 0x76, 0x63, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x56, 0x63, 0x73, 0x54, 0x79, 0x70, 0x65, 0x52, 0x03, 0x76, 0x63, 0x73,
	0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x72, 0x65, 0x66, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65, 0x66, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x75, 0x69, 0x6c,
	0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x66, 0x75, 0x6c, 0x6c, 0x5f, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x66, 0x75, 0x6c,
	0x6c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x13, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x6b, 0x0a,
	0x12, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x2d, 0x0a, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2a, 0x35, 0x0a, 0x07, 0x56, 0x63,
	0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x14, 0x56, 0x43, 0x53, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x10, 0x0a, 0x0c, 0x56, 0x43, 0x53, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x47, 0x49, 0x54, 0x10,
	0x01, 0x32, 0x62, 0x0a, 0x0e, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x50, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x1d, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1e, 0x2e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x03, 0x90, 0x02, 0x01, 0x42, 0xb9, 0x01, 0x0a, 0x0e, 0x63, 0x6f, 0x6d, 0x2e, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x50, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x65, 0x6f, 0x6e, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x69, 0x75,
	0x73, 0x2f, 0x67, 0x6f, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2d, 0x74, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61,
	0x70, 0x69, 0x67, 0x65, 0x6e, 0x2f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31,
	0x3b, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x56, 0x58, 0x58,
	0xaa, 0x02, 0x0a, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0a,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x16, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x0b, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_version_v1_version_proto_rawDescOnce sync.Once
	file_version_v1_version_proto_rawDescData = file_version_v1_version_proto_rawDesc
)

func file_version_v1_version_proto_rawDescGZIP() []byte {
	file_version_v1_version_proto_rawDescOnce.Do(func() {
		file_version_v1_version_proto_rawDescData = protoimpl.X.CompressGZIP(file_version_v1_version_proto_rawDescData)
	})
	return file_version_v1_version_proto_rawDescData
}

var file_version_v1_version_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_version_v1_version_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_version_v1_version_proto_goTypes = []interface{}{
	(VcsType)(0),               // 0: version.v1.VcsType
	(*Version)(nil),            // 1: version.v1.Version
	(*GetVersionRequest)(nil),  // 2: version.v1.GetVersionRequest
	(*GetVersionResponse)(nil), // 3: version.v1.GetVersionResponse
	(*v1.Error)(nil),           // 4: shared.v1.Error
}
var file_version_v1_version_proto_depIdxs = []int32{
	0, // 0: version.v1.Version.vcs:type_name -> version.v1.VcsType
	4, // 1: version.v1.GetVersionResponse.error:type_name -> shared.v1.Error
	1, // 2: version.v1.GetVersionResponse.version:type_name -> version.v1.Version
	2, // 3: version.v1.VersionService.GetVersion:input_type -> version.v1.GetVersionRequest
	3, // 4: version.v1.VersionService.GetVersion:output_type -> version.v1.GetVersionResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_version_v1_version_proto_init() }
func file_version_v1_version_proto_init() {
	if File_version_v1_version_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_version_v1_version_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Version); i {
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
		file_version_v1_version_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersionRequest); i {
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
		file_version_v1_version_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVersionResponse); i {
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
			RawDescriptor: file_version_v1_version_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_version_v1_version_proto_goTypes,
		DependencyIndexes: file_version_v1_version_proto_depIdxs,
		EnumInfos:         file_version_v1_version_proto_enumTypes,
		MessageInfos:      file_version_v1_version_proto_msgTypes,
	}.Build()
	File_version_v1_version_proto = out.File
	file_version_v1_version_proto_rawDesc = nil
	file_version_v1_version_proto_goTypes = nil
	file_version_v1_version_proto_depIdxs = nil
}
