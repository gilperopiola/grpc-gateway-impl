# gRPC Gateway Implementation ;)

# General Variables
PROTOS_DIR := ./protos
PBS_OUT_BASE_DIR := ./pkg
DOCS_OUT_DIR := ./docs

# Users Variables
USERS := users
USERS_PROTO_FILE := $(PROTOS_DIR)/$(USERS).proto
USERS_PBS_OUT_DIR := $(PBS_OUT_BASE_DIR)/$(USERS)

# Default target.
# 
# Generates the gRPC and gRPC Gateway code from the .proto files.
# Also generates the Swagger documentation.
# Also runs the gRPC Server and the gRPC Gateway.
all: clean generate test run

# Generates the gRPC and gRPC Gateway files, as well as the Swagger documentation.
generate: generate-pbs generate-swagger

# Runs both the gRPC server and the gRPC Gateway.
# The gRPC Server usually listens on port :50051 and the gRPC Gateway on port :8080.
run:
	go mod tidy
	go run cmd/main.go

# Generates the gRPC and gRPC Gateway code from the .proto files.
generate-pbs:
	protoc -I=$(PROTOS_DIR) --go_out=$(USERS_PBS_OUT_DIR) --go-grpc_out=$(USERS_PBS_OUT_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(USERS_PROTO_FILE)
	protoc -I=$(PROTOS_DIR) --grpc-gateway_out=$(USERS_PBS_OUT_DIR) --grpc-gateway_opt=paths=source_relative $(USERS_PROTO_FILE)

# Generates the Swagger documentation.
generate-swagger:
	protoc -I=$(PROTOS_DIR) --openapiv2_out=$(DOCS_OUT_DIR) $(USERS_PROTO_FILE)

# Clean cache.
clean:
	go clean -cache -modcache -testcache

# Test the app.
test:
	go test ./... -cover