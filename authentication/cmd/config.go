package cmd

import "os"

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

func NewConfig() Config {
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "741953"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "test"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	return Config{
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBHost:     dbHost,
		DBPort:     dbPort,
	}
}
