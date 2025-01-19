# ⭐ —— GRPC Gateway Implementation

NAME 			:= grpc-gateway-impl
VERSION 		:= v1.1.2

DOCS_OUT_DIR 	:= ./docs 				# Where the auto-generated swagger files will go
PBS_OUT_DIR 	:= ./app/core/pbs 		# Where the auto-generated .pb.go files will go
PROTOS_DIR 		:= ./app/core/protos 	# Where the .proto files are

PROTO_FILES 	:= $(shell find $(PROTOS_DIR) -maxdepth 1 -name "*.proto")

### ——> make
all:
	@'$(MAKE)' install
	@'$(MAKE)' walk

### ——> make run
run:
	@echo ">>> Running $(NAME) $(VERSION)..."
	go run main.go

### ——> make walk
# Like run, but slower as it also regenerates code
walk:
	go mod tidy
	@'$(MAKE)' generate
	@'$(MAKE)' run

### ——> make test
test:
	go test ./... -cover

### ——> make generate
# Runs go generate for some stuff
# Also auto-generates the .pb.go files and swagger based on the .protos
generate:
	@echo ">>> Generating code for $(NAME) $(VERSION)..."
	go generate ./...

	@echo ">>> And now based on the .protos..."
	protoc -I=$(PROTOS_DIR) --go_out=$(PBS_OUT_DIR) --go-grpc_out=$(PBS_OUT_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(PROTO_FILES)
	protoc -I=$(PROTOS_DIR) --grpc-gateway_out=$(PBS_OUT_DIR) --grpc-gateway_opt=paths=source_relative $(PROTO_FILES)
	protoc -I=$(PROTOS_DIR) --openapiv2_out=$(DOCS_OUT_DIR) $(PROTO_FILES)

### ——> make build
build:
	@echo ">>> Building $(NAME) $(VERSION)..."
	go build -ldflags="-s -w" -trimpath -o $(NAME).exe main.go

### ——> make graph
graph:
	@echo ">>> Graphing dependencies for $(NAME) $(VERSION)..."
	./scripts/graph_dependencies.sh

### ——> make install
# To be able to run the code auto-generation
install:
	@echo ">>> Installing dependencies for $(NAME) $(VERSION)..."
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

### ——> make push (optional msg='msg')
push:
	git add .
	git commit -m "[$(VERSION)] <3 $(msg)"
	git push origin master

### ——> make clean
clean:
	go clean -cache -modcache -testcache

### ——> make git-log
git-log:
	git log --oneline --graph --decorate --all

### ——> make update-deps
update-deps:
	go get -u ./...
	go mod tidy

### ——> make proinhanssr
# Enhances the project -> ???
proinhanssr: 
	go run etc/tools/proinhanssr/proinhanssr.go