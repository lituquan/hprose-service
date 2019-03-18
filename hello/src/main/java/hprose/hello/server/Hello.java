package hprose.hello.server;

import java.util.concurrent.TimeUnit;

import org.springframework.context.annotation.Bean;

import brave.Tracing;
import brave.context.log4j2.ThreadContextCurrentTraceContext;
import brave.http.HttpTracing;
import hprose.config.annotation.RpcService;
import zipkin2.codec.SpanBytesEncoder;
import zipkin2.reporter.AsyncReporter;
import zipkin2.reporter.Sender;
import zipkin2.reporter.okhttp3.OkHttpSender;

@RpcService(IEcho.class)
public class Hello implements IEcho{
    public String sayHello(String name) {
        return "Hello " + name + "!";
    }

	@Bean(name="httpTracing")
	public HttpTracing getHttptracing(){
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
