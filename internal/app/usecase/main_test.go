package usecase

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var (
	envPath = "../../../.env"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(envPath); err != nil {
		fmt.Printf("env not found: %s", err)
		os.Exit(1)
	}

	if os.Getenv("BENCHMARK") != "1" {
		fmt.Println("bench-test skipped; set BENCHMARK=1 to run")
		os.Exit(1)
	}

	code := m.Run()

	os.Exit(code)
}
