package database

import (
	"log"

	"github.com/RomanGhost/pull-request-avito-test/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Автоматическая миграция схемы
	if err := db.AutoMigrate(
		&domain.Team{},
		&domain.User{},
		&domain.PullRequest{},
		&domain.PRReviewer{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}
