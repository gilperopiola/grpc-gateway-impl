# gRPC Gateway Implementation ;)

# Variables
PKG_DIR := ./pkg/users
OUT_DIR := ./out
DOCS_DIR := ./docs
PROTO_DIR := ./protos
USERS_PROTO_FILE := $(PROTO_DIR)/users.proto

# Targets
# External targets
all: generate run

generate: generate-protos generate-swagger

run:
	go run cmd/main.go

# Internal targets

generate-protos: prepare
	protoc -I=$(PROTO_DIR) --go_out=$(PROTO_DIR) --go-grpc_out=$(PROTO_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(USERS_PROTO_FILE)
	protoc -I=$(PROTO_DIR) --grpc-gateway_out=$(PROTO_DIR) --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative $(USERS_PROTO_FILE)

	make move-protos --no-print-directory
	make clean --no-print-directory

generate-swagger: prepare
	protoc -I=$(PROTO_DIR) --openapiv2_out=$(OUT_DIR) $(USERS_PROTO_FILE)

	make move-swagger --no-print-directory
	make clean --no-print-directory

prepare:
	mkdir -p $(OUT_DIR)

clean:
	rm -rf $(OUT_DIR)

move-protos:
	mv "$(OUT_DIR)/users.pb.go"      "$(PKG_DIR)"
	mv "$(OUT_DIR)/users_grpc.pb.go" "$(PKG_DIR)"
	mv "$(OUT_DIR)/users.pb.gw.go"   "$(PKG_DIR)"

move-swagger:
	mv "$(OUT_DIR)/users.swagger.json" "$(DOCS_DIR)"
