# gRPC Gateway Implementation ;)

# Variables
# ... (not needed for now)

# Targets
all: gen run

run:
	go run cmd/main.go

generate:
	make generate-protos
	make generate-swagger

prepare:
	mkdir -p ./out

clean:
	rm -rf ./out

generate-protos: 
	make prepare

	protoc -I=./protos --go_out=./out --go-grpc_out=./out --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative protos/users.proto
	protoc -I=./protos --grpc-gateway_out=./out --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative protos/users.proto

	make move-protos
	make clean

generate-swagger:
	make prepare

	protoc -I=./protos --openapiv2_out=./out protos/users.proto

	make move-swagger
	make clean

move-protos:
	mv "./pkg/users/users.pb.go"      "./pkg/users"
	mv "./pkg/users/users_grpc.pb.go" "./pkg/users"
	mv "./pkg/users/users.pb.gw.go"   "./pkg/users"

move-swagger:
	mv "./out/users.swagger.json" "./docs"
