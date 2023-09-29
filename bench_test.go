package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	greetv1 "github.com/ktakehara-icd/rpc-compare/gen/greet/v1"
	"github.com/ktakehara-icd/rpc-compare/gen/greet/v1/greetv1connect"
	"github.com/stretchr/testify/assert"
	"github.com/twitchtv/twirp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

var (
	insecureHTTPClient = &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	smallReq = &greetv1.GreetRequest{
		Name: "Jane",
		Age:  34,
		Address: &greetv1.GreetRequest_Address{
			PostalCode:    "0000000",
			StateProvince: "Foo State",
			City:          "Foo City",
			Street:        "Bar St.",
			BuildingName:  "BLDG 1234",
			Note: []string{
				"This is example address.",
			},
		},
	}
	bigReq = &greetv1.GreetRequest{
		Name: "Jane",
		Age:  34,
		Address: &greetv1.GreetRequest_Address{
			PostalCode:    "0000000",
			StateProvince: "Foo State",
			City:          "Foo City",
			Street:        "Bar St.",
			BuildingName:  "BLDG 1234",
			Note: []string{
				"This is example address.",
			},
		},
	}
)

func init() {
	for i := 0; i < 10000; i++ {
		bigReq.Appendix = append(bigReq.Appendix, fmt.Sprintf("This is appendix text: %d", i))
	}
}

func BenchmarkAll(b *testing.B) {
	grpcLis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		b.Fatalf("cannot listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	greetv1.RegisterGreetServiceServer(grpcSrv, &GreetServiceGRPC{})
	go grpcSrv.Serve(grpcLis)

	connectMux := http.NewServeMux()
	connectMux.Handle(greetv1connect.NewGreetServiceHandler(&GreetServiceConnect{}))
	twirpH := greetv1.NewGreetServiceServer(&GreetServiceTwirp{})
	mux := http.NewServeMux()
	mux.Handle(greetv1connect.NewGreetServiceHandler(&GreetServiceConnect{})) // for gRPC client
	mux.Handle("/connect/", http.StripPrefix("/connect", connectMux))
	mux.Handle(twirpH.PathPrefix(), twirpH)
	mux.HandleFunc("/rest", RESTHandler)
	srv := httptest.NewUnstartedServer(h2c.NewHandler(mux, &http2.Server{}))
	srv.EnableHTTP2 = true
	srv.Start()
	b.Cleanup(srv.Close)

	type Client interface {
		Greet(context.Context, *greetv1.GreetRequest) (*greetv1.GreetResponse, error)
	}

	// name is "server-client"
	tests := []struct {
		name   string
		client Client
	}{
		{
			name:   "gRPC-gRPC",
			client: &wrapGRPCClient{raw: newGRPCClient(b, grpcLis)},
		},
		{
			name: "gRPC-connect",
			client: &wrapConnectClient{raw: greetv1connect.NewGreetServiceClient(
				insecureHTTPClient, "http://"+grpcLis.Addr().String(), connect.WithGRPC())},
		},
		{
			name:   "connect-gRPC",
			client: &wrapGRPCClient{raw: newGRPCClient(b, srv.Listener)},
		},
		{
			name:   "connect-connect",
			client: &wrapConnectClient{raw: greetv1connect.NewGreetServiceClient(srv.Client(), srv.URL+"/connect")},
		},
		{
			name: "connect-twirp(json)",
			client: greetv1.NewGreetServiceJSONClient(srv.URL, srv.Client(),
				twirp.WithClientPathPrefix("connect")),
		},
		{
			name: "connect-twirp(proto)",
			client: greetv1.NewGreetServiceProtobufClient(srv.URL, srv.Client(),
				twirp.WithClientPathPrefix("connect"),
				twirp.WithClientHooks(&twirp.ClientHooks{
					RequestPrepared: func(ctx context.Context, r *http.Request) (context.Context, error) {
						r.Header.Set("Content-Type", "application/proto")
						return ctx, nil
					},
				})),
		},
		{
			name:   "twirp-twirp(json)",
			client: greetv1.NewGreetServiceJSONClient(srv.URL, srv.Client()),
		},
		{
			name:   "twirp-twirp(proto)",
			client: greetv1.NewGreetServiceProtobufClient(srv.URL, srv.Client()),
		},
		{
			name:   "twirp-connect(json)",
			client: &wrapConnectClient{raw: greetv1connect.NewGreetServiceClient(srv.Client(), srv.URL+"/twirp", connect.WithProtoJSON())},
		},
		{
			name:   "twirp-connect(proto)",
			client: &wrapConnectClient{raw: greetv1connect.NewGreetServiceClient(srv.Client(), srv.URL+"/twirp", connect.WithCodec(&protobufBinaryCodec{}))},
		},
		{
			name:   "REST-REST",
			client: NewRESTClient(srv.Client(), srv.URL+"/rest"),
		},
	}
	for _, reqs := range []struct {
		name string
		req  *greetv1.GreetRequest
	}{
		{name: "small", req: smallReq}, {name: "big", req: bigReq},
	} {
		for _, tt := range tests {
			b.Run(reqs.name+"/"+tt.name, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					resp, err := tt.client.Greet(context.Background(), reqs.req)
					_ = assert.NoError(b, err, "request error") &&
						assert.Equal(b, "Hello, Jane!", resp.GetGreeting(), "response mismatch")
				}
			})
		}
	}
}

