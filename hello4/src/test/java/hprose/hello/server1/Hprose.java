package hprose.hello.server1;

import hprose.hello.server.IEcho2;
import hprose.hello.server.go.IEcho;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
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
public class Hprose {
	Logger logger = Logger.getLogger(Hprose.class.getName());
	Logger log=logger;

	@Autowired
	IEcho2 h;
	@Autowired
	IEcho hello;
	@Test
	public void get(){
		System.out.println(hello.sayHello111("11111"));
	}
}
