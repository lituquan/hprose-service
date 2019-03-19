package main

import (
	"hprose/hello/server"
	"hprose-golang/rpc"
	"fmt"
	"hprose-golang/register/zookeeper"
	"reflect"
	"os"
	"net"
)
var helloService=new(server.IEcho)
var zookeeper_address="127.0.0.1:2181"

func main() {
	var  helloservice *hello
	rpc.SetServiceName("go-client")
	rpc.InitRegister(zookeeper_address)
	rpc.Create(&helloservice)
	fmt.Println(helloservice.ServiceName())
}

func GetIntranetIp() string{
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("ip:", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}
	panic("没有ip")
	return ""
}
type hello struct{

}
func (*hello)ServiceName() string{
	return "hprose.hello.server.go.IEcho"
}

func (*hello)SayHello111(name string) string {
	rpc.InitRegister(zookeeper_address)
	rpc.Create(&helloService)
	return helloService.SayHello("world")
}

func register(serviceObj interface{},port string){
	register:=zookeeper.GetZooKeeperServiceRegistry(zookeeper_address)
	v:=reflect.ValueOf(serviceObj).Elem()
	mv := v.MethodByName("ServiceName")
	results:=mv.Call(nil)
	serviceName:=results[0].String()
	addr:=fmt.Sprintf("http://"+GetIntranetIp()+":%s",port)
	register.Register(serviceName,addr)
}