package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/vo1dFl0w/auth-service/internal/domain"
	"github.com/vo1dFl0w/auth-service/internal/repository"
	pggen "github.com/vo1dFl0w/auth-service/internal/adapters/storage/postgres/pggen"
)

type PostgresAuthRepo struct {
	queries *pggen.Queries
}

func NewPostgresAuthRepo(q *pggen.Queries) *PostgresAuthRepo {
	return &PostgresAuthRepo{
		queries: q,
	}
}

func (r *PostgresAuthRepo) CreateUser(ctx context.Context, email string, passwordHash string) (*domain.User, error) {
	var pqErr *pq.Error

	u, err := r.queries.CreateUser(ctx, pggen.CreateUserParams{
		Email: email,
		PasswordHash: sql.NullString{
			String: passwordHash,
			Valid:  passwordHash != "",
		},
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

func (r *PostgresAuthRepo) CreateOAuthUser(ctx context.Context, email string, providerID string, oauthProvider string) (*domain.User, error) {
	u, err := r.queries.CreateOAuthUser(ctx, pggen.CreateOAuthUserParams{
		Email:         email,
		ProviderID:    sql.NullString{String: providerID, Valid: providerID != ""},
		OauthProvider: sql.NullString{String: oauthProvider, Valid: oauthProvider != ""},
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, repository.ErrGatewayTimeout
		} else {
			return nil, err
		}
	}

	return &domain.User{
		UserID:        u.UserID,
		Email:         u.Email,
		CreatedAt:     u.CreatedAt,
		IsActive:      u.IsActive,
		ProviderID:    u.ProviderID.String,
		OAuthProvider: u.OauthProvider.String,
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
		UserID:        u.UserID,
		Email:         u.Email,
		CreatedAt:     u.CreatedAt,
		IsActive:      u.IsActive,
		ProviderID:    u.ProviderID.String,
		OAuthProvider: u.OauthProvider.String,
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
		PasswordHash: u.PasswordHash.String,
		CreatedAt:    u.CreatedAt,
		IsActive:     u.IsActive,
	}, nil
}

func (r *PostgresAuthRepo) FindUserByProviderID(ctx context.Context, providerID string, oauthProvider string) (*domain.User, error) {
	u, err := r.queries.FindUserByProviderID(ctx, pggen.FindUserByProviderIDParams{
		ProviderID: sql.NullString{
			String: providerID,
			Valid:  providerID != "",
		}, OauthProvider: sql.NullString{
			String: oauthProvider,
			Valid:  oauthProvider != "",
		},
	})
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
		UserID:        u.UserID,
		Email:         u.Email,
		CreatedAt:     u.CreatedAt,
		IsActive:      u.IsActive,
		ProviderID:    u.ProviderID.String,
		OAuthProvider: u.OauthProvider.String,
	}, nil
}
