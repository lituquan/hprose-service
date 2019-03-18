package hprose.hello.server;

import hprose.config.annotation.RpcService;
import org.springframework.remoting.hprose.OkhttpHproseClient;
import hprose.register.ServiceDiscovery;
import hprose.register.zookeeper.ZooKeeperServiceDiscovery;

public interface IEcho1 {
	public String sayHello1(String name);

	@RpcService(IEcho.class)
	class Client implements IEcho{
		public String sayHello(String name) {
			OkhttpHproseClient client3=new OkhttpHproseClient();
			ServiceDiscovery zooKeeperServiceDiscovery = new ZooKeeperServiceDiscovery("192.168.6.31:31089");
			String discover1 = zooKeeperServiceDiscovery.discover(IEcho1.class.getName());
			IEcho1 iecho = client3.useService(discover1,IEcho1.class);
			return iecho.sayHello1("中国");
		}

	}

	interface IEcho {
		public String sayHello(String name);
	}
}