package main

import (
	"fmt"
	"hprose-golang/rpc"
	"hprose/hello/server"
	"github.com/afex/hystrix-go/hystrix"
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
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		Timeout:               2000,
		MaxConcurrentRequests: 100,
		ErrorPercentThreshold: 25,
	})

	output := make(chan string, 1)
	errors := hystrix.Go("my_command", func() error {
		rpc.SetServiceName(service_name, zipkinHTTPEndpoint)
		rpc.InitRegister(zookeeper_address)
		rpc.Create(&helloService)
		output <- helloService.SayHello("123")
		return nil
	}, nil)

	select {
	case out := <-output:
		fmt.Println(out)
	case err := <-errors:
		fmt.Println("get an error, handle it"+err.Error())
	}
}
