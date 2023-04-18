package main

import (
	"fmt"
	"os"

	"crowdfund/pkg/database"
	"crowdfund/pkg/helpers"
	"crowdfund/pkg/helpers/logger"
	"crowdfund/pkg/migration"
	"crowdfund/pkg/oauth"
	"crowdfund/pkg/server"

	"github.com/joho/godotenv"
)

func main() {
	cmds := os.Args[1:]
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Can't find env file")
	}

	database.Init()
	if helpers.StringInArr("--migrate", cmds) {
		migration.Run()
	}
	oauth.Init()
	logger.Init()
	server.Init()
}
