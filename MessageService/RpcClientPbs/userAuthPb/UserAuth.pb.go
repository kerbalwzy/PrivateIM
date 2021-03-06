// Code generated by protoc-gen-go. DO NOT EDIT.
// source: UserAuth.proto

package userAuthPb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Token struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Token) Reset()         { *m = Token{} }
func (m *Token) String() string { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()    {}
func (*Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce91c7f1b4618f2e, []int{0}
}

func (m *Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Token.Unmarshal(m, b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Token.Marshal(b, m, deterministic)
}
func (m *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(m, src)
}
func (m *Token) XXX_Size() int {
	return xxx_messageInfo_Token.Size(m)
}
func (m *Token) XXX_DiscardUnknown() {
	xxx_messageInfo_Token.DiscardUnknown(m)
}

var xxx_messageInfo_Token proto.InternalMessageInfo

func (m *Token) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type Id struct {
	Value                int64    `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Id) Reset()         { *m = Id{} }
func (m *Id) String() string { return proto.CompactTextString(m) }
func (*Id) ProtoMessage()    {}
func (*Id) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce91c7f1b4618f2e, []int{1}
}

func (m *Id) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Id.Unmarshal(m, b)
}
func (m *Id) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Id.Marshal(b, m, deterministic)
}
func (m *Id) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Id.Merge(m, src)
}
func (m *Id) XXX_Size() int {
	return xxx_messageInfo_Id.Size(m)
}
func (m *Id) XXX_DiscardUnknown() {
	xxx_messageInfo_Id.DiscardUnknown(m)
}

var xxx_messageInfo_Id proto.InternalMessageInfo

func (m *Id) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*Token)(nil), "userAuthPb.Token")
	proto.RegisterType((*Id)(nil), "userAuthPb.Id")
}

func init() { proto.RegisterFile("UserAuth.proto", fileDescriptor_ce91c7f1b4618f2e) }

var fileDescriptor_ce91c7f1b4618f2e = []byte{
	// 122 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x0b, 0x2d, 0x4e, 0x2d,
	0x72, 0x2c, 0x2d, 0xc9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2a, 0x85, 0xf2, 0x03,
	0x92, 0x94, 0x64, 0xb9, 0x58, 0x43, 0xf2, 0xb3, 0x53, 0xf3, 0x84, 0x44, 0xb8, 0x58, 0xcb, 0x12,
	0x73, 0x4a, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x20, 0x1c, 0x25, 0x29, 0x2e, 0x26,
	0xcf, 0x14, 0x54, 0x39, 0x66, 0xa8, 0x9c, 0x91, 0x3d, 0x17, 0x07, 0xcc, 0x60, 0x21, 0x63, 0x2e,
	0x3e, 0xe7, 0x8c, 0xd4, 0xe4, 0x6c, 0x10, 0x07, 0x62, 0x9e, 0xa0, 0x1e, 0xc2, 0x16, 0x3d, 0xb0,
	0x90, 0x14, 0x1f, 0xb2, 0x90, 0x67, 0x4a, 0x12, 0x1b, 0xd8, 0x39, 0xc6, 0x80, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x8f, 0xe1, 0x00, 0x74, 0xa0, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UserAuthClient is the client API for UserAuth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserAuthClient interface {
	CheckAuthToken(ctx context.Context, in *Token, opts ...grpc.CallOption) (*Id, error)
}

type userAuthClient struct {
	cc *grpc.ClientConn
}

func NewUserAuthClient(cc *grpc.ClientConn) UserAuthClient {
	return &userAuthClient{cc}
}

func (c *userAuthClient) CheckAuthToken(ctx context.Context, in *Token, opts ...grpc.CallOption) (*Id, error) {
	out := new(Id)
	err := c.cc.Invoke(ctx, "/userAuthPb.UserAuth/CheckAuthToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserAuthServer is the server API for UserAuth service.
type UserAuthServer interface {
	CheckAuthToken(context.Context, *Token) (*Id, error)
}

// UnimplementedUserAuthServer can be embedded to have forward compatible implementations.
type UnimplementedUserAuthServer struct {
}

func (*UnimplementedUserAuthServer) CheckAuthToken(ctx context.Context, req *Token) (*Id, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAuthToken not implemented")
}

func RegisterUserAuthServer(s *grpc.Server, srv UserAuthServer) {
	s.RegisterService(&_UserAuth_serviceDesc, srv)
}

func _UserAuth_CheckAuthToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAuthServer).CheckAuthToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/userAuthPb.UserAuth/CheckAuthToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAuthServer).CheckAuthToken(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserAuth_serviceDesc = grpc.ServiceDesc{
	ServiceName: "userAuthPb.UserAuth",
	HandlerType: (*UserAuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckAuthToken",
			Handler:    _UserAuth_CheckAuthToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "UserAuth.proto",
}
