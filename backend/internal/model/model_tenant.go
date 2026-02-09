package model

import (
	"time"
)

type Tenant struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Slug                string    `json:"slug"`
	OwnerID             string    `json:"owner_id"`
	IsActive            bool      `json:"is_active"`
	SubscriptionEndDate time.Time `json:"subscription_end_date"`
	CreatedAt           time.Time `json:"created_at"`
}