// Package persistence expõe a função Migrate que cria/atualiza as tabelas.
package persistence

import (
	"log"

	"gorm.io/gorm"
)

// Migrate executa o AutoMigrate de todos os modelos GORM.
func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(
		&gormUser{},
		&gormLab{},
		&gormComputer{},
		&gormAlert{},
	); err != nil {
		log.Fatalf("[migrate] falha: %v", err)
	}
	log.Println("[migrate] tabelas sincronizadas")
}
