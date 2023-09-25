package main

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	greetv1 "github.com/ktakehara-icd/rpc-compare/gen/greet/v1"
	"github.com/ktakehara-icd/rpc-compare/gen/greet/v1/greetv1connect"
)

type GreetServiceConnect struct {
}

// Greet implements greetv1connect.GreetServiceHandler.
func (*GreetServiceConnect) Greet(_ context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	return connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.GetName()),
	}), nil
}

var _ (greetv1connect.GreetServiceHandler) = (*GreetServiceConnect)(nil)
