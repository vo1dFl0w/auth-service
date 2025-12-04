package postgres

import (
	"database/sql"
	"sync"

	"github.com/vo1dFl0w/auth-service/internal/app/repository"
	"github.com/vo1dFl0w/auth-service/internal/gen"
)

type Storage struct {
	db        *sql.DB
	authOnce  sync.Once
	authRepo  repository.AuthRepository
	tokenOnce sync.Once
	tokenRepo repository.TokenRepository
}

func New(db *sql.DB) *Storage {
	q := gen.New(db)

	return &Storage{
		db:        db,
		authRepo:  NewPostgresAuthRepo(q),
		tokenRepo: NewPostgresTokenRepo(q),
	}
}

func (s *Storage) Auth() repository.AuthRepository {
	s.authOnce.Do(func() {
		q := gen.New(s.db)
		s.authRepo = NewPostgresAuthRepo(q)
	})
	return s.authRepo
}

func (s *Storage) Token() repository.TokenRepository {
	s.tokenOnce.Do(func() {
		q := gen.New(s.db)
		s.tokenRepo = NewPostgresTokenRepo(q)
	})
	return s.tokenRepo
}
