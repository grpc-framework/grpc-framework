// Please use the following editor setup for this file:
// Tab size=2; Tabs as spaces; Clean up trailing whitepsace
//
// In vim add: au FileType proto setl sw=2 ts=2 expandtab list
//
// In vscode install vscode-proto3 extension and add this to your settings.json:
//    "[proto3]": {
//        "editor.tabSize": 2,
//        "editor.insertSpaces": true,
//        "editor.rulers": [80],
//        "editor.detectIndentation": true,
//        "files.trimTrailingWhitespace": true
//    }
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.3
// source: api/hello.proto

package api

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// These values are constants that can be used by the
// client and server applications
type HelloVersion_Version int32

const (
	// Must be set in the proto file; ignore.
	HelloVersion_VERSION_UNSPECIFIED HelloVersion_Version = 0
	// Version major value of this specification
	HelloVersion_MAJOR HelloVersion_Version = 0
	// Version minor value of this specification
	HelloVersion_MINOR HelloVersion_Version = 0
	// Version patch value of this specification
	HelloVersion_PATCH HelloVersion_Version = 1
)

// Enum value maps for HelloVersion_Version.
var (
	HelloVersion_Version_name = map[int32]string{
		0: "VERSION_UNSPECIFIED",
		// Duplicate value: 0: "MAJOR",
		// Duplicate value: 0: "MINOR",
		1: "PATCH",
	}
	HelloVersion_Version_value = map[string]int32{
		"VERSION_UNSPECIFIED": 0,
		"MAJOR":               0,
		"MINOR":               0,
		"PATCH":               1,
	}
)

func (x HelloVersion_Version) Enum() *HelloVersion_Version {
	p := new(HelloVersion_Version)
	*p = x
	return p
}

func (x HelloVersion_Version) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (HelloVersion_Version) Descriptor() protoreflect.EnumDescriptor {
	return file_api_hello_proto_enumTypes[0].Descriptor()
}

func (HelloVersion_Version) Type() protoreflect.EnumType {
	return &file_api_hello_proto_enumTypes[0]
}

func (x HelloVersion_Version) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use HelloVersion_Version.Descriptor instead.
func (HelloVersion_Version) EnumDescriptor() ([]byte, []int) {
	return file_api_hello_proto_rawDescGZIP(), []int{4, 0}
}

// The request message containing the user's name.
type HelloGreeterSayHelloRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name to say hello to
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *HelloGreeterSayHelloRequest) Reset() {
	*x = HelloGreeterSayHelloRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_hello_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloGreeterSayHelloRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloGreeterSayHelloRequest) ProtoMessage() {}

func (x *HelloGreeterSayHelloRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_hello_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloGreeterSayHelloRequest.ProtoReflect.Descriptor instead.
func (*HelloGreeterSayHelloRequest) Descriptor() ([]byte, []int) {
	return file_api_hello_proto_rawDescGZIP(), []int{0}
}

func (x *HelloGreeterSayHelloRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// The response message containing the greetings
type HelloGreeterSayHelloResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Message returned from server containing 'name'
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *HelloGreeterSayHelloResponse) Reset() {
	*x = HelloGreeterSayHelloResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_hello_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloGreeterSayHelloResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloGreeterSayHelloResponse) ProtoMessage() {}

func (x *HelloGreeterSayHelloResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_hello_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloGreeterSayHelloResponse.ProtoReflect.Descriptor instead.
func (*HelloGreeterSayHelloResponse) Descriptor() ([]byte, []int) {
	return file_api_hello_proto_rawDescGZIP(), []int{1}
}

func (x *HelloGreeterSayHelloResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// Empty request
type HelloIdentityVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HelloIdentityVersionRequest) Reset() {
	*x = HelloIdentityVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_hello_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloIdentityVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloIdentityVersionRequest) ProtoMessage() {}

func (x *HelloIdentityVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_hello_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloIdentityVersionRequest.ProtoReflect.Descriptor instead.
func (*HelloIdentityVersionRequest) Descriptor() ([]byte, []int) {
	return file_api_hello_proto_rawDescGZIP(), []int{2}
}

// Defines the response to version
type HelloIdentityVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Hello application version
	ServerVersion *HelloVersion `protobuf:"bytes,1,opt,name=server_version,json=serverVersion,proto3" json:"server_version,omitempty"`
}

func (x *HelloIdentityVersionResponse) Reset() {
	*x = HelloIdentityVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_hello_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloIdentityVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloIdentityVersionResponse) ProtoMessage() {}

func (x *HelloIdentityVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_hello_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloIdentityVersionResponse.ProtoReflect.Descriptor instead.
func (*HelloIdentityVersionResponse) Descriptor() ([]byte, []int) {
	return file_api_hello_proto_rawDescGZIP(), []int{3}
}

func (x *HelloIdentityVersionResponse) GetServerVersion() *HelloVersion {
	if x != nil {
		return x.ServerVersion
	}
	return nil
}

// Hello version in Major.Minor.Patch format. The goal of this
// message is to provide clients a method to determine the server
// and client versions.
type HelloVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Version major number
	Major int32 `protobuf:"varint,1,opt,name=major,proto3" json:"major,omitempty"`
	// Version minor number
	Minor int32 `protobuf:"varint,2,opt,name=minor,proto3" json:"minor,omitempty"`
	// Version patch number
	Patch int32 `protobuf:"varint,3,opt,name=patch,proto3" json:"patch,omitempty"`
	// String representation of the version. Must be
	// in `major.minor.patch` format.
	Version string `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *HelloVersion) Reset() {
	*x = HelloVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_hello_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HelloVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HelloVersion) ProtoMessage() {}

func (x *HelloVersion) ProtoReflect() protoreflect.Message {
	mi := &file_api_hello_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HelloVersion.ProtoReflect.Descriptor instead.
func (*HelloVersion) Descriptor() ([]byte, []int) {
	return file_api_hello_proto_rawDescGZIP(), []int{4}
}

func (x *HelloVersion) GetMajor() int32 {
	if x != nil {
		return x.Major
	}
	return 0
}

func (x *HelloVersion) GetMinor() int32 {
	if x != nil {
		return x.Minor
	}
	return 0
}

func (x *HelloVersion) GetPatch() int32 {
	if x != nil {
		return x.Patch
	}
	return 0
}

func (x *HelloVersion) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

var File_api_hello_proto protoreflect.FileDescriptor

var file_api_hello_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x61, 0x70, 0x69, 0x2f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x08, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x31, 0x0a, 0x1b, 0x48, 0x65, 0x6c,
	0x6c, 0x6f, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x53, 0x61, 0x79, 0x48, 0x65, 0x6c, 0x6c,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x38, 0x0a, 0x1c,
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x53, 0x61, 0x79, 0x48,
	0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x1d, 0x0a, 0x1b, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x49,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x5d, 0x0a, 0x1c, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x49, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x22, 0xb3, 0x01, 0x0a, 0x0c, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6d, 0x61, 0x6a, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x6d,
	0x69, 0x6e, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6d, 0x69, 0x6e, 0x6f,
	0x72, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x61, 0x74, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x70, 0x61, 0x74, 0x63, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x22, 0x47, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a, 0x13,
	0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x4d, 0x41, 0x4a, 0x4f, 0x52, 0x10, 0x00,
	0x12, 0x09, 0x0a, 0x05, 0x4d, 0x49, 0x4e, 0x4f, 0x52, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x50,
	0x41, 0x54, 0x43, 0x48, 0x10, 0x01, 0x1a, 0x02, 0x10, 0x01, 0x32, 0x8a, 0x01, 0x0a, 0x0c, 0x48,
	0x65, 0x6c, 0x6c, 0x6f, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x12, 0x7a, 0x0a, 0x08, 0x53,
	0x61, 0x79, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x25, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e,
	0x76, 0x31, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x53,
	0x61, 0x79, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x26,
	0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x47,
	0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x53, 0x61, 0x79, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1f, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x19, 0x22, 0x14,
	0x2f, 0x76, 0x31, 0x2f, 0x67, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x3a, 0x73, 0x61, 0x79, 0x48,
	0x65, 0x6c, 0x6c, 0x6f, 0x3a, 0x01, 0x2a, 0x32, 0x94, 0x01, 0x0a, 0x0d, 0x48, 0x65, 0x6c, 0x6c,
	0x6f, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x82, 0x01, 0x0a, 0x0d, 0x53, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x25, 0x2e, 0x68, 0x65,
	0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x49, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x26, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x48, 0x65,
	0x6c, 0x6c, 0x6f, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x22, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x1c, 0x12, 0x1a, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x3a, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x27,
	0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x76, 0x31, 0x42, 0x0a,
	0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x09, 0x2e, 0x2f,
	0x61, 0x70, 0x69, 0x3b, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_hello_proto_rawDescOnce sync.Once
	file_api_hello_proto_rawDescData = file_api_hello_proto_rawDesc
)

func file_api_hello_proto_rawDescGZIP() []byte {
	file_api_hello_proto_rawDescOnce.Do(func() {
		file_api_hello_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_hello_proto_rawDescData)
	})
	return file_api_hello_proto_rawDescData
}

var file_api_hello_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_hello_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_hello_proto_goTypes = []interface{}{
	(HelloVersion_Version)(0),            // 0: hello.v1.HelloVersion.Version
	(*HelloGreeterSayHelloRequest)(nil),  // 1: hello.v1.HelloGreeterSayHelloRequest
	(*HelloGreeterSayHelloResponse)(nil), // 2: hello.v1.HelloGreeterSayHelloResponse
	(*HelloIdentityVersionRequest)(nil),  // 3: hello.v1.HelloIdentityVersionRequest
	(*HelloIdentityVersionResponse)(nil), // 4: hello.v1.HelloIdentityVersionResponse
	(*HelloVersion)(nil),                 // 5: hello.v1.HelloVersion
}
var file_api_hello_proto_depIdxs = []int32{
	5, // 0: hello.v1.HelloIdentityVersionResponse.server_version:type_name -> hello.v1.HelloVersion
	1, // 1: hello.v1.HelloGreeter.SayHello:input_type -> hello.v1.HelloGreeterSayHelloRequest
	3, // 2: hello.v1.HelloIdentity.ServerVersion:input_type -> hello.v1.HelloIdentityVersionRequest
	2, // 3: hello.v1.HelloGreeter.SayHello:output_type -> hello.v1.HelloGreeterSayHelloResponse
	4, // 4: hello.v1.HelloIdentity.ServerVersion:output_type -> hello.v1.HelloIdentityVersionResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_hello_proto_init() }
func file_api_hello_proto_init() {
	if File_api_hello_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_hello_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloGreeterSayHelloRequest); i {
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
		file_api_hello_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloGreeterSayHelloResponse); i {
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
		file_api_hello_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloIdentityVersionRequest); i {
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
		file_api_hello_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloIdentityVersionResponse); i {
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
		file_api_hello_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HelloVersion); i {
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
			RawDescriptor: file_api_hello_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_api_hello_proto_goTypes,
		DependencyIndexes: file_api_hello_proto_depIdxs,
		EnumInfos:         file_api_hello_proto_enumTypes,
		MessageInfos:      file_api_hello_proto_msgTypes,
	}.Build()
	File_api_hello_proto = out.File
	file_api_hello_proto_rawDesc = nil
	file_api_hello_proto_goTypes = nil
	file_api_hello_proto_depIdxs = nil
}
