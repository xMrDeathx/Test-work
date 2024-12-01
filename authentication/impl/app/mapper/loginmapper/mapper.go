package loginmapper

import (
	"TestWork/authentication/impl/app/commands"
	"TestWork/authentication/impl/domain/model"
)

func NewLoginDataToDomainLoginData(command commands.LoginCommand) model.LoginData {
	return model.LoginData{
		Email:    command.Email,
		Password: command.Password,
	}
}
