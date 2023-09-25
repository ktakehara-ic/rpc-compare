package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"connectrpc.com/connect"
	greetv1 "github.com/ktakehara-icd/rpc-compare/gen/greet/v1"
	"github.com/ktakehara-icd/rpc-compare/gen/greet/v1/greetv1connect"
	"golang.org/x/net/http2"
)

func main() {
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8081",
		connect.WithGRPC(),
	)
	res, err := client.Greet(
		context.Background(),
		connect.NewRequest(&greetv1.GreetRequest{Name: "Jane"}),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.Msg.GetGreeting())
}

func insecureGRPCClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
}
