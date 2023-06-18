include .env
export $(shell sed 's/=.*//' .env)

LOCAL_BIN:=$(CURDIR)/bin

install-go-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.2

generate:
	make generate-chat-api

generate-chat-api:
	mkdir -p pkg/chat_v1
	protoc --proto_path api/chat_v1 --proto_path vendor.protogen \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/chat_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	api/chat_v1/chat.proto

CERTS_DIR=./certs
CA_KEY=$(CERTS_DIR)/ca.key
CA_CRT=$(CERTS_DIR)/ca.crt
KEY_SIZE=4096
DAYS=1024
SHA=sha256
CLIENT_DN="/C=US/ST=California/L=San Francisco/O=MyOrg/OU=MyDepartment/CN=localhost"
CLIENT_KEY=$(CERTS_DIR)/client.key
CLIENT_CSR=$(CERTS_DIR)/client.csr
CLIENT_CRT=$(CERTS_DIR)/client.crt

generate-client-key:
	openssl genrsa -out $(CLIENT_KEY) $(KEY_SIZE)

generate-client-csr:
	openssl req -new -key $(CLIENT_KEY) -out $(CLIENT_CSR) -subj $(CLIENT_DN)

generate-client-crt:
	openssl x509 -req -in $(CLIENT_CSR) -CA $(CA_CRT) -CAkey $(CA_KEY) -CAcreateserial \
 		-out $(CLIENT_CRT) -days $(DAYS) -$(SHA) -extensions client_ext -extfile $(CERTS_DIR)/v3.ext

generate-openssl-keys:
	make generate-client-key
	make generate-client-csr
	make generate-client-crt


vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi