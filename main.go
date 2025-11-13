package main

import (
	"log"
	"os"

	"github.com/RomanGhost/pull-request-avito-test/internal/database"
	"github.com/RomanGhost/pull-request-avito-test/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://pruser:prpassword@db:5432/prdb?sslmode=disable"
	}

	database := database.InitDB(dsn)

	r := gin.Default()

	h := handler.RegisterHandlers(database)
	h.RegisterRoutes(r)

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
