package postgres

import (
	"database/sql"

	"github.com/vo1dFl0w/auth-service/internal/repository"
	pggen "github.com/vo1dFl0w/auth-service/internal/adapters/storage/postgres/pggen"
)

type Storage struct {
	db        *sql.DB
	authRepo  repository.AuthRepository
	tokenRepo repository.TokenRepository
}

func New(db *sql.DB) *Storage {
	q := pggen.New(db)

	return &Storage{
		db:        db,
		authRepo:  NewPostgresAuthRepo(q),
		tokenRepo: NewPostgresTokenRepo(q),
	}
}

func (s *Storage) Auth() repository.AuthRepository {
	return s.authRepo
}

func (s *Storage) Token() repository.TokenRepository {
	return s.tokenRepo
}
