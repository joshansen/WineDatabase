package main

import (
	"github.com/joshansen/WineDatabase/utils"
	"github.com/joshansen/WineDatabase/web"
	"github.com/stretchr/graceful"
	"os"
)

func main() {
	isDevelopment := os.Getenv("ENVIRONMENT") == "development"
	dbURL := os.Getenv("DB_PORT_27017_TCP_ADDR")

	dbAccessor := utils.NewDatabaseAccessor(dbURL, os.Getenv("DATABASE_NAME"), 0)
	s := web.NewServer(*dbAccessor, os.Getenv("SESSION_SECRET"), isDevelopment)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	graceful.Run(":"+port, 0, s)
}
