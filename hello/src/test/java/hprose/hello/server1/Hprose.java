package hprose.hello.server1;

import hprose.hello.server.IEcho;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.TestExecutionListeners;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;
import org.springframework.test.context.support.DependencyInjectionTestExecutionListener;
import org.springframework.test.context.support.DirtiesContextTestExecutionListener;

import javax.annotation.Resource;
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

	@Resource()
	IEcho iecho;

	@Test
	public void get(){
		System.out.println(iecho.sayHello("111"));
	}
	//工程名
	@Value("${rpc.path}")
	private String path = null;
	//工程名
	@Value("${rpc.trace.path}")
	private String trace_path = null;
	

}
