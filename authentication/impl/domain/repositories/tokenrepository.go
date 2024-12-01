package repositories

import (
	"TestWork/authentication/impl/domain/model"
	"context"
	"github.com/google/uuid"
)

type TokenRepository interface {
	GetToken(context context.Context, userID uuid.UUID) (model.RefreshToken, string, error)
	SaveToken(context context.Context, session model.Session) error
	UpdateToken(context context.Context, oldToken []byte, newToken model.RefreshToken, requestIP string) error
}
