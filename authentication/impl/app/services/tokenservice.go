package services

import (
	"TestWork/authentication/impl/app/commands"
	"TestWork/authentication/impl/app/mapper/tokenmapper"
	"TestWork/authentication/impl/domain/repositories"
	"TestWork/authentication/impl/domain/services"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/smtp"
	"time"
)

var ErrUserWithTokenNotFound = errors.New("user with token not found")

const SecretKey = "secret"

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
	hashedRefreshToken, err := hashRefreshToken(refreshToken)
	if err != nil {
		return commands.TokensResult{}, err
	}
	//Создание и сохранение новой сессии с refresh token в базу данных
	session := tokenmapper.NewDomainSession(userID, requestIP, hashedRefreshToken)
	err = service.repository.SaveToken(context, session)
	if err != nil {
		return commands.TokensResult{}, err
	}

	return tokenmapper.NewTokensResultFromEntity(accessToken, encodeRefreshToken(refreshToken)), nil
}

// Обработка запроса на обновление токенов
func (service *tokenService) RefreshTokens(context context.Context, userID uuid.UUID, token string, requestIP string) (commands.TokensResult, error) {
	requestToken, err := decodeRefreshToken(token)
	if err != nil {
		return commands.TokensResult{}, err
	}

	domainToken, userIP, err := service.repository.GetToken(context, userID)
	if errors.Is(err, ErrUserWithTokenNotFound) {
		return commands.TokensResult{}, ErrUserWithTokenNotFound
	}
	if err != nil {
		return commands.TokensResult{}, err
	}

	if err = validateToken(requestToken, domainToken.Token); err != nil {
		return commands.TokensResult{}, errors.New("invalid refresh token")
	}

	if err = checkExpiration(domainToken.ExpiresIn); err != nil {
		return commands.TokensResult{}, err
	}

	if requestIP != userIP {
		if msgErr := sendWarning(); msgErr != nil {
			log.Printf("Error while sending warning message: %s", msgErr)
		}
	}

	newAccessToken, err := createAccessToken(userID, requestIP)
	if err != nil {
		return commands.TokensResult{}, err
	}

	newRefreshToken := uuid.New()
	hashedRefreshToken, err := hashRefreshToken(newRefreshToken)
	if err != nil {
		return commands.TokensResult{}, err
	}

	newDomainRefreshToken := tokenmapper.NewDomainRefreshToken(hashedRefreshToken)
	err = service.repository.UpdateToken(context, domainToken.Token, newDomainRefreshToken, requestIP)
	if err != nil {
		return commands.TokensResult{}, err
	}

	return tokenmapper.NewTokensResultFromEntity(newAccessToken, encodeRefreshToken(newRefreshToken)), nil
}

func createAccessToken(userID uuid.UUID, requestIP string) (string, error) {
	//Создание payload для access token
	accessTokenPayload := jwt.MapClaims{
		"exp":       time.Now().Add(30 * time.Minute).UnixMilli(),
		"user":      userID,
		"requestIP": requestIP,
	}
	//Создание access token
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenPayload).SignedString([]byte(SecretKey))
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

func hashRefreshToken(token uuid.UUID) ([]byte, error) {
	byteToken, err := token.MarshalBinary()
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword(byteToken, 12)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func validateToken(oldToken uuid.UUID, hashedOldToken []byte) error {
	byteOldToken, err := oldToken.MarshalBinary()
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword(hashedOldToken, byteOldToken)
}

func checkExpiration(expirationTime int64) error {
	if time.Now().UnixMilli() > expirationTime {
		return errors.New("token had expired")
	}

	return nil
}

func sendWarning() error {
	username := "b6496edc9f8b17"
	password := "9d02e85e5638d0"
	host := "sandbox.smtp.mailtrap.io"
	port := "25"

	// Subject and body
	subject := "WARNING"
	body := "Entering from another device!"

	// Sender and receiver
	from := "from@example.com"
	to := []string{
		"to@example.com",
	}

	// Build the message
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += fmt.Sprintf("\r\n%s\r\n", body)

	// Authentication.
	auth := smtp.PlainAuth("", username, password, host)

	// Send email
	err := smtp.SendMail(host+":"+port, auth, from, to, []byte(message))
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Email sent successfully.")
	return nil
}
