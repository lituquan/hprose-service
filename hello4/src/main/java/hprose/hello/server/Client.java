package hprose.hello.server;

import hprose.config.annotation.RpcService;

import javax.annotation.Resource;

@RpcService(IEcho3.class)
public class Client implements IEcho3{

	@Resource(name="iecho")
	IEcho2 iecho;

	public String sayHello1(String name) {
		return iecho.sayHello("中国");
	}
}
