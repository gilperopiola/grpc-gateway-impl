###---------------------------------------------###
###                                             ###
###       - gRPC Gateway Implementation -       ###
###                                             ###
###--------------------------- by @gilperopiola ###

#------------------------------#
#          Variables           #
#------------------------------#

# Get the current git tag or commit hash.
VERSION := $(shell git describe --tags --always --dirty) 

# Set the path to the .proto's dir, the output dir for the generated code and the output dir for the Swagger documentation.
PROTOS_DIR := ./protos
PBS_OUT_BASE_DIR := ./pkg
DOCS_OUT_DIR := ./docs

# Users Service Variables
USERS := users
USERS_PROTO_FILE := $(PROTOS_DIR)/$(USERS).proto
USERS_PBS_OUT_DIR := $(PBS_OUT_BASE_DIR)/$(USERS)

# Don't print unnecessary output.
MAKEFLAGS += --no-print-directory

#------------------------------#
#         Main Targets         #
#------------------------------#

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

#------------------------------#

all: ## Cleans env, generates code and Swagger, runs tests, starts the application. Add fast=1 to skip cleaning and testing.
ifeq ($(fast),)
	@'$(MAKE)' clean generate test run
else 
	@'$(MAKE)' generate run
endif

#------------------------------#

run: ## Updates dependencies and starts the application.
	@echo ''
	go mod tidy
	go run cmd/main.go

#------------------------------#

version: ## Shows version.
	@echo $(VERSION)

#------------------------------#
#      Secondary Targets       #
#------------------------------#

generate: ## Generates gRPC and gRPC Gateway code + the Swagger documentation. 
	@'$(MAKE)' generate-pbs generate-swagger

generate-pbs: ## Generates gRPC and gRPC Gateway code.
	@echo ''
	protoc -I=$(PROTOS_DIR) --go_out=$(USERS_PBS_OUT_DIR) --go-grpc_out=$(USERS_PBS_OUT_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(USERS_PROTO_FILE)
	protoc -I=$(PROTOS_DIR) --grpc-gateway_out=$(USERS_PBS_OUT_DIR) --grpc-gateway_opt=paths=source_relative $(USERS_PROTO_FILE)
 
generate-swagger: ## Generates the Swagger documentation.
	@echo ''
	protoc -I=$(PROTOS_DIR) --openapiv2_out=$(DOCS_OUT_DIR) $(USERS_PROTO_FILE)

clean: ## Cleans the environment.
	@echo ''
	go clean -cache -modcache -testcache

test: ## Runs the tests.
	@echo ''
	go test ./... -cover