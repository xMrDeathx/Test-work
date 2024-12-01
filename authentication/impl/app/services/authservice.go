package services

import (
	"TestWork/authentication/impl/app/commands"
	"TestWork/authentication/impl/app/mapper/loginmapper"
	"TestWork/authentication/impl/domain/repositories"
	"TestWork/authentication/impl/domain/services"
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

func NewAuthService(repository repositories.AuthRepository, tokenService services.TokenService) services.AuthService {
	return &authService{
		repository:   repository,
		tokenService: tokenService,
	}
}

type authService struct {
	repository   repositories.AuthRepository
	tokenService services.TokenService
}

// Обработка запроса на логинацию
func (service *authService) Login(context context.Context, loginData commands.LoginCommand, requestIP string) (commands.LoginResult, error) {
	domainLoginData := loginmapper.NewLoginDataToDomainLoginData(loginData)

	//Получение пароля пользователя из базы данных по логину (эл. почте)
	password, err := service.repository.GetPassword(context, domainLoginData.Email)
	if err != nil {
		return commands.LoginResult{}, err
	}
	//Проверка правильности пароля с тем, который был введён
	if password != domainLoginData.Password {
		return commands.LoginResult{}, ErrUserNotFound
	}
	//Получение uuid пользователя по логину
	userID, err := service.repository.GetUserByEmail(context, domainLoginData.Email)
	if err != nil {
		return commands.LoginResult{}, err
	}

	//Генерация токенов для пользователя
	tokens, err := service.tokenService.GenerateTokens(context, userID, requestIP)
	if err != nil {
		return commands.LoginResult{}, err
	}

	return commands.LoginResult{
		Tokens: commands.TokensResult{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
		UserID: userID,
	}, nil
}
