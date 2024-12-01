package main

import (
	"TestWork/di"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Panic("Ошибка при загрузке .env файла", err)
	}

	di.Migrate()
	di.InitAppModule()
}
