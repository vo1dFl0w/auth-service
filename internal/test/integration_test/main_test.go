package integrationtest

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	migrate "github.com/golang-migrate/migrate/v4"
	migratepq "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	postgresUser     = "test"
	postgresPassword = "test"
	postgresDB       = "testdb"
	postgresPort     = "5432/tcp"

	migrationsSource = "file://../../../migrations"
)

var (
	DSN    string
	TestDB *sql.DB
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		fmt.Printf("env not found: %s", err)
		os.Exit(1)
	}

	if os.Getenv("INTEGRATION") != "1" {
		fmt.Println("integration tests skipped; set INTEGRATION=1 to run")
		os.Exit(1)
	}

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14-alpine",
		ExposedPorts: []string{postgresPort},
		Env: map[string]string{
			"POSTGRES_USER":     postgresUser,
			"POSTGRES_PASSWORD": postgresPassword,
			"POSTGRES_DB":       postgresDB,
		},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", postgresUser, postgresPassword, port.Port(), postgresDB)
		}).WithStartupTimeout(time.Second * 60),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Printf("generic container: %s", err)
		os.Exit(1)
	}

	host, err := c.Host(ctx)
	if err != nil {
		log.Printf("host: %s", err)
		panic(err)
	}

	p, err := c.MappedPort(ctx, "5432")
	if err != nil {
		log.Printf("mapped port: %s", err)
		panic(err)
	}

	DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPassword, host, p.Port(), postgresDB)

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		log.Printf("open db: %s", err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Printf("ping db: %s", err)
		panic(err)
	}

	driver, err := migratepq.WithInstance(db, &migratepq.Config{})
	if err != nil {
		log.Printf("with instance: %s", err)
		panic(err)
	}

	migration, err := migrate.NewWithDatabaseInstance(migrationsSource, "postgres", driver)
	if err != nil {
		log.Printf("new with database instance: %s", err)
		panic(err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("migration up: %s", err)
		panic(err)
	}

	TestDB = db

	code := m.Run()

	if err := c.Terminate(ctx); err != nil {
		log.Printf("terminate container: %s", err)
		panic(err)
	}

	os.Exit(code)
}
