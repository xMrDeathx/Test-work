package repositories

import (
	"context"
	"github.com/google/uuid"
)

type TokenRepository interface {
	ValidateRefreshToken(context context.Context, userID uuid.UUID, oldToken uuid.UUID) (bool, error)
	SaveToken(context context.Context, token uuid.UUID, userID uuid.UUID) error
	UpdateToken(context context.Context, token uuid.UUID, userID uuid.UUID) error
}
