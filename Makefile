# 🔻 GRPC Gateway Implementation — Makefile 🔻

VERSION := $(shell git describe --tags --always --dirty) # Current git tag or commit
MAKEFLAGS += --no-print-directory # Don't print unnecessary output

version:
	@echo $(VERSION)

DOCS_OUT_DIR := ./docs # Path where the auto-generated swagger files will go
PBS_OUT_DIR := ./app/core/pbs # Path where the auto-generated .pb.go files will go
PROTOS_DIR := ./app/core/protos # Path where the .proto files are

# All .proto files in the PROTOS_DIR. Just top-level, not subfolders
PROTO_FILES := $(shell find $(PROTOS_DIR) -maxdepth 1 -name "*.proto")

# 🔻 Main Commands 🔻

all:
	@'$(MAKE)' clean generate test run

run:
	go mod tidy
	go run main.go

run-gen:
	go mod tidy
	go generate ./...
	go run main.go

test:
	go test ./... -cover

# Runs go generate — then auto-generates the .pb.go files and swagger based on the .protos
generate:
	go generate ./...
	@'$(MAKE)' generate-pbs generate-swagger

graph:
	./scripts/graph_dependencies.sh

# 🔻 Other Commands 🔻

# Adds, commits and pushes to master. 
# make push msg='msg'.
push:
	git add .
	git commit -m "[@gilperopiola] — $(msg)"
	git push origin master

generate-pbs:
	@echo "Proto files: $(PROTO_FILES)"
	protoc -I=$(PROTOS_DIR) --go_out=$(PBS_OUT_DIR) --go-grpc_out=$(PBS_OUT_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(PROTO_FILES)
	protoc -I=$(PROTOS_DIR) --grpc-gateway_out=$(PBS_OUT_DIR) --grpc-gateway_opt=paths=source_relative $(PROTO_FILES)
 
generate-swagger:
	protoc -I=$(PROTOS_DIR) --openapiv2_out=$(DOCS_OUT_DIR) $(PROTO_FILES)

# Cleans the project.
clean:
	go clean -cache -modcache -testcache

git-log:
	git log --oneline --graph --decorate --all

update-deps:
	go get -u ./...
	go mod tidy

# Enhances the project - ???
proinhanssr: 
	go run etc/tools/proinhanssr/proinhanssr.go

# 🔻 Resources 🔻

# Besides installing protoc and adding it to the path,
# to use the code auto-generation you'll also need to run:
#
# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway 
# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 
# go install google.golang.org/protobuf/cmd/protoc-gen-go 
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
