.PHONY: test genall ogen sqlc install-tools

ogen:
	ogen --target ./internal/gen --package gen --clean ./api/v1/openapi.yaml

sqlc:
	sqlc generate -f db/sqlc.yaml

mocks:
	mockery --config ./internal/test/mocks/.mockery.yaml --log-level=debug

genall: ogen sqlc mocks

install-tools:
	go install -v github.com/ogen-go/ogen/cmd/ogen@latest
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
	go install github.com/vektra/mockery/v3@v3.6.1

test:
	go test ./internal/app/transport/http
	go test ./internal/app/usecase
	go test ./internal/test/integration_test