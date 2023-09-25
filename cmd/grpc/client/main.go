package main

import (
	"context"
	"log"

	greetv1 "github.com/ktakehara-icd/rpc-compare/gen/greet/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := greetv1.NewGreetServiceClient(conn)

	res, err := c.Greet(context.Background(), &greetv1.GreetRequest{Name: "Jane"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s", res.GetGreeting())
}
