package main

import (
	"fmt"
	"hprose-golang/rpc"
	"hprose/hello/server"
)

//var helloService=new(server.IEcho)
var zookeeper_address = "127.0.0.1:2181"
var helloservice *server.IEcho

func main() {
	rpc.SetServiceName("go-client")
	rpc.InitZipkin()
	rpc.InitRegister(zookeeper_address)
	rpc.Create(&helloservice)
	fmt.Println(helloservice.SayHello("1111"))
}
