package sql

import (
	"TestWork/authentication/impl/domain/repositories"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewAuthRepository(conn *pgxpool.Pool) repositories.AuthRepository {
	return &authRepository{conn: conn}
}

type authRepository struct {
	conn *pgxpool.Pool
}

func (repo *authRepository) GetPassword(context context.Context, email string) (string, error) {
	var password string

	err := repo.conn.QueryRow(context, `
		SELECT password FROM auth_user 
        WHERE email=$1
	`, email).Scan(&password)

	if err == pgx.ErrNoRows {
		return "", errors.New("user not found")
	} else if err != nil {
		return "", err
	}

	return password, nil
}

func (repo *authRepository) GetUserByEmail(context context.Context, email string) (uuid.UUID, error) {
	var userID uuid.UUID

	err := repo.conn.QueryRow(context, `
		SELECT id FROM auth_user 
        WHERE email=$1
	`, email).Scan(&userID)

	if err == pgx.ErrNoRows {
		return uuid.UUID{}, errors.New("user not found")
	} else if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}
