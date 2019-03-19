package hprose.hello.server;

import hprose.config.annotation.RpcService;
import org.springframework.beans.factory.annotation.Autowired;

@RpcService(IEcho1.class)
public class Client implements IEcho1{
	@Autowired
	IEcho iecho;

	public String sayHello1(String name) {
		return iecho.sayHello("中国");
	}
}