type wrapGRPCClient struct {
	raw  greetv1.GreetServiceClient
	opts []grpc.CallOption
}

func (c *wrapGRPCClient) Greet(ctx context.Context, req *greetv1.GreetRequest) (*greetv1.GreetResponse, error) {
	return c.raw.Greet(ctx, req, c.opts...)
}

type wrapConnectClient struct {
	raw greetv1connect.GreetServiceClient
}

func (c *wrapConnectClient) Greet(ctx context.Context, req *greetv1.GreetRequest) (*greetv1.GreetResponse, error) {
	resp, err := c.raw.Greet(ctx, connect.NewRequest(req))
	if err != nil {
		return nil, err
	}
	return resp.Msg, nil
}

func newGRPCClient(b *testing.B, l net.Listener) greetv1.GreetServiceClient {
	b.Helper()

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		b.Fatalf("cannot create gRPC connection: %v", err)
	}
	b.Cleanup(func() { conn.Close() })

	return greetv1.NewGreetServiceClient(conn)
}

// copy from https://github.com/connectrpc/connect-go/blob/065b6ad19e9243bb66472c20f88504fc3c8e9fbb/codec.go
// Difference from original code is only Name() returning value.
type protobufBinaryCodec struct{}

var _ connect.Codec = (*protobufBinaryCodec)(nil)

func (c *protobufBinaryCodec) Name() string { return "protobuf" }

func (c *protobufBinaryCodec) Marshal(message any) ([]byte, error) {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return nil, errNotProto(message)
	}
	return proto.Marshal(protoMessage)
}

func (c *protobufBinaryCodec) MarshalAppend(dst []byte, message any) ([]byte, error) {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return nil, errNotProto(message)
	}
	return proto.MarshalOptions{}.MarshalAppend(dst, protoMessage)
}

func (c *protobufBinaryCodec) Unmarshal(data []byte, message any) error {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return errNotProto(message)
	}
	return proto.Unmarshal(data, protoMessage)
}

func (c *protobufBinaryCodec) MarshalStable(message any) ([]byte, error) {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return nil, errNotProto(message)
	}
	// protobuf does not offer a canonical output today, so this format is not
	// guaranteed to match deterministic output from other protobuf libraries.
	// In addition, unknown fields may cause inconsistent output for otherwise
	// equal messages.
	// https://github.com/golang/protobuf/issues/1121
	options := proto.MarshalOptions{Deterministic: true}
	return options.Marshal(protoMessage)
}

func (c *protobufBinaryCodec) IsBinary() bool {
	return true
}

func errNotProto(message any) error {
	if _, ok := message.(protoiface.MessageV1); ok {
		return fmt.Errorf("%T uses github.com/golang/protobuf, but connect-go only supports google.golang.org/protobuf: see https://go.dev/blog/protobuf-apiv2", message)
	}
	return fmt.Errorf("%T doesn't implement proto.Message", message)
}
