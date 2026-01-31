// Package domain - Domain layer
package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	FullName     string    `json:"full_name"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}
