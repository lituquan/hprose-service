package hprose.hello.server1;

import hprose.hello.server.IEcho1;
import org.springframework.remoting.hprose.OkhttpHproseClient;
import hprose.register.ServiceDiscovery;
import hprose.register.zookeeper.ZooKeeperServiceDiscovery;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.TestExecutionListeners;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;
import org.springframework.test.context.support.DependencyInjectionTestExecutionListener;
import org.springframework.test.context.support.DirtiesContextTestExecutionListener;

import java.util.logging.Logger;


/**
 * Unit test for simple App.
 */

@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(locations={"classpath:hprose-consumer.xml"})
@TestExecutionListeners(
		{ 
			DependencyInjectionTestExecutionListener.class,
			DirtiesContextTestExecutionListener.class 
		})
public class Hprose{
	Logger logger = Logger.getLogger(Hprose.class.getName());
	Logger log=logger;

	@Test
	public void get(){
		String discover="http://localhost:8090/hello.server/hprose";
		OkhttpHproseClient client3=new OkhttpHproseClient();
		ServiceDiscovery zooKeeperServiceDiscovery = new ZooKeeperServiceDiscovery("192.168.6.31:31089");
		String discover1 = zooKeeperServiceDiscovery.discover(IEcho1.class.getName());
		IEcho1 h = client3.useService(discover1,IEcho1.class);
		System.out.println(h.sayHello1("11111"));
	}

	//工程名
	@Value("${rpc.path}")
	private String path = null;
	//工程名
	@Value("${rpc.trace.path}")
	private String trace_path = null;
}
