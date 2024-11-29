protoc-compile:
	protoc -I . -I proto/googleapis \
		--go_out=. \
		--go-grpc_out=. \
		--grpc-gateway_out=. \
		proto/service.proto

	