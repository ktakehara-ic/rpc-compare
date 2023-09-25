package main

import (
	"context"
	"fmt"

	greetv1 "github.com/ktakehara-icd/rpc-compare/gen/greet/v1"
)

type GreetServiceTwirp struct{}

// Greet implements greetv1.GreetService.
func (*GreetServiceTwirp) Greet(_ context.Context, req *greetv1.GreetRequest) (*greetv1.GreetResponse, error) {
	return &greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

var _ (greetv1.GreetService) = (*GreetServiceTwirp)(nil)
