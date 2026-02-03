// Package dto provides data transfer objects for the application.
package dto

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}