package config

import (
	"fmt"
	"log"
	"time"

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

	var db *gorm.DB
	var err error
	for i := 1; i <= 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(lvl)})
		if err == nil {
			break
		}
		log.Printf("[db] tentativa %d/5 falhou: %v — aguardando 2s...", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("[db] não foi possível conectar após 5 tentativas: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[db] obter sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(App.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(App.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(App.DBConnMaxLifetimeMin) * time.Minute)

	log.Println("[db] conectado com sucesso!")
	DB = db
}
