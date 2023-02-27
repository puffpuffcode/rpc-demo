protoc --proto_path=proto \
--go_out=proto --go_opt=paths=source_relative  \
--go-grpc_out=proto --go-grpc_opt=paths=source_relative \
book/price.proto book/book.proto author/author.proto