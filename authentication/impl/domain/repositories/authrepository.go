package repositories

import (
	"context"
	"github.com/google/uuid"
)

type AuthRepository interface {
	GetPassword(context context.Context, email string) (string, error)
	GetUserByEmail(context context.Context, email string) (uuid.UUID, error)
}
