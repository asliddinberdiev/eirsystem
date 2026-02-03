// Package hasher provides functions for hashing and verifying data.
package hasher

import "golang.org/x/crypto/bcrypt"

func Hash(data string) (string, error) {
	hashedData, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedData), nil
}

func Verify(data string, hashedData string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedData), []byte(data))
}
