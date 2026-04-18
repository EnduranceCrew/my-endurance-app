package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=America/Sao_Paulo",
		App.DBHost, App.DBPort, App.DBUser, App.DBPassword, App.DBName, App.DBSSLMode,
	)

	lvl := logger.Info
	if App.GinMode == "release" {
		lvl = logger.Warn
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(lvl),
	})
	if err != nil {
		log.Fatalf("[db] falha ao conectar: %v", err)
	}

	log.Println("[db] conectado com sucesso!")
	DB = db
}
