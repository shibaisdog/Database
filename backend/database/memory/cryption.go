package memory

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Hash(value string) (string, error) {
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedValue), nil
}

func Compare(hashedValue, value string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
