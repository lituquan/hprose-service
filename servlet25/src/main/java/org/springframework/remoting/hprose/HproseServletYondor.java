package org.springframework.remoting.hprose;

import hprose.common.HproseMethods;
import hprose.server.HproseServlet;

import java.util.Iterator;
import java.util.Map;

/**
 * 使用java发布对象
 * @author tuquan
 *
 */
public class HproseServletYondor extends HproseServlet{
	//存放 服务名 与 服务对象 之间的映射关系
	private static Map<String, Object> handlerMap = null;
	@Override
	public void setGlobalMethods(HproseMethods methods) {
		super.setGlobalMethods(methods);
		//注册HelloService下所有的public方法
		handlerMap=SpringContextUtil.getServiceMap();
		Iterator<String> iterator = handlerMap.keySet().iterator();
		try{        	
			while(iterator.hasNext()){
				Class<?> clazz = Class.forName( iterator.next());
				Object object=SpringContextUtil.getBean(clazz);
				methods.addInstanceMethods(object,clazz);
			}
		}catch(Exception e){
			e.printStackTrace();
		}
	}
//
//	@Override
//	public void init(ServletConfig config) throws ServletException {
//		super.init(config);
//		String contextPath = config.getServletContext().getContextPath();
//
//		Sender sender = OkHttpSender.create(SpringContextUtil.getTracePath());
//		AsyncReporter asyncReporter = AsyncReporter.builder(sender).closeTimeout(500, TimeUnit.MILLISECONDS)
//				.build(SpanBytesEncoder.JSON_V2);
//
//		Tracing tracing = Tracing.newBuilder()
//				.localServiceName(contextPath.substring(1))
//				.spanReporter(asyncReporter)
//				.currentTraceContext(ThreadContextCurrentTraceContext.create())
//				.build();
//
//		HttpTracing httpTracing = HttpTracing.create(tracing);
//		config.getServletContext().setAttribute("httpTracing", httpTracing);
//
//		BraveServerInterceptor braveServerInterceptor = new BraveServerInterceptor(httpTracing);
//		service.setFilter(new StatFilter());
//		service.setFilter(braveServerInterceptor);
//	}
}