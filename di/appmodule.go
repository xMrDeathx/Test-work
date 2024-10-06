package di

import (
	"TestWork/authentication/cmd"
	"fmt"
	"log"
	"net/http"
	"os"
)

func InitAppModule() {
	authConn, err := cmd.InitAuthModule(cmd.NewConfig())
	if authConn == nil {
		fmt.Printf("Can't connect to company database")
		return
	}
	defer authConn.Close()

	appPort := os.Getenv("BACKEND_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	err = http.ListenAndServe(":"+appPort, nil)
	if err != nil {
		log.Panic("ListenAndServe: ", err)
	}
}

func Migrate() {
	err := cmd.Migrate(cmd.NewConfig())
	if err != nil {
		log.Fatal("failed migrate auth module:", err)
	}
}
