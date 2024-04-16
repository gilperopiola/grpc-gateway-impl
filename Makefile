###---------------------------------------------###
###                                             ###
###       - gRPC Gateway Implementation -       ###
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
PBS_OUT_DIR := ./app/pbs
PROTOS_DIR := ./app/protos
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
	@echo ''
	go mod tidy
	go run main.go

test:
	@echo ''
	go test ./... -race -cover

push:
	@echo ''
	git add .
	git commit -m "[@gilperopiola] - $(msg)"
	git push origin master

#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#
#      Secondary Targets       #
#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#

generate:
	@'$(MAKE)' generate-pbs generate-swagger

generate-pbs:
	@echo ''
	protoc -I=$(PROTOS_DIR) --go_out=$(PBS_OUT_DIR) --go-grpc_out=$(PBS_OUT_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(USERS_PROTO)
	protoc -I=$(PROTOS_DIR) --grpc-gateway_out=$(PBS_OUT_DIR) --grpc-gateway_opt=paths=source_relative $(USERS_PROTO)
 
generate-swagger:
	@echo ''
	protoc -I=$(PROTOS_DIR) --openapiv2_out=$(DOCS_OUT_DIR) $(USERS_PROTO)

clean:
	@echo ''
	go clean -cache -modcache -testcache

proinhanssr: 
	@echo ''
	go run etc/tools/proinhanssr/proinhanssr.go


