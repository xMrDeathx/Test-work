package services

import (
	"TestWork/authentication/impl/app/commands"
	"TestWork/authentication/impl/domain/repositories"
	"TestWork/authentication/impl/domain/services"
	"context"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

var ErrUserWithTokenNotFound = errors.New("user with token not found")

func NewTokenService(repository repositories.TokenRepository) services.TokenService {
	return &tokenService{repository: repository}
}

type tokenService struct {
	repository repositories.TokenRepository
}

// Обработка запроса на генерацию токенов
func (service *tokenService) GenerateTokens(context context.Context, userID uuid.UUID, requestIP string) (commands.TokensResult, error) {
	accessToken, err := createAccessToken(userID, requestIP)
	if err != nil {
		return commands.TokensResult{}, err
	}

	refreshToken := uuid.New()
	//Сохранение refresh token в базу данных
	err = service.repository.SaveToken(context, refreshToken, userID)
	if err != nil {
		return commands.TokensResult{}, err
	}

	return commands.TokensResult{
		AccessToken:  accessToken,
		RefreshToken: encodeRefreshToken(refreshToken),
	}, nil
}

// Обработка запроса на обновление токенов
func (service *tokenService) RefreshTokens(context context.Context, userID uuid.UUID, token string, requestIP string) (commands.TokensResult, error) {
	domainToken, err := decodeRefreshToken(token)
	if err != nil {
		return commands.TokensResult{}, err
	}

	tokenValidated, err := service.repository.ValidateRefreshToken(context, userID, domainToken)
	if err != nil {
		return commands.TokensResult{}, err
	}

	if !tokenValidated {
		return commands.TokensResult{}, errors.New("Invalid refresh token")
	}

	newAccessToken, err := createAccessToken(userID, requestIP)
	if err != nil {
		return commands.TokensResult{}, err
	}

	newRefreshToken := uuid.New()
	err = service.repository.UpdateToken(context, newRefreshToken, userID)
	if err != nil {
		return commands.TokensResult{}, err
	}

	return commands.TokensResult{
		AccessToken:  newAccessToken,
		RefreshToken: encodeRefreshToken(newRefreshToken),
	}, nil
}

func createAccessToken(userID uuid.UUID, requestIP string) (string, error) {
	secretKey := "secret"
	//Создание payload для access token
	accessTokenPayload := jwt.MapClaims{
		"exp":       time.Now().Add(30 * time.Minute).UnixMilli(),
		"user":      userID,
		"requestIP": requestIP,
	}

	//Создание access token
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenPayload).SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// base64 формат для refresh token
func encodeRefreshToken(token uuid.UUID) string {
	return base64.StdEncoding.EncodeToString([]byte(token.String()))
}

func decodeRefreshToken(token string) (uuid.UUID, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return uuid.UUID{}, err
	}
	result, err := uuid.ParseBytes(decodedToken)
	if err != nil {
		return uuid.UUID{}, err
	}

	return result, nil
}
