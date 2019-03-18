package org.springframework.remoting.hprose;

import hprose.register.ServiceDiscovery;
import org.springframework.beans.BeansException;
import org.springframework.remoting.hprose.HproseProxyFactoryBean;
import org.springframework.remoting.hprose.SpringContextUtil;
import org.springframework.stereotype.Component;

import java.util.HashMap;
import java.util.Map;

@Component
public class HproseClientFactory{
	static Map<String,HproseProxyFactoryBean> map=new HashMap<>();

	public static <T> T create(Class<T> serviceInterface) throws BeansException{
		if (serviceInterface != null && !serviceInterface.isInterface()) {
			throw new IllegalArgumentException("'serviceInterface' must be an interface");
		}
		HproseProxyFactoryBean factory=null;
		ServiceDiscovery discovery= SpringContextUtil.getBean(ServiceDiscovery.class);//注册中心
		if(!map.containsKey(serviceInterface.getName())){
			factory=new HproseProxyFactoryBean();
			factory.setDiscovery(discovery);
			factory.setServiceInterface(serviceInterface);
			factory.afterPropertiesSet();
			map.put(serviceInterface.getName(),factory);
		}else{
			factory=map.get(serviceInterface.getName());
		}
		return (T)factory.create();
	}
}
