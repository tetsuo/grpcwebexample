GO := go
GOFLAGS ?= -mod=readonly -ldflags "-s -w"
OUT_MODULE := example.com/grpcwebexample
PROTO_DIR := proto
PBSCHEMAS := $(shell find $(PROTO_DIR) -name '*.proto')
JS_OUT_DIR := web/gen

build: distclean
	CGO_ENABLED=0 $(GO) build -trimpath -buildvcs=false $(GOFLAGS) -o bin/ ./...

distclean:
	@rm -fr bin
	@rm -fr build
	@rm -fr dist

clean:
	@$(GO) clean

gofmt:
	@mkdir -p output
	@rm -f output/lint.log

	gofmt -d -s . 2>&1 | tee output/lint.log

	@[ ! -s output/lint.log ]

	@rm -fr output

pbclean:
	@find . -name *.pb.go -delete
	@find . -name *_pb.js -delete

pb: pbclean
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=. \
		--go_opt=module=$(OUT_MODULE) \
		--go-grpc_out=. \
		--go-grpc_opt=module=$(OUT_MODULE) \
		--js_out=import_style=commonjs,binary:$(JS_OUT_DIR) \
		--grpc-web_out=import_style=commonjs,mode=grpcweb:$(JS_OUT_DIR) \
		$(PBSCHEMAS)

tidy: pb
	@$(GO) mod tidy -v

.PHONY: build
