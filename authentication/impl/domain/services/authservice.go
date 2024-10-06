package services

import (
	"TestWork/authentication/impl/app/commands"
	"context"
)

type AuthService interface {
	Login(context context.Context, loginData commands.LoginCommand, requestIP string) (commands.LoginResult, error)
}
