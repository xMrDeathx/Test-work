package tokenmapper

import (
	"TestWork/authentication/impl/app/commands"
	"TestWork/authentication/impl/domain/model"
	"github.com/google/uuid"
	"time"
)

func NewDomainSession(userID uuid.UUID, requestIP string, token []byte) model.Session {
	return model.Session{
		ID:     uuid.New(),
		UserID: userID,
		UserIP: requestIP,
		Token:  NewDomainRefreshToken(token),
	}
}

func NewDomainRefreshToken(token []byte) model.RefreshToken {
	return model.RefreshToken{
		Token:     token,
		ExpiresIn: time.Now().Add(3 * 24 * time.Hour).UnixMilli(),
	}
}

func NewTokensResultFromEntity(accessToken string, refreshToken string) commands.TokensResult {
	return commands.TokensResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
