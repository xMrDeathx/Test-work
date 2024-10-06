package transport

import (
	frontendapi "TestWork/authentication/api/frontend"
	"TestWork/authentication/impl/app/commands"
	"TestWork/authentication/impl/domain/services"
	"encoding/json"
	"errors"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"io"
	"net/http"
	"time"
)

func NewAuthServer(
	authService services.AuthService,
	tokenService services.TokenService,
) frontendapi.ServerInterface {
	return &authServer{
		authService:  authService,
		tokenService: tokenService,
	}
}

type authServer struct {
	authService  services.AuthService
	tokenService services.TokenService
}

// Запрос на логинацию
func (server *authServer) Login(w http.ResponseWriter, r *http.Request) {
	requestIP := r.RemoteAddr
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var loginData commands.LoginCommand
	err = json.Unmarshal(requestBody, &loginData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loginResponse, err := server.authService.Login(r.Context(), loginData, requestIP)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	server.setRefreshTokenToCookie(w, loginResponse.Tokens.RefreshToken)

	response, err := json.Marshal(frontendapi.LoginResponse{
		AccessToken: loginResponse.Tokens.AccessToken,
		UserId:      loginResponse.UserID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Запрос на обновление токенов
func (server *authServer) RefreshToken(w http.ResponseWriter, r *http.Request, userId openapi_types.UUID) {
	requestIP := r.RemoteAddr
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokens, err := server.tokenService.RefreshTokens(r.Context(), userId, cookie.Value, requestIP)
	if err == errors.New("user with token not found") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	server.setRefreshTokenToCookie(w, tokens.RefreshToken)

	response, err := json.Marshal(frontendapi.TokenResponse{
		AccessToken: tokens.AccessToken,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Вставка refresh token в куки, возвращаемую по запросу
func (server *authServer) setRefreshTokenToCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(14 * 24 * time.Hour), // 2 weeks
	}
	http.SetCookie(w, &cookie)
}
