package services

import (
	"TestWork/authentication/impl/app/commands"
	"context"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateTokens(context context.Context, userID uuid.UUID, requestIP string) (commands.TokensResult, error)
	RefreshTokens(context context.Context, userID uuid.UUID, token string, requestIP string) (commands.TokensResult, error)
}
