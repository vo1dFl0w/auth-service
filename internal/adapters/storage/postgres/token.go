package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/auth-service/internal/domain"
	"github.com/vo1dFl0w/auth-service/internal/repository"
	pggen "github.com/vo1dFl0w/auth-service/internal/adapters/storage/postgres/pggen"
)

type PostgresTokenRepo struct {
	queries *pggen.Queries
}

func NewPostgresTokenRepo(q *pggen.Queries) *PostgresTokenRepo {
	return &PostgresTokenRepo{
		queries: q,
	}
}

func (r *PostgresTokenRepo) FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	ref, err := r.queries.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, repository.ErrGatewayTimeout
		} else if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		} else {
			return nil, err
		}
	}

	return &domain.RefreshToken{
		UserID:       ref.UserID,
		RefreshToken: ref.RefreshTokenHash,
		ExpiresAt:    ref.ExpiresAt,
	}, nil
}

func (r *PostgresTokenRepo) SaveHashedRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.queries.SaveHashedRefreshToken(ctx, pggen.SaveHashedRefreshTokenParams{
		UserID:           userID,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        expiresAt,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return repository.ErrGatewayTimeout
		} else {
			return err
		}
	}

	return nil
}

func (r *PostgresTokenRepo) DeleteRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	ref, err := r.queries.DeleteRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, repository.ErrGatewayTimeout
		} else if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNoRowDeleted
		} else {
			return nil, err
		}
	}

	return &domain.RefreshToken{
		UserID:       ref.UserID,
		RefreshToken: ref.RefreshTokenHash,
		ExpiresAt:    ref.ExpiresAt,
	}, nil
}
