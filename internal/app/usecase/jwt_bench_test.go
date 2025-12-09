package usecase_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("jwt-secret-key-test")

func BenchmarkJWTSign(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Subject:   "userID",
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		})
		_, _ = token.SignedString(jwtSecret)
	}
}

func BenchmarkJWTParse(b *testing.B) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   "userID",
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	})

	s, _ := token.SignedString(jwtSecret)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = jwt.Parse(s, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
	}
}
