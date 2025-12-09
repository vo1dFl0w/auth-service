package usecase_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func BenchmarkBcryptCost4(b *testing.B) {
	benchBcryptCost(b, 4)
}

func BenchmarkBcryptCost10(b *testing.B) {
	benchBcryptCost(b, 10)
}

func BenchmarkBcryptCost12(b *testing.B) {
	benchBcryptCost(b, 12)
}

func benchBcryptCost(b *testing.B, cost int) {
	b.ReportAllocs()
	jwtSecret := []byte("jwt-secret-key-test")
	b.ResetTimer()
	for b.Loop() {
		_, _ = bcrypt.GenerateFromPassword(jwtSecret, cost)
	}
}
