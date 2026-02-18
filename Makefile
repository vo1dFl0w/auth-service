.PHONY: genall ogen sqlc install-tools testunit testbench testintegration testall 

ogen:
	ogen --target ./internal/transport/http/httpgen --package httpgen --clean ./api/v1/openapi.yaml

sqlc:
	sqlc generate

mocks:
	mockery --log-level=debug

genall: ogen sqlc mocks

install-tools:
	go install -v github.com/ogen-go/ogen/cmd/ogen@latest
	go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
	go install github.com/vektra/mockery/v3@v3.6.1

testunit:
	go test ./internal/transport/http
	go test ./internal/usecase

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

testall: testunit testbench testintegration
