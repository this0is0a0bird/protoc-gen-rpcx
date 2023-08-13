// Code generated by david client code gen. DO NOT EDIT.
// versions:
// - protoc-gen-rpcx v0.3.0
// - protoc          v3.17.0
// source: helloworld.proto

package greeter_service

import (
	context "context"
	rpcclient "github.com/cctip/cctip-service-client/rpcclient"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = context.TODO
var _ = rpcclient.GetRpcClient

// ================== interface skeleton ===================
type GreeterAble interface {
	// GreeterAble can be used for interface verification.

	// SayHello is server rpc method as defined
	SayHello(ctx context.Context, args *HelloRequest, reply *HelloReply) (err error)
}

// ================== client stub ===================
var (
	_greeterClient *greeterClient
	_clientOnce    sync.Once
)

// Greeter is a client wrapped XClient.
type greeterClient struct {
}

// NewGreeterClient wraps a XClient as greeterClient.
// You can pass a shared XClient object created by NewXClientForGreeter.
func GetGreeterClient() *greeterClient {
	if _greeterClient == nil {
		_clientOnce.Do(func() {
			_greeterClient = &greeterClient{}
		})
	}
	return _greeterClient
}

// SayHello is client rpc method as defined
func (c *greeterClient) SayHello(ctx context.Context, args *HelloRequest) (reply *HelloReply, err error) {
	reply = &HelloReply{}
	err = rpcclient.GetRpcClient().Call(ctx, ServiceName, "SayHello", args, reply)
	return reply, err
}
