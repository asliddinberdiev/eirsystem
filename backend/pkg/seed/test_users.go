package seed

import (
	"fmt"
	"time"

	"github.com/asliddinberdiev/eirsystem/internal/model"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	casbinlib "github.com/casbin/casbin/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedTestUsers(log logger.Logger, db *gorm.DB, enforcer *casbinlib.Enforcer) error {
	log.Info("Seeding Test Users (Owner, Doctor, Nurse)...")
	
	tenantID := uuid.New().String()
	testTenant := &model.Tenant{
		ID:                  tenantID,
		Name:                "Test Tenant",
		Slug:                "test-tenant",
		OwnerID:             uuid.New().String(),
		IsActive:            true,
		SubscriptionEndDate: time.Now().AddDate(1, 0, 0),
	}

	if err := db.FirstOrCreate(testTenant, model.Tenant{Slug: "test-tenant"}).Error; err != nil {
		return fmt.Errorf("error seeding tenant: %w", err)
	}
	
	tenantID = testTenant.ID

	users := []struct {
		Username string
		Password string
		Role     string
		FullName string
		ClinicID string
	}{
		{"system", "password", "system", "Test System", tenantID},
		{"owner", "password", "owner", "Test Owner", tenantID},
		{"admin", "password", "admin", "Test Admin", tenantID},
		{"doctor", "password", "doctor", "Test Doctor", tenantID},
		{"nurse", "password", "nurse", "Test Nurse", tenantID},
		{"technician", "password", "technician", "Test Technician", tenantID},
		{"reception", "password", "reception", "Test Reception", tenantID},
	}

	for _, u := range users {
		var exists bool
		checkSQL := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`
		if err := db.Raw(checkSQL, u.Username).Scan(&exists).Error; err != nil {
			return fmt.Errorf("error checking user %s: %w", u.Username, err)
		}

		if exists {
			log.Info("User already exists. Skipping.", logger.String("username", u.Username))
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error hashing password for %s: %w", u.Role, err)
		}

		insertSQL := `
            INSERT INTO users (tenant_id, full_name, username, password_hash, role, phone, is_active)
            VALUES (?, ?, ?, ?, ?, '123456789', TRUE) RETURNING id
        `
		var userID string
		if err := db.Raw(insertSQL, u.ClinicID, u.FullName, u.Username, string(hashedPassword), u.Role).Scan(&userID).Error; err != nil {
			return fmt.Errorf("error creating user %s: %w", u.Role, err)
		}

		if _, err := enforcer.AddGroupingPolicy(userID, u.Role, u.ClinicID); err != nil {
			return fmt.Errorf("error adding grouping policy for %s: %w", u.Role, err)
		}

		var policies [][]string

		switch u.Role {
		case "owner":
			policies = append(policies, []string{"owner", u.ClinicID, "/api/v1/test/owner", "GET"})
			policies = append(policies, []string{"owner", u.ClinicID, "/api/v1/test/doctor", "GET"})
			policies = append(policies, []string{"owner", u.ClinicID, "/api/v1/test/nurse", "GET"})
		case "doctor":
			policies = append(policies, []string{"doctor", u.ClinicID, "/api/v1/test/doctor", "GET"})
		case "nurse":
			policies = append(policies, []string{"nurse", u.ClinicID, "/api/v1/test/nurse", "GET"})
		}

		for _, policy := range policies {
			if _, err := enforcer.AddPolicy(policy[0], policy[1], policy[2], policy[3]); err != nil {
				log.Error("Error adding permission policy", logger.Error(err))
				return err
			}
		}

		log.Info("User and policies successfully created!", logger.String("role", u.Role))
	}

	if err := enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save policies: %w", err)
	}

	return nil
}