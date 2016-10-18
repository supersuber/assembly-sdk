// Code generated by protoc-gen-go.
// source: api.proto
// DO NOT EDIT!

/*
Package api is a generated protocol buffer package.

Package api specifies a distributed ledger interface.

This package describes the Ledger API, which defines how users will interact
with the ledger. Various implementations of the actual ledgers can be used,
the basic requirements being strict ordering and append only semantics, the
interface is completely agnostic to things like single-node vs networking or
various degrees of crash or byzantine fault tolerance.

It is generated from these files:
	api.proto

It has these top-level messages:
	ReadRequest
	ReadResult
	SequencedTransaction
	AppendRequest
	UnsequencedTransaction
	AppendResult
	Empty
	ServerStatusResult
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// ReadRequest is a request to read certain transactions from the ledger.
type ReadRequest struct {
	// NetworkSeed identifies the ledger. The read will be rejected if this is
	// set and doesn't match what the ledger has.
	NetworkSeed []byte `protobuf:"bytes,1,opt,name=network_seed,json=networkSeed,proto3" json:"network_seed,omitempty"`
	// Index is the index of the first transactions to read.
	Index int64 `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	// Count is the maximum number of transactions to read (if available).
	Count int64 `protobuf:"varint,3,opt,name=count" json:"count,omitempty"`
}

func (m *ReadRequest) Reset()                    { *m = ReadRequest{} }
func (m *ReadRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()               {}
func (*ReadRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// ReadResult is the result of a ReadTransactions call.
type ReadResult struct {
	// NetworkSeed identifies the ledger. It will always stay the same for a
	// given ledger; if it has a surprising value, transactions were read from
	// a different (or potentially reset) ledger.
	NetworkSeed []byte `protobuf:"bytes,1,opt,name=network_seed,json=networkSeed,proto3" json:"network_seed,omitempty"`
	// Transactions are sequenced transactions read from the ledger.
	Transactions []*SequencedTransaction `protobuf:"bytes,2,rep,name=transactions" json:"transactions,omitempty"`
}

func (m *ReadResult) Reset()                    { *m = ReadResult{} }
func (m *ReadResult) String() string            { return proto.CompactTextString(m) }
func (*ReadResult) ProtoMessage()               {}
func (*ReadResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ReadResult) GetTransactions() []*SequencedTransaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

// SequencedTransaction is a transaction that has been appended to the ledger
// and assigned index, timestamp and state hash. It's unique in its index and
// its fields are guaranteed to never change and be equal on all honest ledger
// nodes.
type SequencedTransaction struct {
	// Type is the transaction type. This field allows clients to filter out
	// transactions they don't care about.
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	// Index is the transactions index on the ledger. It's unique and without
	// gaps. Counts from 1 (<= 0 are invalid values).
	Index int64 `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	// Timestamp is time this transaction was appended to the ledger.
	// Monotonically increasing with Index. Accuracy is implementation
	// specific, unit is nanoseconds since the Unix epoch.
	Timestamp int64 `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
	// Data is the payload of the transaction. The ledger treats this as
	// arbitrary, unparseable binary data.
	Data []byte `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	// Hash is the SHA256 hash of the concatenation of type and data.
	Hash []byte `protobuf:"bytes,5,opt,name=hash,proto3" json:"hash,omitempty"`
	// StateHash is the SHA256 hash of the concatenation of the previous state
	// hash and the hash of this transaction. For the first transaction it will
	// be simply its hash. This provides paranoid clients with a mechanism for
	// verifying ledger correctness.
	//
	// `StateHash = sha256.Sum256(append(previousStateHash, Hash...))`
	StateHash []byte `protobuf:"bytes,6,opt,name=state_hash,json=stateHash,proto3" json:"state_hash,omitempty"`
}

func (m *SequencedTransaction) Reset()                    { *m = SequencedTransaction{} }
func (m *SequencedTransaction) String() string            { return proto.CompactTextString(m) }
func (*SequencedTransaction) ProtoMessage()               {}
func (*SequencedTransaction) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

// AppendRequest contains transaction to append to the ledger.
type AppendRequest struct {
	// NetworkSeed identifies the ledger. The request will be rejected if this
	// is set and doesn't match what the ledger has.
	NetworkSeed []byte `protobuf:"bytes,1,opt,name=network_seed,json=networkSeed,proto3" json:"network_seed,omitempty"`
	// Transactions are transactions to be appended to the ledger in an
	// unspecified order.
	Transactions []*UnsequencedTransaction `protobuf:"bytes,2,rep,name=transactions" json:"transactions,omitempty"`
}

func (m *AppendRequest) Reset()                    { *m = AppendRequest{} }
func (m *AppendRequest) String() string            { return proto.CompactTextString(m) }
func (*AppendRequest) ProtoMessage()               {}
func (*AppendRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AppendRequest) GetTransactions() []*UnsequencedTransaction {
	if m != nil {
		return m.Transactions
	}
	return nil
}

// UnsequencedTransaction is a transaction that hasn't been assigned an index
// on the ledger yet. No assumptions should be made about when or whether it
// will be appended to the ledger.
type UnsequencedTransaction struct {
	// Type is the transaction type. This field allows clients to filter out
	// transactions they don't care about. It can also be used to pass
	// administration messages to ledger nodes.
	Type string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	// Data is the payload of the transaction. The ledger treats this as
	// arbitrary, unparseable binary data.
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// Hash is the SHA256 hash of the concatenation of type and data. The
	// ledger may fill this in if empty, or reject the transaction if it
	// doesn't match.
	Hash []byte `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (m *UnsequencedTransaction) Reset()                    { *m = UnsequencedTransaction{} }
func (m *UnsequencedTransaction) String() string            { return proto.CompactTextString(m) }
func (*UnsequencedTransaction) ProtoMessage()               {}
func (*UnsequencedTransaction) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

// AppendResult is the result of an append call.
type AppendResult struct {
	// NetworkSeed identifies the ledger. It will always stay the same for a
	// given ledger; if it has a surprising value, transactions were attempted
	// appended to a different (or potentially reset) ledger. Such attempts
	// will be rejected.
	NetworkSeed []byte `protobuf:"bytes,1,opt,name=network_seed,json=networkSeed,proto3" json:"network_seed,omitempty"`
	// LastIndex is the index assigned by the ledger to the last of the
	// provided transactions. Any subsequent appends are guaranteed to be order
	// after this index.
	LastIndex int64 `protobuf:"varint,2,opt,name=last_index,json=lastIndex" json:"last_index,omitempty"`
}

func (m *AppendResult) Reset()                    { *m = AppendResult{} }
func (m *AppendResult) String() string            { return proto.CompactTextString(m) }
func (*AppendResult) ProtoMessage()               {}
func (*AppendResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

// Empty as an empty message.
type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

// ServerStatus is the status of the local ledger node.
type ServerStatusResult struct {
	// NetworkSeed identifies the ledger. It will always stay the same for a
	// given ledger; if it has an unexpected value, the request was handled by
	// a different (or potentially reset) ledger.
	NetworkSeed []byte `protobuf:"bytes,1,opt,name=network_seed,json=networkSeed,proto3" json:"network_seed,omitempty"`
	// NetworkType is an arbitrary string describing the type of the ledger
	// (eg. "production" or "testing").
	NetworkType string `protobuf:"bytes,2,opt,name=network_type,json=networkType" json:"network_type,omitempty"`
	// LastIndex is the index of the last transaction appended to the ledger. A
	// low number indicates that the local node is behind the rest of the
	// network.
	LastIndex int64 `protobuf:"varint,3,opt,name=last_index,json=lastIndex" json:"last_index,omitempty"`
	// ServerTime is the time as seen by the local ledger node, in nanoseconds
	// since the Unix epoch.
	ServerTime int64 `protobuf:"varint,4,opt,name=server_time,json=serverTime" json:"server_time,omitempty"`
	// Ready is a flag indicating if the local node deems itself ready to
	// handle read and append requests. It can be false if the node is in the
	// process of catching up to the rest of the network or is experiencing
	// some other issue.
	Ready bool `protobuf:"varint,5,opt,name=ready" json:"ready,omitempty"`
}

func (m *ServerStatusResult) Reset()                    { *m = ServerStatusResult{} }
func (m *ServerStatusResult) String() string            { return proto.CompactTextString(m) }
func (*ServerStatusResult) ProtoMessage()               {}
func (*ServerStatusResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func init() {
	proto.RegisterType((*ReadRequest)(nil), "api.ReadRequest")
	proto.RegisterType((*ReadResult)(nil), "api.ReadResult")
	proto.RegisterType((*SequencedTransaction)(nil), "api.SequencedTransaction")
	proto.RegisterType((*AppendRequest)(nil), "api.AppendRequest")
	proto.RegisterType((*UnsequencedTransaction)(nil), "api.UnsequencedTransaction")
	proto.RegisterType((*AppendResult)(nil), "api.AppendResult")
	proto.RegisterType((*Empty)(nil), "api.Empty")
	proto.RegisterType((*ServerStatusResult)(nil), "api.ServerStatusResult")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Ledger service

type LedgerClient interface {
	// ReadTransactions reads the requested transactions from the ledger. If no
	// data is immediately available the request will stall up to the provided
	// deadline before returning an empty object. If data becomes available
	// while waiting it will immediately be returned.
	ReadTransactions(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResult, error)
	// AppendTransactions append transactions to the ledger in an unspecified
	// order. Multiple simultaneous append calls will similarly result in
	// unspecified ordering between the calls. Once the append succeeds,
	// transactions in subsequent appends are guaranteed to be order strictly
	// after those in this append.
	AppendTransactions(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResult, error)
	// ServerStatus returns info about and status of the local ledger node.
	ServerStatus(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ServerStatusResult, error)
}

type ledgerClient struct {
	cc *grpc.ClientConn
}

func NewLedgerClient(cc *grpc.ClientConn) LedgerClient {
	return &ledgerClient{cc}
}

func (c *ledgerClient) ReadTransactions(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResult, error) {
	out := new(ReadResult)
	err := grpc.Invoke(ctx, "/api.Ledger/ReadTransactions", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) AppendTransactions(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResult, error) {
	out := new(AppendResult)
	err := grpc.Invoke(ctx, "/api.Ledger/AppendTransactions", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ledgerClient) ServerStatus(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ServerStatusResult, error) {
	out := new(ServerStatusResult)
	err := grpc.Invoke(ctx, "/api.Ledger/ServerStatus", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Ledger service

type LedgerServer interface {
	// ReadTransactions reads the requested transactions from the ledger. If no
	// data is immediately available the request will stall up to the provided
	// deadline before returning an empty object. If data becomes available
	// while waiting it will immediately be returned.
	ReadTransactions(context.Context, *ReadRequest) (*ReadResult, error)
	// AppendTransactions append transactions to the ledger in an unspecified
	// order. Multiple simultaneous append calls will similarly result in
	// unspecified ordering between the calls. Once the append succeeds,
	// transactions in subsequent appends are guaranteed to be order strictly
	// after those in this append.
	AppendTransactions(context.Context, *AppendRequest) (*AppendResult, error)
	// ServerStatus returns info about and status of the local ledger node.
	ServerStatus(context.Context, *Empty) (*ServerStatusResult, error)
}

func RegisterLedgerServer(s *grpc.Server, srv LedgerServer) {
	s.RegisterService(&_Ledger_serviceDesc, srv)
}

func _Ledger_ReadTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).ReadTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/ReadTransactions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).ReadTransactions(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_AppendTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).AppendTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/AppendTransactions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).AppendTransactions(ctx, req.(*AppendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ledger_ServerStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServer).ServerStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Ledger/ServerStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServer).ServerStatus(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Ledger_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Ledger",
	HandlerType: (*LedgerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReadTransactions",
			Handler:    _Ledger_ReadTransactions_Handler,
		},
		{
			MethodName: "AppendTransactions",
			Handler:    _Ledger_AppendTransactions_Handler,
		},
		{
			MethodName: "ServerStatus",
			Handler:    _Ledger_ServerStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 441 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x94, 0xc1, 0x6f, 0x9b, 0x30,
	0x14, 0xc6, 0x45, 0x48, 0xb2, 0xf1, 0x60, 0x5a, 0x66, 0x45, 0x1b, 0xcb, 0x36, 0x6d, 0xe3, 0x94,
	0x53, 0x0e, 0x89, 0x76, 0x9a, 0xa6, 0x69, 0x87, 0x49, 0xad, 0xd4, 0x43, 0x65, 0xe8, 0xad, 0x12,
	0x72, 0x83, 0xd5, 0xa0, 0x26, 0x86, 0x62, 0xd3, 0x36, 0x7f, 0x50, 0xef, 0xed, 0x7f, 0x58, 0xfb,
	0x91, 0x28, 0x10, 0x51, 0x35, 0xbd, 0xd9, 0x3f, 0x7f, 0xb6, 0xdf, 0xf7, 0x3d, 0x03, 0x38, 0x2c,
	0x4f, 0x27, 0x79, 0x91, 0xa9, 0x8c, 0xd8, 0x7a, 0x18, 0x9c, 0x83, 0x4b, 0x39, 0x4b, 0x28, 0xbf,
	0x2e, 0xb9, 0x54, 0xe4, 0x27, 0x78, 0x82, 0xab, 0xdb, 0xac, 0xb8, 0x8a, 0x25, 0xe7, 0x89, 0x6f,
	0xfd, 0xb0, 0xc6, 0x1e, 0x75, 0x37, 0x2c, 0xd4, 0x88, 0x0c, 0xa1, 0x97, 0x8a, 0x84, 0xdf, 0xf9,
	0x1d, 0xbd, 0x66, 0xd3, 0x6a, 0x62, 0xe8, 0x3c, 0x2b, 0x85, 0xf2, 0xed, 0x8a, 0xe2, 0x24, 0x10,
	0x00, 0xd5, 0xe9, 0xb2, 0x5c, 0x1e, 0x74, 0xf8, 0x1f, 0xf0, 0x54, 0xc1, 0x84, 0x64, 0x73, 0x95,
	0x66, 0x42, 0xea, 0x3b, 0xec, 0xb1, 0x3b, 0xfd, 0x3c, 0x31, 0x55, 0x87, 0xa6, 0x46, 0x31, 0xe7,
	0x49, 0xb4, 0x53, 0xd0, 0x86, 0x3c, 0xb8, 0xb7, 0x60, 0xd8, 0x26, 0x23, 0x04, 0xba, 0x6a, 0x9d,
	0x73, 0xbc, 0xd2, 0xa1, 0x38, 0x7e, 0xc6, 0xc8, 0x57, 0x70, 0x54, 0xba, 0xd2, 0x59, 0xb0, 0x55,
	0xbe, 0x31, 0xb3, 0x03, 0xe6, 0x9c, 0x84, 0x29, 0xe6, 0x77, 0xb1, 0x74, 0x1c, 0x1b, 0xb6, 0x60,
	0x72, 0xe1, 0xf7, 0x2a, 0x66, 0xc6, 0xe4, 0x1b, 0x80, 0xde, 0xa0, 0x78, 0x8c, 0x2b, 0x7d, 0x5c,
	0x71, 0x90, 0x1c, 0x69, 0x10, 0x48, 0x78, 0xf7, 0x2f, 0xcf, 0xb9, 0x78, 0x4d, 0xee, 0x7f, 0x5b,
	0xa3, 0xf9, 0x82, 0xd1, 0x9c, 0x09, 0xf9, 0x72, 0x38, 0x11, 0x7c, 0x6c, 0xd7, 0xb5, 0xa6, 0xb3,
	0x75, 0xda, 0x69, 0x71, 0x6a, 0xef, 0x9c, 0x06, 0xa7, 0xe0, 0x6d, 0xad, 0x1c, 0xda, 0x64, 0x1d,
	0xce, 0x92, 0x49, 0x15, 0xd7, 0xd3, 0x77, 0x0c, 0x39, 0x36, 0x20, 0x78, 0x03, 0xbd, 0xff, 0xab,
	0x5c, 0xad, 0x83, 0x07, 0x0b, 0x48, 0xc8, 0x8b, 0x1b, 0x5e, 0x84, 0x3a, 0xb9, 0x52, 0x1e, 0x7e,
	0x43, 0x4d, 0x82, 0xc6, 0x3a, 0x68, 0x6c, 0x2b, 0x89, 0x8c, 0xbf, 0x66, 0x11, 0xf6, 0x5e, 0x11,
	0xe4, 0x3b, 0xb8, 0x12, 0xaf, 0x8e, 0x4d, 0xf3, 0xb1, 0xdf, 0x36, 0x85, 0x0a, 0x45, 0x9a, 0x98,
	0xd7, 0x53, 0xe8, 0xa7, 0xbd, 0xc6, 0xb6, 0xbf, 0xa5, 0xd5, 0x64, 0xfa, 0x68, 0x41, 0xff, 0x84,
	0x27, 0x97, 0xbc, 0x20, 0xbf, 0x60, 0x60, 0xde, 0x7e, 0x2d, 0x67, 0x49, 0x06, 0xd8, 0xad, 0xda,
	0x07, 0x37, 0x7a, 0x5f, 0x23, 0xe8, 0xee, 0x37, 0x90, 0x2a, 0xcf, 0xc6, 0x46, 0x82, 0xb2, 0xc6,
	0x9b, 0x19, 0x7d, 0x68, 0x30, 0xdc, 0x3c, 0x03, 0xaf, 0x1e, 0x18, 0x01, 0x94, 0x60, 0x9a, 0xa3,
	0x4f, 0x9b, 0x8f, 0x68, 0x3f, 0xcf, 0x8b, 0x3e, 0xfe, 0x0e, 0x66, 0x4f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xb9, 0xa2, 0xcb, 0xd4, 0x1b, 0x04, 0x00, 0x00,
}
