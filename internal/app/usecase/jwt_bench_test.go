package usecase_test

import (
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("jwt-secret-key-test")

func BenchmarkJWTSign(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   "userID",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute*15)),
		})
		_, _ = token.SignedString(jwtSecret)
	}
}

func BenchmarkJWTParse(b *testing.B) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "userID",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute*15)),
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
