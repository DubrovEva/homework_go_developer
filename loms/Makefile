BINDIR=${CURDIR}/bin
PACKAGE=route256/loms

LOMS_PROTO_PATH:="api/loms/v1"

.PHONY: bindir
bindir:
	mkdir -p ${BINDIR}

.PHONY: build
build:
	echo "build loms"

.PHONY: run-migrations
run-migrations:
	echo "run migrations"

# Устанавливаем v1 описания protoc-gen-openapiv2/options
vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor-proto/grpc-ecosystem && \
 	cd vendor-proto/grpc-ecosystem && \
	git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
	git checkout
	mkdir -p vendor-proto/protoc-gen-openapiv2
	mv vendor-proto/grpc-ecosystem/protoc-gen-openapiv2/options vendor-proto/protoc-gen-openapiv2
	rm -rf vendor-proto/grpc-ecosystem

.PHONY: protoc-generate
protoc-generate: pkg/api/loms/v1/api.pb.gw.go pkg/api/loms/v1/api.pb.go pkg/api/loms/v1/order.pb.go

pkg/api/loms/v1/%.pb.go: api/loms/v1/%.proto vendor-proto/protoc-gen-openapiv2/options ${BINDIR}/protoc-gen-go ${BINDIR}/protoc-gen-validate ${BINDIR}/goose ${BINDIR}/sqlc pkg/api/loms/v1
	protoc \
    	-I ${LOMS_PROTO_PATH} \
    	-I vendor-proto \
    	--plugin=protoc-gen-go=${BINDIR}/protoc-gen-go \
    	--go_out pkg/${LOMS_PROTO_PATH} \
    	--go_opt paths=source_relative \
    	--plugin=protoc-gen-validate=${BINDIR}/protoc-gen-validate \
    	--validate_out="lang=go,paths=source_relative:pkg/api/loms/v1" \
    	$<
	go mod tidy

pkg/api/loms/v1/api.pb.gw.go: vendor-proto/protoc-gen-openapiv2/options ${BINDIR}/protoc-gen-go-grpc ${BINDIR}/protoc-gen-grpc-gateway ${BINDIR}/protoc-gen-openapiv2  ${BINDIR}/goose ${BINDIR}/sqlc pkg/api/loms/v1 pkg/api/loms/v1
	mkdir -p api/openapiv2
	protoc \
		-I ${LOMS_PROTO_PATH} \
		-I vendor-proto \
		--plugin=protoc-gen-go-grpc=${BINDIR}/protoc-gen-go-grpc \
		--go-grpc_out pkg/${LOMS_PROTO_PATH} \
		--go-grpc_opt paths=source_relative \
		--plugin=protoc-gen-validate=${BINDIR}/protoc-gen-validate \
		--validate_out="lang=go,paths=source_relative:pkg/api/loms/v1" \
		--plugin=protoc-gen-grpc-gateway=${BINDIR}/protoc-gen-grpc-gateway \
		--grpc-gateway_out pkg/${LOMS_PROTO_PATH} \
		--grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
		--plugin=protoc-gen-openapiv2=${BINDIR}/protoc-gen-openapiv2 \
		--openapiv2_out api/openapiv2 \
		--openapiv2_opt logtostderr=true \
		${LOMS_PROTO_PATH}/api.proto
	go mod tidy

${BINDIR}/protoc-gen-go:
	GOBIN=$(BINDIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1

${BINDIR}/protoc-gen-go-grpc:
	GOBIN=$(BINDIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

${BINDIR}/protoc-gen-validate:
	GOBIN=$(BINDIR) go install github.com/envoyproxy/protoc-gen-validate@v1.0.4

${BINDIR}/protoc-gen-grpc-gateway:
	GOBIN=$(BINDIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1

${BINDIR}/protoc-gen-openapiv2:
	GOBIN=$(BINDIR) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1

${BINDIR}/goose:
	GOBIN=$(BINDIR) go install github.com/pressly/goose/v3/cmd/goose@v3.24.1

${BINDIR}/sqlc:
	GOBIN=$(BINDIR) go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.28.0


create-migration-file:
	$(BINDIR)/goose -dir migrations create -s $(NAME) sql

up-migrations:
	$(BINDIR)/goose -dir migrations postgres "postgresql://loms-user:loms-password@127.0.0.1:5432/loms_db?sslmode=disable" up

compile-sql:
	$(BINDIR)/sqlc generate

pkg/api/loms/v1:
	mkdir -p pkg/api/loms/v1

.PHONY: clean
clean:
	rm -rf vendor-proto
	rm -rf api/openapiv2
