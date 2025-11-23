package utils

import (
	"encoding/hex"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateID() string {
	b := make([]byte, 8)
	rand.Seed(time.Now().UnixNano())
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateOrderNumber() string {
	rand.Seed(time.Now().UnixNano())
	timestamp := time.Now().Format("20060102")
	randomPart := rand.Intn(9999) + 1000
	return "ORD-" + timestamp + "-" + strings.ToUpper(hex.EncodeToString([]byte{byte(randomPart)}))
}

func GeneratePaymentCode() string {
	rand.Seed(time.Now().UnixNano())
	timestamp := time.Now().Format("150405")
	randomPart := rand.Intn(99999999)
	return "PAY" + timestamp + hex.EncodeToString([]byte{byte(randomPart)})
}

// Role permission checking functions
func HasPermission(permissions []string, requiredPermission string) bool {
	for _, p := range permissions {
		if p == requiredPermission {
			return true
		}
	}
	return false
}

func HasAnyPermission(permissions []string, requiredPermissions []string) bool {
	for _, required := range requiredPermissions {
		if HasPermission(permissions, required) {
			return true
		}
	}
	return false
}

func HasAllPermissions(permissions []string, requiredPermissions []string) bool {
	for _, required := range requiredPermissions {
		if !HasPermission(permissions, required) {
			return false
		}
	}
	return true
}
