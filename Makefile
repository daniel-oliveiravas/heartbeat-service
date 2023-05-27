MOCKERY_VERSION = 2.23.1

install-tools:
	go install github.com/vektra/mockery/v2@v${MOCKERY_VERSION}

test-unit:
	go test -count=1 ./...

test-integration:
	docker-compose -f ./infra/compose/docker-compose.yaml up -d
	INTEGRATION=1 go test -count=1 ./...

test: test-integration
