package main

import (
	"hprose-golang/rpc"
	"hprose/hello/server"
	"net/http"
)

var helloService = new(server.IEcho)
var zookeeper_address = "127.0.0.1:2181"

func SayHello(name string) string {
	return helloService.SayHello("world")
}
func main() {
	//初始化zipkin、zookeeper、远程服务对象
	rpc.InitRegister(zookeeper_address)
	rpc.Create(&helloService)
	//建立服务对象
	service := rpc.CreateHTTPService()
	serviceObj := &server.Hello{
		SayHello:SayHello,
	}
	port := "8095"
	service.AddInstanceMethods(serviceObj)
	//服务注册
	rpc.StartRegister(&serviceObj, port)
	http.ListenAndServe(":"+port, service)

}
