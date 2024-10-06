package commands

import "github.com/google/uuid"

type LoginResult struct {
	Tokens TokensResult
	UserID uuid.UUID
}
