package integrationtest

import (
	"context"
	"os"
	"testing"

	db "github.com/vo1dFl0w/auth-service/internal/gen"
)

func BenchmarkFindUserByEmail(b *testing.B) {
    if os.Getenv("INTEGRATION") != "1" || os.Getenv("BENCHMARK") != "1" {
		b.Skip("integration bench-test skipped; set INTEGRATION=1 and BENCHMARK=1 to run")
	}
    ctx := context.Background()

    email := "user@example.com"
    q := db.New(TestDB)
    _, _ = q.CreateUser(ctx, db.CreateUserParams{Email: email, PasswordHash: "hash"})
    b.ReportAllocs()
    b.ResetTimer()

    for b.Loop() {
        _, err := q.FindUserByEmail(ctx, email)
        if err != nil {
            b.Fatalf("find failed: %v", err)
        }
    }
}

