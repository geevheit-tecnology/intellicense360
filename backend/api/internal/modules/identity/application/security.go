package application

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
)

func newID(prefix string) string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return prefix + "_" + hex.EncodeToString(b[:])
}

func hashValue(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func hashPassword(password string) string {
	salt := newID("salt")
	sum := sha256.Sum256([]byte(salt + ":" + password))
	return salt + ":" + hex.EncodeToString(sum[:])
}

func verifyPassword(hash string, password string) bool {
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
