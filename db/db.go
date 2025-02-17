package main

import (
	"fmt"
	"log"

	"github.com/abdisetiakawan/go-ecommerce/internal/config"
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	evententity "github.com/abdisetiakawan/go-ecommerce/internal/entity/event_entity"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


func Migrate(db *gorm.DB, log *logrus.Logger) {
    if err := db.AutoMigrate(&entity.User{}); err != nil {
        log.Fatalf("failed to migrate User entity: %v", err)
    }

    if err := db.AutoMigrate(&entity.Profile{}); err != nil {
        log.Fatalf("failed to migrate Profile entity: %v", err)
    }

    if err := db.AutoMigrate(&entity.Store{}); err != nil {
        log.Fatalf("failed to migrate Store entity: %v", err)
    }

    if err := db.AutoMigrate(&entity.Product{}); err != nil {
        log.Fatalf("failed to migrate Product entity: %v", err)
    }

    if err := db.AutoMigrate(&entity.Order{}); err != nil {
        log.Fatalf("failed to migrate Order entity: %v", err)
    }

    if err := db.AutoMigrate(&entity.OrderItem{}); err != nil {
        log.Fatalf("failed to migrate OrderItem entity: %v", err)
    }

    if err := db.AutoMigrate(&entity.Payment{}); err != nil {
        log.Fatalf("failed to migrate Payment entity: %v", err)
    }

    if err := db.AutoMigrate(&entity.Shipping{}); err != nil {
        log.Fatalf("failed to migrate Shipping entity: %v", err)
    }

	// Event
	if err := db.AutoMigrate(&evententity.OrderEvent{}); err != nil {
		log.Fatalf("failed to migrate OrderEvent entity: %v", err)
	}

    log.Info("Migration completed successfully")
}

func main() {
	viperConfig := config.NewViper()
	username := viperConfig.GetString("DATABASE_USERNAME")
	password := viperConfig.GetString("DATABASE_PASSWORD")
	host := viperConfig.GetString("DATABASE_HOST")
	port := viperConfig.GetInt("DATABASE_PORT")
	database := viperConfig.GetString("DATABASE_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database : %v", err)
	}

	logger := config.NewLogger(viperConfig)
	Migrate(db, logger)
}
