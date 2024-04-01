.PHONY: proto

proto:
	protoc \
		--go-grpc_out="${PWD}"/api \
		--go_out="${PWD}"/api \
		-I "${PWD}/api/hlcv1" \
		hlcv1.proto
