package database

import (
	"fmt"
	"log"
	"userManagement/internal/config"
	"userManagement/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(env config.Env){
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", env.DB_USER, env.DB_PASSWORD, env.DB_HOST, env.DB_PORT, env.DB_NAME)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection falied due to %v", err)
	}

	err = Automigrate(db)
	if err != nil {
		log.Fatalf("Database Automigration failed due to %v", err)

	}

}

func Automigrate(db *gorm.DB) error{
	return db.AutoMigrate(
		&models.User{},
	)
}
