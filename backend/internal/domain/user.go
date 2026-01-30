// Package domain - Domain layer
package domain

import "time"

type User struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	Username  string    `json:"username" gorm:"column:username"`
	Email     string    `json:"email" gorm:"column:email"`
	Password  string    `json:"-" gorm:"column:password_hash"`
	FullName  string    `json:"full_name" gorm:"column:full_name"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}
