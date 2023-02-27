This is a demo repository for learning rpc.
## Demo
+ http-api-rpc (simple remote procedure call using http)
+ go-rpc (implementing rpc with the standard library in golang)
    + /http implement by http protocol.
    + /tcp implement by tcp protocol.
    + /gob-json rpc in golang uses gob parsing by default, but can also be parsed by json.

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
```