package main

import (
	"hprose-golang/rpc"
	"hprose/hello/server"
	"net/http"
	"github.com/afex/hystrix-go/hystrix"
	"net"
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
		Timeout:               1000,
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
	service := rpc.CreateHTTPService()
	serviceObj := &server.Hello2{ //实现接口
		SayHello: SayHello,
	}

	service.AddInstanceMethods(serviceObj)
	//服务注册
	rpc.StartRegister(&serviceObj, port)
	http.ListenAndServe(":"+port, service)
}
