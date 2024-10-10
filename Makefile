# 🔻 GRPC Gateway Implementation 🔻

# 🔽 Make Setup 🔽

# Current git tag (or commit)
VERSION := $(shell git describe --tags --always --dirty)

# On 'make version', prints the version.
version:
	@echo $(VERSION)

# Path where the auto-generated swagger files will go.
DOCS_OUT_DIR := ./docs

# Path where the auto-generated .pb.go files will go.
PBS_OUT_DIR := ./app/core/pbs

# Path where the .proto files are.
PROTOS_DIR := ./app/core/protos

# All .proto files in the PROTOS_DIR.
# Just top-level, not subfolders.
PROTO_FILES := $(wildcard $(PROTOS_DIR)/*.proto)

# Don't print unnecessary output.
MAKEFLAGS += --no-print-directory 

# 🔻 Main Commands 🔻

# On 'make all --fast', auto-generates code + runs program.
# On 'make all', cleans the project + auto-generates code + runs tests + runs program.
all:
ifeq ($(fast),)
	@'$(MAKE)' generate run
else 
	@'$(MAKE)' clean generate test run
endif

# On 'make run' tidies the go modules and runs the program.
run:
	go mod tidy
	go run main.go

rungen:
	go mod tidy
	go generate ./...
	go run main.go

# On 'make test' executes tests.
test:
	go test ./... -cover

# On 'make generate', auto-generates the .pb.go files and the swagger
# based on the .proto files. Also runs scripts called by go generate.
#
# A succesful 'make generate' command should output something like this:
#
# 	protoc -I=./protos 	--go_out=./pbs  --go-grpc_out=./pbs 	--go_opt=paths=so...... 	./protos/auth.proto ./protos/common.proto ./protos/users.proto
# 	protoc -I=./protos 	--grpc-gateway_out=./pbs 				--grpc-gateway_op......		./protos/auth.proto ./protos/common.proto ./protos/users.proto
# 	protoc -I=./protos 	--openapiv2_out=./etc/docs 											./protos/auth.proto ./protos/common.proto ./protos/users.proto
#
# - This example was shortened, and every "./app/core/" was replaced by just "./" for brevity.
generate:
	go generate ./...
	@'$(MAKE)' generate-pbs generate-swagger

# 🔻 Other Commands 🔻

# On 'make clean', cleans the project.
clean:
	go clean -cache -modcache -testcache

# On 'make generate-pbs', auto-generates .pb.go files.
generate-pbs:
	protoc -I=$(PROTOS_DIR) --go_out=$(PBS_OUT_DIR) --go-grpc_out=$(PBS_OUT_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(PROTO_FILES)
	protoc -I=$(PROTOS_DIR) --grpc-gateway_out=$(PBS_OUT_DIR) --grpc-gateway_opt=paths=source_relative $(PROTO_FILES)
 
# On 'make generate-swagger', auto-generates swagger files.
generate-swagger:
	protoc -I=$(PROTOS_DIR) --openapiv2_out=$(DOCS_OUT_DIR) $(PROTO_FILES)

# On 'make push', adds, commits, and pushes to master.
# Add msg="something" to add a custom commit message.
push:
	git add .
	git commit -m "[@gilperopiola] - $(msg)"
	git push origin master

git-log:
	git log --oneline --graph --decorate --all

# On 'make graph', generates a graph of the dependencies,
# in .dot and .png formats.
#
# Oh yes I totally understand this code, so intuitive.
graph:
	echo "digraph dependencies {" > external_deps.dot
	go mod graph | awk '{print "\"" $$1 "\" -> \"" $$2 "\";"}' >> external_deps.dot
	echo "}" >> external_deps.dot
	dot -Tpng external_deps.dot -o external_deps.png

	echo "digraph dependencies {" > internal_deps.dot
	go list -f '{{.ImportPath}} {{join .Imports " "}}' ./... | awk ' \
	{ \
	  for (i = 2; i <= NF; i++) { \
	    print "\"" $$1 "\" -> \"" $$i "\""; \
	  } \
	}' >> internal_deps.dot
	echo "}" >> internal_deps.dot
	dot -Tpng internal_deps.dot -o internal_deps.png

update-deps:
	go get -u ./...
	go mod tidy

# On 'make proinhanssr', enhances the project.
proinhanssr: 
	go run etc/tools/proinhanssr/proinhanssr.go

# 🔻 Resources 🔻

# Besides installing protoc and adding it to the PATH,
# to use the code auto-generation you'll also need to run:
#
# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway 
# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 
# go install google.golang.org/protobuf/cmd/protoc-gen-go 
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
