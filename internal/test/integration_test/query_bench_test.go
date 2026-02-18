package integrationtest

import (
	"context"
	"database/sql"
	"os"
	"testing"

	pggen "github.com/vo1dFl0w/auth-service/internal/adapters/storage/postgres/pggen"
)

func BenchmarkFindUserByEmail(b *testing.B) {
    if os.Getenv("INTEGRATION") != "1" || os.Getenv("BENCHMARK") != "1" {
		b.Skip("integration bench-test skipped; set INTEGRATION=1 and BENCHMARK=1 to run")
	}
    ctx := context.Background()

    email := "user@example.com"
    q := pggen.New(TestDB)
    _, _ = q.CreateUser(ctx, pggen.CreateUserParams{Email: email, PasswordHash: sql.NullString{String: "hash", Valid: "hash" != ""},})
    b.ReportAllocs()
    b.ResetTimer()

    for b.Loop() {
        _, err := q.FindUserByEmail(ctx, email)
        if err != nil {
            b.Fatalf("find failed: %v", err)
        }
    }
}

