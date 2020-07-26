# grpc-gateway支持

分别表示两类功能的支持：
* tag.pb.go
* tag.pb.gw.go

annotations.proto文件中的核心使用的http.proto文件中的一部分，总体来说，主要是对http转换提供支持，
定义了proto文件中可用于定义API服务的HTTP的相关配置，并且可以指定把每个rpc方法都映射到一个或多个HTTP rest方法上