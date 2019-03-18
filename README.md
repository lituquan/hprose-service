# hprose-service

hprose 是一个跨语言rpc 框架。提供了发布服务、调用服务功能。

为实现服务治理，本项目引入服务注册、服务跟踪、服务熔断。

服务注册使用zookeeper、服务跟踪使用okhttp+zipkin、服务熔断还没有做。业务基于java、go ,java要做spring 整合。

服务注册使用zookeeper
  参考https://gitee.com/huangyong/rpc.git
  
服务调用跟踪使用zipkin
  参考https://gitee.com/mozhu/zipkin-learning.git
  由于参考的接口主要是基于OKhttp调用，所以不使用hprose-java的HproseHttpClient。
  （HproseHttpClient + Filter传递上下文还没调通）
