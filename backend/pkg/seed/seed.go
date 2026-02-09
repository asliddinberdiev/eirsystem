// Package seed - Seed layer
package seed

import (
	"fmt"

	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedSystemAdmin(log logger.Logger, db *gorm.DB, cfg config.SeedSystemAdmin) error {
	log.Info("Seeding System Admin...", logger.Any("config", cfg))
	if cfg.Username == "" || cfg.Password == "" {
		log.Warn("System Admin username/password is empty. Seeding skipped.")
		return nil
	}

	var exists bool

	checkSQL := `SELECT EXISTS(SELECT 1 FROM users WHERE role = 'system')`

	if err := db.Raw(checkSQL).Scan(&exists).Error; err != nil {
		return fmt.Errorf("error checking admin: %w", err)
	}

	if exists {
		log.Info("System Admin already exists. Skipping seeding.")
		return nil
	}

	log.Info("System Admin not found. Creating...")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	insertSQL := `
		INSERT INTO users (tenant_id, full_name, username, password_hash, role, phone, is_active)
		VALUES (NULL, ?, ?, ?, 'system', ?, TRUE)
	`

	if err := db.Exec(insertSQL, cfg.FullName, cfg.Username, string(hashedPassword), cfg.Phone).Error; err != nil {
		return fmt.Errorf("error saving admin: %w", err)
	}

	log.Info("System Admin successfully created!", logger.String("username", cfg.Username))
	return nil
}
