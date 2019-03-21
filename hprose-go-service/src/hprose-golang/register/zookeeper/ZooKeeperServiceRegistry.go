package zookeeper

import (
	"fmt"
	"log"

	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const ZK_SESSION_TIMEOUT = 5000
const ZK_CONNECTION_TIMEOUT = 1000
const ZK_REGISTRY_PATH = "/registry"

/**
 * 基于 ZooKeeper 的服务注册接口实现
 *
 * @author huangyong
 * @since 1.0.0
 */
type ZooKeeperServiceRegistry struct {
	conn *zk.Conn
}
var ZKADDRESS =""
//获取zk连接
func GetZooKeeperServiceRegistry(zkAddress string) *ZooKeeperServiceRegistry {
	c, _, err := zk.Connect([]string{zkAddress}, ZK_CONNECTION_TIMEOUT*time.Millisecond) //*10)
	if err != nil {
		panic(err)
	}
	ZKADDRESS=zkAddress
	return &ZooKeeperServiceRegistry{c}
}

func (s *ZooKeeperServiceRegistry) Discover(name string) (string, error) {
	c, _, err := zk.Connect([]string{ZKADDRESS}, ZK_CONNECTION_TIMEOUT*time.Millisecond) //*10)
	if err != nil {
		panic(err)
	}
	s.conn=c
	path := ZK_REGISTRY_PATH + "/" + name
	// 获取字节点名称
	childs, _, err := s.conn.Children(path)
	if err != nil || len(childs) == 0 {
		if err == zk.ErrNoNode {
			return "", nil
		}
		log.Println(err)
		return "", err
	}
	index := 0
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

func (this *ZooKeeperServiceRegistry) Register(serviceName, serviceAddress string) error {
	// 创建 registry 节点（持久）
	registryPath := ZK_REGISTRY_PATH
	conn := this.conn
	exist, _, err := conn.Exists(registryPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !exist {
		err = ensureRoot(conn)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	// 创建 service 节点（持久）
	servicePath := registryPath + "/" + serviceName
	exist, _, err = conn.Exists(servicePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err := ensureName(servicePath, conn); err != nil {
		return err
	}
	// 创建 address 节点（临时）
	addressPath := servicePath + "/address-"
	data := serviceAddress
	_, err = conn.CreateProtectedEphemeralSequential(addressPath, []byte(data), zk.WorldACL(zk.PermAll))
	return err
}

func ensureName(path string, conn *zk.Conn) error {
	exists, _, err := conn.Exists(path)
	if err != nil {
		return err
	}
	if !exists {
		_, err := conn.Create(path, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func ensureRoot(conn *zk.Conn) error {
	exists, _, err := conn.Exists(ZK_REGISTRY_PATH)
	if err != nil {
		return err
	}
	if !exists {
		_, err := conn.Create(ZK_REGISTRY_PATH, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}
