package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	greetv1 "github.com/ktakehara-icd/rpc-compare/gen/greet/v1"
)

func RESTHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user greetv1.GreetRequest
	decoder.Decode(&user)
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", user.Name),
	})
}

type RESTClient struct {
	client  *http.Client
	baseURL string
}

func NewRESTClient(client *http.Client, baseURL string) *RESTClient {
	return &RESTClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *RESTClient) Greet(ctx context.Context, req *greetv1.GreetRequest) (*greetv1.GreetResponse, error) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(req)

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, &buf)
	if err != nil {
		return nil, err
	}

	hresp, err := c.client.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	defer hresp.Body.Close()

	// We need to parse response to have a fair comparison as gRPC does it
	var resp greetv1.GreetResponse
	decodeErr := json.NewDecoder(hresp.Body).Decode(&resp)
	if decodeErr != nil {
		return nil, fmt.Errorf("unable to decode json: %w", decodeErr)
	}
	return &resp, nil
}
