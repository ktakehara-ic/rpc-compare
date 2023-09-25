package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"connectrpc.com/connect"
	greetv1 "github.com/ktakehara-icd/rpc-compare/gen/greet/v1"
	"github.com/ktakehara-icd/rpc-compare/gen/greet/v1/greetv1connect"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	logrustrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/sirupsen/logrus"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const addr = "localhost:8081"

var logger *logrus.Logger

func init() {
	tracer.Start()

	logger = logrus.New()
	logger.AddHook(&logrustrace.DDContextLogHook{})
}

func main() {
	greeter := &GreetServer{}
	mux := httptrace.NewServeMux()
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logger.WithContext(r.Context())
			span, ok := tracer.SpanFromContext(r.Context())
			if ok {
				debugSpan(span)
				logger.Printf("span: %+v", span)
			} else {
				logger.Printf("no span")
			}
			logger.Printf("[%s %s] request", r.Method, r.URL)
			defer logger.Printf("[%s %s] response", r.Method, r.URL)
			next.ServeHTTP(w, r)
		})
	}
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)
	mux.Handle(path,
		httptrace.WrapHandler(
			middleware(handler),
			"test-service", "test-resouce",
		),
	)
	http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{}))
}

type GreetServer struct {
}

// Greet implements greetv1connect.GreetServiceHandler.
func (*GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	logger := logger.WithContext(ctx)
	logger.Println("Request headers: ", req.Header())

	span, ok := tracer.SpanFromContext(ctx)
	if ok {
		debugSpan(span)
		logger.Printf("span: %+v", span)
	} else {
		logger.Printf("no span")
	}

	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.GetName()),
	})
	return res, nil
}

var _ (greetv1connect.GreetServiceHandler) = (*GreetServer)(nil)

func debugSpan(span tracer.Span) {
	v := reflect.ValueOf(span)
	logger.Printf("%s", v.Type())
	if v.Type().Kind() == reflect.Pointer {
		v = v.Elem()
	}
	meta := v.FieldByName("Meta")
	m, ok := meta.Interface().(map[string]string)
	if ok {
		l := logger.WithTime(time.Now())
		for k, v := range m {
			l = l.WithField(k, v)
		}
		l.Printf("meta")
	}
	resourceName := v.FieldByName("Resource")
	logger.Printf("%v", resourceName)
}
