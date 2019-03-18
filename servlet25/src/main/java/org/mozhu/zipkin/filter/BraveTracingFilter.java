package org.mozhu.zipkin.filter;

import brave.Tracing;
import brave.context.log4j2.ThreadContextCurrentTraceContext;
import brave.http.HttpTracing;
import brave.servlet.TracingFilter;
import zipkin2.codec.SpanBytesEncoder;
import zipkin2.reporter.AsyncReporter;
import zipkin2.reporter.Sender;
import zipkin2.reporter.okhttp3.OkHttpSender;

import javax.servlet.*;
import java.io.IOException;
import java.util.concurrent.TimeUnit;

/**
 * 参考https://gitee.com/mozhu/zipkin-learning.git
 */
public class BraveTracingFilter implements Filter {
    Filter tracingFilter;

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
        Sender sender = OkHttpSender.create("http://192.168.6.30:30550/api/v2/spans");
        AsyncReporter asyncReporter = AsyncReporter.builder(sender).closeTimeout(500, TimeUnit.MILLISECONDS)
                .build(SpanBytesEncoder.JSON_V2);

        Tracing tracing = Tracing.newBuilder()
                .localServiceName(System.getProperty("zipkin.service", filterConfig.getServletContext().getContextPath().substring(1)))
                .spanReporter(asyncReporter)
                .currentTraceContext(ThreadContextCurrentTraceContext.create())
                .build();

        HttpTracing httpTracing = HttpTracing.create(tracing);
        filterConfig.getServletContext().setAttribute("HttpTracing", httpTracing);
        tracingFilter = TracingFilter.create(httpTracing);
        tracingFilter.init(filterConfig);
    }

    @Override
    public void doFilter(ServletRequest servletRequest, ServletResponse servletResponse, FilterChain filterChain) throws IOException, ServletException {
        tracingFilter.doFilter(servletRequest, servletResponse, filterChain);
    }

    @Override
    public void destroy() {
        tracingFilter.destroy();
    }

}
