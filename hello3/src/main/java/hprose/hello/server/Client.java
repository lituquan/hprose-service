package hprose.hello.server;

import hprose.config.annotation.RpcService;

import javax.annotation.Resource;

@RpcService(IEcho2.class)
public class Client implements IEcho2{

	@Resource
	IEcho1 iecho;

	public String sayHello(String name) {
		return iecho.sayHello1("中国");
	}
}
