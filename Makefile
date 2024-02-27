# gRPC Gateway Implementation ;)
#
# This Makefile is used to generate the gRPC and gRPC Gateway files from the .proto files.
# It also generates the Swagger file for the gRPC Gateway.
# The Makefile also has a target to run the server.
# The server not only serves gRPC requests but also serves HTTP requests through the 
# gRPC Gateway.

# General Variables
# 
# PROTOS_DIR: The directory where the .proto files are located.
# PBS_OUT_BASE_DIR: The directory where the .pb files will be generated.
# SWAGGER_OUT_DIR: The directory where the Swagger file will be generated.
PROTOS_DIR := ./protos
PBS_OUT_BASE_DIR := ./pkg
SWAGGER_OUT_DIR := ./docs

# User Variables
#
# USERS: The name of the entity / service / package.
# USERS_PROTO_FILE: The corresponding .proto file, defining the service and its methods
# alongside the Request and Response types.
# USERS_PKG_DIR: The directory where the generated files will be moved to,
# inside of the actual source code. These will be the structs that you'll end up using.
USERS := users
USERS_PROTO_FILE := $(PROTOS_DIR)/$(USERS).proto
USERS_PBS_OUT_DIR := $(PBS_OUT_BASE_DIR)/$(USERS)

# Default target.
# 
# The default target is the "all" target. 
# It generates the gRPC and gRPC Gateway files, as well as the Swagger file.
# Finally, it runs the server.
all: generate run

# Generate the gRPC and gRPC Gateway files, as well as the Swagger file.
generate: generate-pbs generate-swagger

# Run both the gRPC server and the gRPC Gateway server.
# The gRPC server usually listens on port :50051 and the gRPC Gateway on port :8080.
run:
	go run cmd/main.go

# Generate the .pb and .pb.gw files for gRPC and gRPC Gateway.
generate-pbs:
	protoc -I=$(PROTOS_DIR) --go_out=$(USERS_PBS_OUT_DIR) --go-grpc_out=$(USERS_PBS_OUT_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(USERS_PROTO_FILE)
	protoc -I=$(PROTOS_DIR) --grpc-gateway_out=$(USERS_PBS_OUT_DIR) --grpc-gateway_opt=paths=source_relative $(USERS_PROTO_FILE)

# Generate the Swagger file for the gRPC Gateway.
generate-swagger:
	protoc -I=$(PROTOS_DIR) --openapiv2_out=$(SWAGGER_OUT_DIR) $(USERS_PROTO_FILE)

# Test the app.
test:
	go test ./... -cover