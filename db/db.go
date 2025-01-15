package main

import (
	"fmt"
	"log"

	"github.com/abdisetiakawan/go-ecommerce/internal/config"
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	viperConfig := config.NewViper()

	cfg, err := config.LoadConfig(viperConfig)
	if err != nil {
		log.Fatalf("Failed to load config.json : %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.Charset,
	)


	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database : %v", err)
	}

	if err := db.AutoMigrate(&entity.User{}, &entity.Profile{}, &entity.Store{}); err != nil {
		log.Fatalf("failed to migrate database : %v", err)
	}

	fmt.Println("Migration completed successfully")
}
