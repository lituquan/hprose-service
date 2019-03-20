/**********************************************************\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: http://www.hprose.com/                 |
|                   http://www.hprose.org/                 |
|                                                          |
\**********************************************************/
/**********************************************************\
 *                                                        *
 * HproseProxyFactoryBean.java                            *
 *                                                        *
 * HproseProxyFactoryBean for Java Spring Framework.      *
 *                                                        *
 * LastModified: Mar 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/
package org.springframework.remoting.hprose;

import org.springframework.remoting.support.UrlBasedRemoteAccessor;

import hprose.client.HproseClient;
import hprose.client.HproseHttpClient;
import hprose.client.HproseTcpClient;
import hprose.common.FilterHandler;
import hprose.common.HproseFilter;
import hprose.common.InvokeHandler;
import hprose.io.HproseMode;
import hprose.register.ServiceRegistry;

public class HproseProxyFactoryBean extends UrlBasedRemoteAccessor{

    private HproseClient client = null;
    private Exception exception = null;
    private boolean keepAlive = true;
    private int keepAliveTimeout = 300;
    private int timeout = 5000; 
    private String proxyHost = null;
    private int proxyPort = 80;
    private String proxyUser = null;
    private String proxyPass = null;
    private HproseMode mode = HproseMode.MemberMode;
    private HproseFilter filter = null;
    private InvokeHandler invokeHandler = null;
    private FilterHandler beforeFilterHandler = null;
    private FilterHandler afterFilterHandler = null;
    private ServiceRegistry discovery = null ;
    
    @Override
    public void afterPropertiesSet() {
        super.afterPropertiesSet();
        try {
            String serviceUrl = getServiceUrl();//serviceUrl 检测空、null
            if(serviceUrl.startsWith("http")){
                client =OkhttpHproseClient.create(serviceUrl,mode);
            }else{
                client = HproseClient.create(serviceUrl, mode);
            }
        }
        catch (Exception ex) {
            exception = ex;
        }
        if (client instanceof HproseHttpClient) {
            HproseHttpClient httpClient = (HproseHttpClient)client;
            httpClient.setKeepAlive(keepAlive);
            httpClient.setKeepAliveTimeout(keepAliveTimeout);
            httpClient.setTimeout(timeout);
            httpClient.setProxyHost(proxyHost);
            httpClient.setProxyPort(proxyPort);
            httpClient.setProxyUser(proxyUser);
            httpClient.setProxyPass(proxyPass);
            httpClient.use(invokeHandler);
            httpClient.beforeFilter.use(beforeFilterHandler);
            httpClient.afterFilter.use(afterFilterHandler);
        }
        if (client instanceof HproseTcpClient) {
            HproseTcpClient tcpClient = (HproseTcpClient)client;
            tcpClient.setTimeout(timeout);
        }
        client.setFilter(filter);
    }

// for HproseHttpClient
    public void setKeepAlive(boolean value) {
        keepAlive = value;
    }

    public void setKeepAliveTimeout(int value) {
        keepAliveTimeout = value;
    }

    public void setProxyHost(String value) {
        proxyHost = value;
    }

    public void setProxyPort(int value) {
        proxyPort = value;
    }

    public void setProxyUser(String value) {
        proxyUser = value;
    }

    public void setProxyPass(String value) {
        proxyPass = value;
    }

// for HproseClient
    public void setTimeout(int value) {
        timeout = value;
    }

    public void setMode(HproseMode value) {
        mode = value;
    }

    public void setFilter(HproseFilter filter) {
        this.filter = filter;
    }

    public void setInvokeHandler(InvokeHandler value) {
        invokeHandler = value;
    }

    public void setBeforeFilterHandler(FilterHandler value) {
        beforeFilterHandler = value;
    }

    public void setAfterFilterHandler(FilterHandler value) {
        afterFilterHandler = value;
    }

    public Object getObject() throws Exception {
        if (exception != null) {
            throw exception;
        }
        return client.useService(getServiceInterface());
    }

    public Object create() {
        if (exception != null) {
        	exception.printStackTrace();
        	return null;
        }
        return client.useService(getServiceInterface());
    }
    
    public ServiceRegistry getDiscovery() {
		return discovery;
	}

	public void setDiscovery(ServiceRegistry discovery) {
		this.discovery = discovery;
	}

	public HproseClient getClient() {
		return client;
	}

	public void setClient(HproseClient client) {
		this.client = client;
	}

	public Exception getException() {
		return exception;
	}

	public void setException(Exception exception) {
		this.exception = exception;
	}

	public boolean isKeepAlive() {
		return keepAlive;
	}

	public int getKeepAliveTimeout() {
		return keepAliveTimeout;
	}

	public int getTimeout() {
		return timeout;
	}

	public String getProxyHost() {
		return proxyHost;
	}

	public int getProxyPort() {
		return proxyPort;
	}

	public String getProxyUser() {
		return proxyUser;
	}

	public String getProxyPass() {
		return proxyPass;
	}

	public HproseMode getMode() {
		return mode;
	}

	public HproseFilter getFilter() {
		return filter;
	}

	public InvokeHandler getInvokeHandler() {
		return invokeHandler;
	}

	public FilterHandler getBeforeFilterHandler() {
		return beforeFilterHandler;
	}

	public FilterHandler getAfterFilterHandler() {
		return afterFilterHandler;
	}

	@Override
    public String getServiceUrl() {
        String serviceName=getServiceInterface().getName();
        super.setServiceUrl(discovery.discover(serviceName));
        return discovery.discover(serviceName);
    }
}