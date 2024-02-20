# gRPC Gateway Implementation ;)

# Variables
PROTO_DIR := ./protos
PKG_DIR := ./pkg/users
DOCS_DIR := ./docs
OUT_DIR := ./out

# Targets
# External targets
all: generate run

generate: generate-protos generate-swagger

run:
	go run cmd/main.go

# Internal targets

generate-protos: prepare
	protoc -I=./protos --go_out=./out --go-grpc_out=./out --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative protos/users.proto
	protoc -I=./protos --grpc-gateway_out=./out --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative protos/users.proto

	make move-protos --no-print-directory
	make clean --no-print-directory

generate-swagger: prepare
	protoc -I=./protos --openapiv2_out=./out protos/users.proto

	make move-swagger --no-print-directory
	make clean --no-print-directory

prepare:
	mkdir -p ./out

clean:
	rm -rf ./out

move-protos:
	mv "./out/users.pb.go"      "./pkg/users"
	mv "./out/users_grpc.pb.go" "./pkg/users"
	mv "./out/users.pb.gw.go"   "./pkg/users"

move-swagger:
	mv "./out/users.swagger.json" "./docs"
