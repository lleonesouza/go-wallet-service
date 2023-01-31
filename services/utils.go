package services

import (
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseStringToInt(s string) (int, error) {
	v, err := strconv.ParseInt(s[0:], 10, 64)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func FormatError(err error) map[string]string {
	m := make(map[string]string)
	m["error"] = err.Error()
	return m
}
