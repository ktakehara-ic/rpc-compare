package main

import (
	"context"
	"fmt"

	greetv1 "github.com/ktakehara-icd/rpc-compare/gen/greet/v1"
)

type GreetServiceGRPC struct {
	greetv1.UnimplementedGreetServiceServer
}

// Greet implements greetv1.GreetServiceServer.
func (*GreetServiceGRPC) Greet(_ context.Context, req *greetv1.GreetRequest) (*greetv1.GreetResponse, error) {
	return &greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

var _ (greetv1.GreetServiceServer) = (*GreetServiceGRPC)(nil)
