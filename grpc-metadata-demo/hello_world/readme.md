# 第一个 gRPC 实例
hello world!

## 三步
+ 编写 protobuf 文件
+ 生成代码(Server And Client)
```
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
pb/hello.proto
```
+ 完成业务逻辑