package hprose.hello.server;

import hprose.config.annotation.RpcService;

@RpcService(IEcho.class)
public class Hello implements IEcho{
    public String sayHello(String name) {
        return "Hello " + name + "!";
    }
}
