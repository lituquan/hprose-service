package hprose.hello.server;

import hprose.config.annotation.RpcService;
import org.springframework.remoting.hprose.OkhttpHproseClient;
import hprose.register.ServiceRegistry;
import hprose.register.zookeeper.ZooKeeperServiceRegistry;

public interface IEcho1 {
	public String sayHello1(String name);

	@RpcService(IEcho.class)
	class Client implements IEcho{
		public String sayHello(String name) {
			OkhttpHproseClient client3=new OkhttpHproseClient();
			ServiceRegistry zooKeeperServiceRegistry = new ZooKeeperServiceRegistry("127.0.0.1:2181");
			String discover1 = zooKeeperServiceRegistry.discover(IEcho1.class.getName());
			IEcho1 iecho = client3.useService(discover1,IEcho1.class);
			return iecho.sayHello1("中国");
		}

	}

	interface IEcho {
		public String sayHello(String name);
	}
}