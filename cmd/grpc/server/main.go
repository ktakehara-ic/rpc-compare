package main

import (
	"context"
	"fmt"
	"log"
	"net"

	greetv1 "github.com/ktakehara-icd/rpc-compare/gen/greet/v1"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		panic(err)
	}
	log.Printf("server listening at %v", lis.Addr())

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
			log.Printf("request: %s: %v", info.FullMethod, req)
			defer func() {
				if err != nil {
					log.Printf("error: %v", err)
				} else {
					log.Printf("response: %v", err)
				}
			}()
			return handler(ctx, req)
		}),
	)
	greetv1.RegisterGreetServiceServer(srv, &GreetService{})
	srv.Serve(lis)
}

type GreetService struct {
	greetv1.UnimplementedGreetServiceServer
}

// Greet implements greetv1.GreetServiceServer.
func (*GreetService) Greet(ctx context.Context, req *greetv1.GreetRequest) (*greetv1.GreetResponse, error) {
	log.Printf("request: %v", req.GetName())
	return &greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

var _ (greetv1.GreetServiceServer) = (*GreetService)(nil)
