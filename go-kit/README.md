# go-kit
1. transport层主要负责网络传输，处理HTTP、grpc、thrift等相关逻辑
2. endpoint层主要负责request/response格式的转换，以及公用拦截器相关的逻辑。是go-kit的核心，采用洋葱模式，提供了对日志、限流、熔断、链路追踪和服务监控等方面的扩展能力
3. service层主要专注于业务逻辑

