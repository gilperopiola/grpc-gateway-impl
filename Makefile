###---------------------------------------------###
###                                             ###
###       - GRPC Gateway Implementation -       ###
###                                             ###
###--------------------------- by @gilperopiola ###

#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#
#          - Set up -          #
#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#

VERSION := $(shell git describe --tags --always --dirty) # Current git tag or commit hash.

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

version:
	@echo $(VERSION)

.DEFAULT_GOAL := help
MAKEFLAGS += --no-print-directory # Don't print unnecessary output.

#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#

DOCS_OUT_DIR := ./etc/docs
PBS_OUT_DIR := ./app/core/pbs
PROTOS_DIR := ./app/core/protos

COMMON_PROTO := $(PROTOS_DIR)/common.proto
AUTH_PROTO := $(PROTOS_DIR)/auth.proto
GROUPS_PROTO := $(PROTOS_DIR)/groups.proto
USERS_PROTO := $(PROTOS_DIR)/users.proto

#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#
#       - Main Targets -       #
#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#

all:
ifeq ($(fast),)
	@'$(MAKE)' clean generate test run
else 
	@'$(MAKE)' generate run
endif

run:
	go mod tidy
	go run main.go

test:
	go test ./... -race -cover

push:
	git add .
	git commit -m "[@gilperopiola] - $(msg)"
	git push origin master

#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#
#      Secondary Targets       #
#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#

generate:
	@'$(MAKE)' generate-pbs generate-swagger

# For this command remember to install protoc (and add to path).
# Also:
# -> go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway 
# -> go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 
# -> go install google.golang.org/protobuf/cmd/protoc-gen-go 
# -> go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
generate-pbs:
	protoc -I=$(PROTOS_DIR) --go_out=$(PBS_OUT_DIR) --go-grpc_out=$(PBS_OUT_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(AUTH_PROTO) $(GROUPS_PROTO) $(USERS_PROTO) $(COMMON_PROTO)
	protoc -I=$(PROTOS_DIR) --grpc-gateway_out=$(PBS_OUT_DIR) --grpc-gateway_opt=paths=source_relative $(AUTH_PROTO) $(GROUPS_PROTO) $(USERS_PROTO) $(COMMON_PROTO)
 
generate-swagger:
	protoc -I=$(PROTOS_DIR) --openapiv2_out=$(DOCS_OUT_DIR)  $(AUTH_PROTO) $(GROUPS_PROTO) $(USERS_PROTO) $(COMMON_PROTO)

clean:
	go clean -cache -modcache -testcache

proinhanssr: 
	go run etc/tools/proinhanssr/proinhanssr.go


