.PHONY: gen help

PROTO_DIR = proto

gen:
	protoc \
	--proto_path=proto \
	--go_out=proto \
	--go_opt=paths=source_relative  \
	--go-grpc_out=proto \
	--go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=proto \
	--grpc-gateway_opt=paths=source_relative \
	$(shell find $(PROTO_DIR) -iname "*.proto" )

help:
	@echo "make gen - 生成代码及 grpc 代码"