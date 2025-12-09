package usecase_test

import (
	"crypto/rand"
	"testing"

	"github.com/vo1dFl0w/auth-service/internal/app/usecase"
)

func BenchmarkHashRefreshToken32(b *testing.B) {
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
    
    b.ReportAllocs()
    b.ResetTimer()
    for b.Loop() {
        _ = usecase.HashRefreshTokenFunc(string(buf))
    }
}

func BenchmarkHashRefreshToken64(b *testing.B) {
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)

    b.ReportAllocs()
    b.ResetTimer()
    for b.Loop() {
        _ = usecase.HashRefreshTokenFunc(string(buf))
    }
}