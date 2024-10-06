package cmd

import (
	"TestWork/authentication/api/frontend"
	"TestWork/authentication/impl/app/services"
	"TestWork/authentication/impl/infrastructure/sql"
	"TestWork/authentication/impl/infrastructure/transport"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	"time"
)

func InitAuthModule(config Config) (
	*pgxpool.Pool,
	error,
) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, 5432, config.DBUser, config.DBPassword, config.DBName)

	conn, _ := ConnectLoop(connStr, 30*time.Second)

	authRepository := sql.NewAuthRepository(conn)
	tokenRepository := sql.NewTokenStorage(conn)
	tokenService := services.NewTokenService(tokenRepository)
	authService := services.NewAuthService(authRepository, tokenService)
	authServer := transport.NewAuthServer(authService, tokenService)

	router := mux.NewRouter()

	options := frontendapi.GorillaServerOptions{
		BaseRouter: router,
		Middlewares: []frontendapi.MiddlewareFunc{func(handler http.Handler) http.Handler {
			return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handler.ServeHTTP(w, r)
			}))
		}},
	}
	r := frontendapi.HandlerWithOptions(authServer, options)
	http.Handle("/authorization/", r)

	return conn, nil
}

func ConnectLoop(connStr string, timeout time.Duration) (*pgxpool.Pool, error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)
	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %s timeout", timeout)

		case <-ticker.C:
			db, err := pgxpool.Connect(context.Background(), connStr)
			if err == nil {
				return db, nil
			}
		}
	}
}
