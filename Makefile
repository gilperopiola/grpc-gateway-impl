# - GRPC
### - Gateway
##### - Implementation 
##### 
##### Made with love, @gilperopiola~

#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#
#        - Make Setup -        #
#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#

# Current git tag (or commit)
VERSION := $(shell git describe --tags --always --dirty)

# On 'make version', prints the version.
version:
	@echo $(VERSION)

# Path where the auto-generated swagger files will go.
DOCS_OUT_DIR := ./etc/docs

# Path where the auto-generated .pb.go files will go.
PBS_OUT_DIR := ./app/core/pbs

# Path where the .proto files are.
PROTOS_DIR := ./app/core/protos

# All .proto files in the PROTOS_DIR.
# Just top-level, not subfolders.
PROTO_FILES := $(shell find $(PROTOS_DIR) -maxdepth 1 -name "*.proto")

# Don't print unnecessary output.
MAKEFLAGS += --no-print-directory 

#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#
#       - Main Commands -      #
#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#

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

# On 'make test' executes tests.
test:
	go test ./... -cover

# On 'make generate', auto-generates the .pb.go files and the swagger
# based on the .proto files.
#
# A succesful 'make generate' command should output something like this:
#
# 	protoc -I=./protos 	--go_out=./pbs  --go-grpc_out=./pbs 	--go_opt=paths=so...... 	./protos/auth.proto ./protos/common.proto ./protos/users.proto
# 	protoc -I=./protos 	--grpc-gateway_out=./pbs 				--grpc-gateway_op......		./protos/auth.proto ./protos/common.proto ./protos/users.proto
# 	protoc -I=./protos 	--openapiv2_out=./etc/docs 											./protos/auth.proto ./protos/common.proto ./protos/users.proto
#
# - This example was shortened, and every "./app/core/" was replaced by just "./" for brevity.
generate:
	@'$(MAKE)' generate-pbs generate-swagger

#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#
#        Other Commands        #
#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#

# On 'make generate-pbs', auto-generates .pb.go files.
generate-pbs:
	protoc -I=$(PROTOS_DIR) --go_out=$(PBS_OUT_DIR) --go-grpc_out=$(PBS_OUT_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(PROTO_FILES)
	protoc -I=$(PROTOS_DIR) --grpc-gateway_out=$(PBS_OUT_DIR) --grpc-gateway_opt=paths=source_relative $(PROTO_FILES)
 
# On 'make generate-swagger', auto-generates swagger files.
generate-swagger:
	protoc -I=$(PROTOS_DIR) --openapiv2_out=$(DOCS_OUT_DIR) $(PROTO_FILES)

# On 'make push', adds + commits + pushes to master.
# On 'make push msg=":)"' adds a commit message.
push:
	git add .
	git commit -m "[@gilperopiola] - $(msg)"
	git push origin master

# On 'make clean', cleans the project.
clean:
	go clean -cache -modcache -testcache

# On 'make proinhanssr', enhances the project.
proinhanssr: 
	go run etc/tools/proinhanssr/proinhanssr.go

#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#
#           Resources          #
#-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~#

# Besides installing protoc and adding it to the PATH,
# to use the code auto-generation you'll also need to run:
#
# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway 
# go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 
# go install google.golang.org/protobuf/cmd/protoc-gen-go 
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
