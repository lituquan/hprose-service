package hprose.hello.server;

import brave.Tracing;
import brave.context.log4j2.ThreadContextCurrentTraceContext;
import brave.http.HttpTracing;
import org.springframework.remoting.hprose.OkhttpHproseClient;
import hprose.register.ServiceDiscovery;
import hprose.register.zookeeper.ZooKeeperServiceDiscovery;
import org.springframework.beans.factory.annotation.Autowired;
import zipkin2.codec.SpanBytesEncoder;
import zipkin2.reporter.AsyncReporter;
import zipkin2.reporter.Sender;
import zipkin2.reporter.okhttp3.OkHttpSender;

import javax.annotation.Resource;
import java.util.concurrent.TimeUnit;

//@RpcService(IEcho1.class)
public class Client implements IEcho1{

	public static void main(String[] args) throws Throwable{
		String discover="http://localhost:8090/hello.server/hprose";
		OkhttpHproseClient client3=new OkhttpHproseClient();
		ServiceDiscovery zooKeeperServiceDiscovery = new ZooKeeperServiceDiscovery("192.168.6.31:31089");
		String discover1 = zooKeeperServiceDiscovery.discover(IEcho1.class.getName());
		IEcho1 h = client3.useService(discover1,IEcho1.class);
		System.out.println(h.sayHello1("11111"));
	}
	
	@Autowired
	HttpTracing httpTracing;

	@Resource(name="iecho")
	IEcho iecho;

	public String sayHello1(String name) {	
		return iecho.sayHello("中国");
	}
	
	public static HttpTracing getHttptracing(){
		Sender sender = OkHttpSender.create("http://192.168.6.30:30550/api/v2/spans");
		AsyncReporter asyncReporter = AsyncReporter.builder(sender).closeTimeout(500, TimeUnit.MILLISECONDS)
				.build(SpanBytesEncoder.JSON_V2);

		Tracing tracing = Tracing.newBuilder()
				.localServiceName("client1234")
				.spanReporter(asyncReporter)
				.currentTraceContext(ThreadContextCurrentTraceContext.create())
				.build();

		return HttpTracing.create(tracing);
	}
}
