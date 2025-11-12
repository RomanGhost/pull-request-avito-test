package main

import (
	"log"
	"os"

	"pr-app/internal/db"
	"pr-app/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://pruser:prpassword@db:5432/prdb?sslmode=disable"
	}

	database := db.InitDB(dsn)

	r := gin.Default()
	router.SetupRoutes(r, database)

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
