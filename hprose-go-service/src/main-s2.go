package main

import (
	"hprose-golang/rpc"
	"hprose/hello/server"
	"net/http"
)

const  (
	zookeeper_address  = "localhost:2181"
	zipkinHTTPEndpoint = "http://localhost:9411/api/v1/spans"
)
var (
	helloService       = new(server.Hello)
	service_name       = "service222"
	port               = "8096"
)

func SayHello(name string) string {
	return helloService.SayHello("world")
}
func main() {
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
