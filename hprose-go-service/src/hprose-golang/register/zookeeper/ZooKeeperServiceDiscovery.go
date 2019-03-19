package zookeeper

import (
	"github.com/samuel/go-zookeeper/zk"
	"fmt"
	"time"
	"log"
)

/**
 * 基于 ZooKeeper 的服务发现接口实现
 *
 * @author huangyong
 * @since 1.0.0
 */

type ZooKeeperServiceDiscovery struct {
	conn *zk.Conn
}

//获取zk连接
func GetZooKeeperServiceDiscovery(zkAddress string) *ZooKeeperServiceDiscovery {
	c, _, err := zk.Connect([]string{zkAddress}, 10*time.Second) //*10
	if err != nil {
		fmt.Println(err.Error())
	}
	return &ZooKeeperServiceDiscovery{c}
}
func (s *ZooKeeperServiceDiscovery) Discover(name string) (string, error) {
	path := ZK_REGISTRY_PATH + "/" + name
	// 获取字节点名称
	childs, _, err := s.conn.Children(path)
	if err != nil || len(childs)==0{
		if err == zk.ErrNoNode {
			return "", nil
		}
		log.Println(err)
		return  "",err
	}
	index:=0
	//if len(childs)>1{
	//	index=rand.Intn(len(childs))
	//}
	fullPath := path + "/" + childs[index]
	data, _, err := s.conn.Get(fullPath)
	if err != nil {
		panic(err)
	}
	node := string(data)
	return node, nil
}
