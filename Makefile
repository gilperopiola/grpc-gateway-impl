# Makefile :)

# Variables
# Targets
protoc-gen: 
	protoc -I=./protos --go_out=./out --go-grpc_out=./out --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative protos/users.proto
	protoc -I=./protos --grpc-gateway_out=./out --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative protos/users.proto
