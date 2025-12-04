package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/auth-service/internal/app/domain"
	"github.com/vo1dFl0w/auth-service/internal/app/repository"
	"github.com/vo1dFl0w/auth-service/internal/gen"
)

type PostgresTokenRepo struct {
	queries *gen.Queries
}

func NewPostgresTokenRepo(q *gen.Queries) *PostgresTokenRepo {
	return &PostgresTokenRepo{
		queries: q,
	}
}

func (r *PostgresTokenRepo) FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	ref, err := r.queries.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		} else {
			return nil, fmt.Errorf("failed to find refresh token: %w", err)
		}
	}

	return &domain.RefreshToken{
		UserID: ref.UserID,
		RefreshToken: ref.RefreshTokenHash,
		ExpiresAt: ref.ExpiresAt,
	}, nil
}

func (r *PostgresTokenRepo) SaveHashedRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.queries.SaveHashedRefreshToken(ctx, gen.SaveHashedRefreshTokenParams{
		UserID: userID,
		RefreshTokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresTokenRepo) DeleteRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	ref, err := r.queries.DeleteRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNoRowDeleted
		}
		return nil, err
	}
	
	return &domain.RefreshToken{
		UserID: ref.UserID,
		RefreshToken: ref.RefreshTokenHash,
		ExpiresAt: ref.ExpiresAt,
	}, nil
}