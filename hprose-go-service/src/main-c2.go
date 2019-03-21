package main

import (
	"fmt"
	"hprose-golang/rpc"
	"hprose/hello/server"
)


const  (
	zookeeper_address  = "localhost:2181"
	zipkinHTTPEndpoint = "http://localhost:9411/api/v1/spans"
)
var (
	helloService       = new(server.Hello2)
	service_name       = "go-client-2"
	port               = "8096"
)

func main() {
	rpc.SetServiceName(service_name,zipkinHTTPEndpoint)
	rpc.InitRegister(zookeeper_address)
	rpc.Create(&helloService)
	fmt.Println(helloService.SayHello("1111"))
}
