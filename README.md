# hprose-service

hprose 是一个跨语言rpc 框架。提供了发布服务、调用服务功能。

为实现服务治理，本项目引入服务注册、服务跟踪、服务熔断。

服务注册使用zookeeper、服务跟踪使用okhttp+zipkin、服务熔断还没有做。业务基于java、go ,java要做spring 整合。

服务注册使用zookeeper
  java参考https://gitee.com/huangyong/rpc.git
  注意：因为java 默认序列化不是直接转byte的,要设置为BytesPushThroughSerializer。
  
  go参考https://www.v2ex.com/t/440662
  
服务调用跟踪使用zipkin
  java参考https://gitee.com/mozhu/zipkin-learning.git,由于参考的接口主要是基于OKhttp调用，所以不使用hprose-java的HproseHttpClient。
  
  go 参考 https://github.com/openzipkin-contrib/zipkin-go-opentracing/tree/master/examples

hprose 的拦截器默认传递的是报文,比较难获取http.request  http.response对象。所以采用比较极端的方式，直接重写服务响应部分和客户端的请求部分。
