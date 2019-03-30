package rpc

import (
	"fmt"
	"hprose-golang/register/zookeeper"
	"net"
	"os"
	"reflect"
)

var register *zookeeper.ZooKeeperServiceRegistry

func InitRegister(zkaddress string) {
	register = zookeeper.GetZooKeeperServiceRegistry(zkaddress)
}

func StartRegister(serviceObj interface{}, port string) {
	v := reflect.ValueOf(serviceObj).Elem()
	mv := v.MethodByName("ServiceName")
	results := mv.Call(nil)
	serviceName := results[0].String()
	addr := fmt.Sprintf("http://"+service_name)
	err := register.Register(serviceName, addr)
	if err != nil {
		panic(err)
	}
}

func getIntranetIp() string {
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
