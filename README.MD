This is a demo repository for learning rpc.  
Course Link: https://study.163.com/course/courseMain.htm?courseId=1212937804
## Demo
+ http-api-rpc 
    simple remote procedure call using http
+ go-rpc (implementing rpc with the standard library in golang)
    + /http  
        implement by http protocol.
    + /tcp  
        implement by tcp protocol.
    + /gob-json  
        rpc in golang uses gob parsing by default, but can also be parsed by json.
+ grpc
    + /hello_world  
        first demo.
    + /add  
        addServer & addClient.
+ probuf-demo1  
    learn protobuf syntax.
+ probuf-demo2-oneof  
    learn protobuf oneof and optional.
+ probuf-demo3-wrap&fieldmask  
    learn how to use warpValue and fieldMask.

## do follow these steps
```
# get gRPC 
go install google.golang.org/grpc@latest

# install Protocol Buffers v3
download from https://github.com/google/protobuf/releases and add bin to your env.

# install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# install protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# install to handle grpc status-err
go install google.golang.org/genproto/googleapis/rpc/errdetails@latest

# install protoc-gen-grpc-gateway v2
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2

```
