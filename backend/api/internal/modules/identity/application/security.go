package application

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
)

func newID(_ string) string { return uuid.NewString() }

func hashValue(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

// HashPassword is shared by controlled provisioning and user creation.
func HashPassword(password string) (string, error) {
	if err := validatePassword(domain.PasswordPolicy{MinLength: 8, RequireNumber: true}, password); err != nil {
		return "", err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
func verifyPassword(hash string, password string) bool {
	if strings.HasPrefix(hash, "$2") {
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
	}
	// Compatibility for the explicitly selected in-memory development adapter.
	parts := strings.Split(hash, ":")
	if len(parts) != 2 {
		return false
	}
	sum := sha256.Sum256([]byte(parts[0] + ":" + password))
	return hmac.Equal([]byte(parts[1]), []byte(hex.EncodeToString(sum[:])))
}

func validatePassword(policy domain.PasswordPolicy, password string) error {
	if len(password) < policy.MinLength {
		return ErrPasswordPolicy
	}
	if policy.RequireNumber && !strings.ContainsAny(password, "0123456789") {
		return ErrPasswordPolicy
	}
	if policy.RequireSpecial && !strings.ContainsAny(password, "!@#$%^&*()-_=+[]{};:,.?/") {
		return ErrPasswordPolicy
	}
	return nil
}

func randomToken() string {
	var b [32]byte
	_, _ = rand.Read(b[:])
	return base64.RawURLEncoding.EncodeToString(b[:])
}

func unixSeconds(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}
