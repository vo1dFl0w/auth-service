package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/vo1dFl0w/auth-service/internal/app/domain"
	"github.com/vo1dFl0w/auth-service/internal/app/repository"
	"github.com/vo1dFl0w/auth-service/internal/gen"
)

type PostgresAuthRepo struct {
	queries *gen.Queries
}

func NewPostgresAuthRepo(q *gen.Queries) *PostgresAuthRepo {
	return &PostgresAuthRepo{
		queries: q,
	}
}

func (r *PostgresAuthRepo) CreateUser(ctx context.Context, email string, passwordHash string) (*domain.User, error) {
	var pqErr *pq.Error

	u, err := r.queries.CreateUser(ctx, gen.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, repository.ErrGatewayTimeout
		} else if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return nil, repository.ErrEmailAlreadyExists
			}
		} else {
			return nil, err
		}
	}

	return &domain.User{
		UserID:    u.UserID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		IsActive:  u.IsActive,
	}, nil
}

func (r *PostgresAuthRepo) GetUserInfo(ctx context.Context, user_id uuid.UUID) (*domain.User, error) {
	u, err := r.queries.GetUserInfo(ctx, user_id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, repository.ErrGatewayTimeout
		} else if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		} else {
			return nil, err
		}
	}

	return &domain.User{
		UserID:    u.UserID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		IsActive:  u.IsActive,
	}, nil
}

func (r *PostgresAuthRepo) FindUserByEmail(ctx context.Context, email string) (*domain.UserWithPassword, error) {
	u, err := r.queries.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, repository.ErrGatewayTimeout
		} else if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		} else {
			return nil, err
		}
	}

	return &domain.UserWithPassword{
		UserID:       u.UserID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		IsActive:     u.IsActive,
	}, nil
}
