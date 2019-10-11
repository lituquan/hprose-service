package main

import (
	"github.com/afex/hystrix-go/hystrix"
	hrpc "github.com/hprose/hprose-golang/rpc"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin-contrib/zipkin-go-opentracing/examples/middleware"
	"hprose-golang/rpc"
	"hprose/hello/server"
	"log"
	"net"
	"net/http"
)

const  (
	zookeeper_address  = "localhost:2181"
	zipkinHTTPEndpoint = "http://localhost:9411/api/v1/spans"
)
var (
	helloService       = new(server.Hello)
	service_name       = "127.0.0.1:8096"
	port               = "8096"
)

func SayHello(name string) string {
	hystrix.ConfigureCommand("SayHello", hystrix.CommandConfig{
		Timeout:               10000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})
	out:=""
	hystrix.Do("SayHello", func() error {
		rpc.SetServiceName(service_name, zipkinHTTPEndpoint)
		rpc.InitRegister(zookeeper_address)
		rpc.Create(&helloService)
		out=helloService.SayHello("123")
		return nil
	}, func(e error) error {
		log.Println(e.Error())
		out= "default"
		return  nil
	})
	return out
}
func main() {
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)

	//初始化zipkin、zookeeper、远程服务对象
	rpc.SetServiceName(service_name, zipkinHTTPEndpoint)
	rpc.InitRegister(zookeeper_address)
	rpc.Create(&helloService)
	//建立服务对象
	service := hrpc.NewHTTPService()
	serviceObj := &server.Hello2{ //实现接口
		SayHello: SayHello,
	}

	service.AddInstanceMethods(serviceObj)
	//服务注册
	rpc.StartRegister(&serviceObj, port)

	_,tracer:=rpc.InitZipkin()
	// create the HTTP Server Handler for the service
	handler := NewHTTPHandler(tracer, service)
	http.ListenAndServe(":"+port, handler)
}

// NewHTTPHandler returns a new HTTP handler our svc2.
func NewHTTPHandler(tracer opentracing.Tracer, service *hrpc.HTTPService) http.Handler {
	// Create the mux.
	mux := http.NewServeMux()
	// Create the Sum handler.
	var sumHandler http.Handler
	sumHandler = http.HandlerFunc(service.ServeHTTP)
	// Wrap the Sum handler with our tracing middleware.
	sumHandler = middleware.FromHTTPRequest(tracer, "Sum")(sumHandler)
	// Wire up the mux.
	mux.Handle("/", sumHandler)
	// Return the mux.
	return mux
}