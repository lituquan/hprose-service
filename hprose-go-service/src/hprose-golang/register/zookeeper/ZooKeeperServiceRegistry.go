package zookeeper

import (
	"fmt"

	"github.com/samuel/go-zookeeper/zk"
	"time"
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
	ZkClient *zk.Conn
}

//获取zk连接
func GetZooKeeperServiceRegistry(zkAddress string) *ZooKeeperServiceRegistry {
	c, _, err := zk.Connect([]string{zkAddress}, ZK_CONNECTION_TIMEOUT*time.Millisecond) //*10)
	if err != nil {
		panic(err)
	}
	return &ZooKeeperServiceRegistry{c}
}

func (this *ZooKeeperServiceRegistry) Register(serviceName, serviceAddress string) error {
	// 创建 registry 节点（持久）
	registryPath := ZK_REGISTRY_PATH
	zkClient := this.ZkClient
	exist, _, err := zkClient.Exists(registryPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !exist {
		err = ensureRoot(zkClient)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	// 创建 service 节点（持久）
	servicePath := registryPath + "/" + serviceName
	exist, _, err = zkClient.Exists(servicePath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if err := ensureName(servicePath, zkClient); err != nil {
		return err
	}
	// 创建 address 节点（临时）
	addressPath := servicePath + "/address-"
	data := serviceAddress
	_, err = zkClient.CreateProtectedEphemeralSequential(addressPath, []byte(data), zk.WorldACL(zk.PermAll))
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
