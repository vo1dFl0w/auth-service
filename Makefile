# ===== ogen =====
.PHONY: install-ogen ogen-gen

install-ogen:
	go install -v github.com/ogen-go/ogen/cmd/ogen@latest

ogen-gen:
	ogen --target ./internal/transport/http/httpgen --package httpgen --clean ./api/v1/openapi.yaml

# ===== sqlc =====
.PHONY: install-sqlc sqlc-gen

install-sqlc:
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

sqlc-gen:
	sqlc generate

# ===== lint =====
.PHONY: install-lint
install-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.10.1

.PHONY: lint
lint:
	golangci-lint run ./...

# ===== test =====
.PHONY: install-mockery mocks-gen testunit testbench testintegration
install-mockery:
	go install github.com/vektra/mockery/v3@v3.6.1

mocks-gen:
	mockery --log-level=debug

testunit:
	go test -race ./internal/transport/http
	go test -race ./internal/usecase

testbench:
	go test ./internal/usecase -bench=BenchmarkBcryptCost4
	go test ./internal/usecase -bench=BenchmarkBcryptCost10
	go test ./internal/usecase -bench=BenchmarkBcryptCost12
	go test ./internal/usecase -bench=BenchmarkJWTSign
	go test ./internal/usecase -bench=BenchmarkJWTParse
	go test ./internal/usecase -bench=BenchmarkHashRefreshToken32
	go test ./internal/usecase -bench=BenchmarkHashRefreshToken64

testintegration:
	go test ./internal/test/integration_test
	go test ./internal/test/integration_test -bench=BenchmarkFindUserByEmail
