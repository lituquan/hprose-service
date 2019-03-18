package org.springframework.remoting.hprose;

import brave.http.HttpTracing;
import brave.okhttp3.TracingInterceptor;
import hprose.client.ClientContext;
import hprose.client.CookieManager;
import hprose.client.HproseClient;
import hprose.common.HproseException;
import hprose.io.HproseMode;
import hprose.util.StrUtil;
import hprose.util.concurrent.Promise;
import hprose.util.concurrent.Threads;
import okhttp3.*;

import javax.net.ssl.HostnameVerifier;
import javax.net.ssl.SSLSocketFactory;
import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.ByteBuffer;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class OkhttpHproseClient extends HproseClient {

    OkHttpClient client;
    HttpTracing httpTracing;

    private static volatile ExecutorService pool = Executors.newCachedThreadPool();
    static {
        Threads.registerShutdownHandler(new Runnable() {
            public void run() {
                ExecutorService p = pool;
                pool = Executors.newCachedThreadPool();
                p.shutdownNow();
            }
        });
    }
    private final ConcurrentHashMap<String, String> headers = new ConcurrentHashMap<String, String>();
    private static boolean disableGlobalCookie = false;
    private static CookieManager globalCookieManager = new CookieManager();
    private final CookieManager cookieManager = disableGlobalCookie ? new CookieManager() : globalCookieManager;
    private boolean keepAlive = true;
    private int keepAliveTimeout = 300;
    private String proxyHost = null;
    private int proxyPort = 80;
    private String proxyUser = null;
    private String proxyPass = null;
    private HostnameVerifier hv = null;
    private SSLSocketFactory sslsf = null;

    public static void setThreadPool(ExecutorService threadPool) {
        pool = threadPool;
    }

    public static void setDisableGlobalCookie(boolean value) {
        disableGlobalCookie = value;
    }

    public static boolean isDisableGlobalCookie() {
        return disableGlobalCookie;
    }

    public OkhttpHproseClient() {
        super();
    }

    public OkhttpHproseClient(String uri) {
        super(uri);
    }

    public OkhttpHproseClient(HproseMode mode) {
        super(mode);
    }

    public OkhttpHproseClient(String uri, HproseMode mode) {
        super(uri, mode);
    }

    public OkhttpHproseClient(String[] uris) {
        super(uris);
    }

    public OkhttpHproseClient(String[] uris, HproseMode mode) {
        super(uris, mode);
    }

    public static HproseClient create(String uri, HproseMode mode) throws IOException, URISyntaxException {
        String scheme = (new URI(uri)).getScheme();
        if (!"http".equalsIgnoreCase(scheme) && !"https".equalsIgnoreCase(scheme)) {
            throw new HproseException("This client doesn't support " + scheme + " scheme.");
        }
        return new OkhttpHproseClient(uri, mode);
    }

    public static HproseClient create(String[] uris, HproseMode mode) throws IOException, URISyntaxException {
        for (int i = 0, n = uris.length; i < n; ++i) {
            String scheme = (new URI(uris[i])).getScheme();
            if (!"http".equalsIgnoreCase(scheme) && !"https".equalsIgnoreCase(scheme)) {
                throw new HproseException("This client doesn't support " + scheme + " scheme.");
            }
        }
        return new OkhttpHproseClient(uris, mode);
    }

    public void setHeader(String name, String value) {
        String nl = name.toLowerCase();
        if (!nl.equals("content-type") &&
                !nl.equals("content-length") &&
                !nl.equals("connection") &&
                !nl.equals("keep-alive") &&
                !nl.equals("host")) {
            if (value == null) {
                headers.remove(name);
            }
            else {
                headers.put(name, value);
            }
        }
    }

    public String getHeader(String name) {
        return headers.get(name);
    }

    public Map<String,String> getHeaders() {
        return headers;
    }

    public boolean isKeepAlive() {
        return keepAlive;
    }

    public void setKeepAlive(boolean keepAlive) {
        this.keepAlive = keepAlive;
    }

    public int getKeepAliveTimeout() {
        return keepAliveTimeout;
    }

    public void setKeepAliveTimeout(int keepAliveTimeout) {
        this.keepAliveTimeout = keepAliveTimeout;
    }

    public String getProxyHost() {
        return proxyHost;
    }

    public void setProxyHost(String proxyHost) {
        this.proxyHost = proxyHost;
    }

    public int getProxyPort() {
        return proxyPort;
    }

    public void setProxyPort(int proxyPort) {
        this.proxyPort = proxyPort;
    }

    public String getProxyUser() {
        return proxyUser;
    }

    public void setProxyUser(String proxyUser) {
        this.proxyUser = proxyUser;
    }

    public String getProxyPass() {
        return proxyPass;
    }

    public void setProxyPass(String proxyPass) {
        this.proxyPass = proxyPass;
    }

    public HostnameVerifier getHostnameVerifier() {
        return hv;
    }

    public void setHostnameVerifier(HostnameVerifier hv) {
        this.hv = hv;
    }

    public SSLSocketFactory getSSLSocketFactory() {
        return sslsf;
    }

    public void setSSLSocketFactory(SSLSocketFactory sslsf) {
        this.sslsf = sslsf;
    }

    @SuppressWarnings({"unchecked"})
    private ByteBuffer syncSendAndReceive(ByteBuffer request, ClientContext context) throws Throwable {
        //跟踪
        httpTracing = SpringContextUtil.getBean(HttpTracing.class);
        client = new OkHttpClient.Builder()
                .dispatcher(new Dispatcher(
                        httpTracing.tracing().currentTraceContext().executorService(
                                new Dispatcher().executorService())
                ))
                .addNetworkInterceptor(TracingInterceptor.create(httpTracing))
                .build();
        MediaType mediaType = MediaType.parse("application/hprose");
        RequestBody requestBody = RequestBody.create(mediaType, StrUtil.toString(request));
        Request request1 = new Request.Builder()
                .url(uri)
                .post(requestBody)
                .build();

        try {
            Response response = client.newCall(request1).execute();
            return  ByteBuffer.wrap(response.body().bytes());
        }
        catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }

    @Override
    protected Promise<ByteBuffer> sendAndReceive(final ByteBuffer request, final ClientContext context) {
        final Promise<ByteBuffer> promise = new Promise<ByteBuffer>();
        pool.submit(new Runnable() {
            public void run() {
                try {
                    promise.resolve(syncSendAndReceive(request, context));
                }
                catch (Throwable ex) {
                    promise.reject(ex);
                }
            }
        });
        return promise;
    }
}
