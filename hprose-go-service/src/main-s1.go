package main

import (
	"hprose-golang/rpc"
	"hprose/hello/server"
	"net/http"
	"time"
)

const  (
	zookeeper_address  = "localhost:2181"
	zipkinHTTPEndpoint = "http://localhost:9411/api/v1/spans"
)
var (
	helloService       = new(server.IEcho)
	service_name       = "127.0.0.1:8095"
	port               = "8095"
)

func SayHello(name string) string {
	//return helloService.SayHello("world")
	time.Sleep(1*time.Second)
	return "2222"
}
func main() {
	//初始化zipkin、zookeeper、远程服务对象
	rpc.SetServiceName(service_name, zipkinHTTPEndpoint)
	rpc.InitRegister(zookeeper_address)
	//rpc.Create(&helloService)
	//建立服务对象
	service := rpc.CreateHTTPService()
	serviceObj := &server.Hello{ //实现接口
		SayHello: SayHello,
	}

	service.AddInstanceMethods(serviceObj)
	//服务注册
	rpc.StartRegister(&serviceObj, port)
	http.ListenAndServe(":"+port, service)
}
