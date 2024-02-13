# gRPC Gateway Implementation ;)

# Variables
# ... (not needed for now)

# Targets
.PHONY: all run gen protoc-gen swagger-gen

all: gen run

run:
	go run cmd/main.go

gen:
	make protoc-gen
	make swagger-gen

protoc-gen: 
	mkdir -p ./out

	protoc -I=./protos --go_out=./out --go-grpc_out=./out --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative protos/users.proto
	protoc -I=./protos --grpc-gateway_out=./out --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative protos/users.proto

	mv "./out/users.pb.go"      "./pkg/users"
	mv "./out/users_grpc.pb.go" "./pkg/users"
	mv "./out/users.pb.gw.go"   "./pkg/users"

	rm -rf ./out

swagger-gen:
	mkdir -p ./out

	protoc -I=./protos --openapiv2_out=./out protos/users.proto
	mv "./out/users.swagger.json" "./docs"

	rm -rf ./out

