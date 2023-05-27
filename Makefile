MOCKERY_VERSION = 2.23.1
PROTOC_VERSION = 1.30.0

install-tools:
	go install github.com/vektra/mockery/v2@v${MOCKERY_VERSION}
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_VERSION}

test-unit:
	go test -count=1 ./...

test-int-down:
	docker-compose -f ./infra/compose/docker-compose.yaml down

test-int-up:
	docker-compose -f ./infra/compose/docker-compose.yaml up -d

test-integration: test-int-up
	INTEGRATION=1 go test -count=1 ./...

test: test-integration

# Requires protoc installed (https://grpc.io/docs/protoc-installation/)
generate-proto:
	protoc --go_out=. --go_opt=paths=source_relative business/event/*.proto
