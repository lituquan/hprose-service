package org.springframework.remoting.hprose;

import brave.Tracing;
import brave.context.log4j2.ThreadContextCurrentTraceContext;
import brave.http.HttpTracing;
import hprose.config.annotation.RpcService;
import hprose.register.ServiceDiscovery;
import hprose.register.ServiceRegistry;
import org.apache.dubbo.common.utils.NetUtils;
import org.springframework.beans.BeansException;
import org.springframework.beans.factory.NoSuchBeanDefinitionException;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;
import zipkin2.codec.SpanBytesEncoder;
import zipkin2.reporter.AsyncReporter;
import zipkin2.reporter.Sender;
import zipkin2.reporter.okhttp3.OkHttpSender;

import javax.management.MBeanServer;
import javax.management.MalformedObjectNameException;
import javax.management.ObjectName;
import javax.management.Query;
import java.lang.annotation.Annotation;
import java.lang.management.ManagementFactory;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.TimeUnit;

/**
 * spring上下文配置
 * @author Mingchenchen
 *
 */
@Component
public class SpringContextUtil implements ApplicationContextAware {
	//注册中心
	public ServiceRegistry registry;
    private static ApplicationContext applicationContext;

    @Value("${rpc.path}")
    private String path;

    @Value("${rpc.trace.path}")
    private String trace_path;

    public static String trace_path_static;
    //存放 服务名 与 服务对象 之间的映射关系
    private static Map<String, Object> handlerMap = new HashMap<>();

    @Override
    public void setApplicationContext(ApplicationContext applicationContext1) throws BeansException
    {
        applicationContext = applicationContext1;
        try {
        	setService();
			register();
			trace_path_static=trace_path;
		} catch (Exception e) {
			e.printStackTrace();
		}
    }

    public void setService(){
    	// 扫描带有 RpcService 注解的类并初始化 handlerMap 对象
        Map<String, Object> serviceBeanMap = getBeansWithAnnotation(RpcService.class);
        if (serviceBeanMap!=null && serviceBeanMap.size()>0) {
            for (Object serviceBean : serviceBeanMap.values()) {
            	//getSuperclass(),使用spring加载的对象被代理了，要用父类获取
                System.out.println("serviceBean.getClass()"+serviceBean.getClass().getName());
                RpcService rpcService =serviceBean.getClass().getSuperclass().getAnnotation(RpcService.class);
                if(rpcService==null){
                    rpcService =serviceBean.getClass().getAnnotation(RpcService.class);
                }
                String serviceName = rpcService.value().getName();
                handlerMap.put(serviceName, serviceBean);
            }
        }
    }

    public static Map<String, Object> getServiceMap() {
        return handlerMap;
    }
 
    public static ApplicationContext getApplicationContext() {
        return applicationContext;
    }
 
    /**
     * 注意 bean name默认 = 类名(首字母小写)
     * 例如: A8sClusterDao = getBean("k8sClusterDao")
     * @param name
     * @return
     * @throws BeansException
     */
    public static Object getBean(String name) throws BeansException {
        return applicationContext.getBean(name);
    }
    
    /**
     * 注意 bean name默认 = 类名(首字母小写)
     * 例如: A8sClusterDao = getBean("k8sClusterDao")
     * @param nameType
     * @return
     * @throws BeansException
     */
    public static<T> T getBean(Class<T> nameType) throws BeansException {
        return applicationContext.getBean(nameType);
    }
 
    public static boolean containsBean(String name) {
        return applicationContext.containsBean(name);
    }
 
    public static boolean isSingleton(String name) throws NoSuchBeanDefinitionException {
        return applicationContext.isSingleton(name);
    }
    
    public static Map<String, Object> getBeansWithAnnotation(Class<? extends Annotation> annotationType)
			throws BeansException{
    	return applicationContext.getBeansWithAnnotation(annotationType);		    	
	}

    public String getPath() {
        return path;
    }
    
    public static String getTracePath() {
        return trace_path_static;
    }

    @Bean
    public  HttpTracing getHttptracing(){
        Sender sender = OkHttpSender.create(trace_path);
        AsyncReporter asyncReporter = AsyncReporter.builder(sender).closeTimeout(500, TimeUnit.MILLISECONDS)
                .build(SpanBytesEncoder.JSON_V2);
        Tracing tracing = Tracing.newBuilder()
                .localServiceName(path)
                .spanReporter(asyncReporter)
                .currentTraceContext(ThreadContextCurrentTraceContext.create())
                .build();
        HttpTracing httpTracing = HttpTracing.create(tracing);
        return httpTracing;
    }
    
    public String getServiceUrl(){
        try{
            String domain=NetUtils.getLocalHost()+":"+getNonSecurePort();
            String address="http://"+domain+"/"+getPath()+"/hprose";
            return address;
        }catch (Exception e){
            e.printStackTrace();
        }
        return "";
    }

    /**
     * path 和 地址的组成要独立出来
     * @throws Exception
     */
	public void register() throws Exception {
		try {
			handlerMap=getServiceMap();
			if(handlerMap.size()>0){
                registry=(ServiceRegistry) getBean(ServiceRegistry.class);
                String address=getServiceUrl();
                Iterator<String> iterator = handlerMap.keySet().iterator();
                while(iterator.hasNext()){
                    String serviceName=iterator.next();
                    registry.register(serviceName,address);
                }
            }
		} catch (Exception e) {
			e.printStackTrace();
		}
	}

	private int getNonSecurePort() throws MalformedObjectNameException{
        int tomcatPort=8080;
        try{
            MBeanServer beanServer = ManagementFactory.getPlatformMBeanServer();
            Set<ObjectName> objectNames = beanServer.queryNames(new ObjectName("*:type=Connector,*"),
                    Query.match(Query.attr("protocol"), Query.value("HTTP/1.1")));
            tomcatPort = Integer.valueOf(objectNames.iterator().next().getKeyProperty("port"));
        }catch (Exception e){
            e.printStackTrace();
        }
        return tomcatPort;
    }

    public  Object getReference(Class clazz) throws Exception{
        HproseProxyFactoryBean factory= (HproseProxyFactoryBean) getBean(HproseProxyFactoryBean.class);
        if(factory==null){
            factory=new HproseProxyFactoryBean();
            ServiceDiscovery discovery= (ServiceDiscovery) getBean(ServiceDiscovery.class);
            if(discovery==null){
                throw new Exception("no discovery");
            }
            factory.setDiscovery(discovery);
            factory.setServiceInterface(clazz);
        }
        return  factory.create();
    }
}
