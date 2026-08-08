package application

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/ports"
)

type TokenService struct {
	issuer string
	secret []byte
}

func NewTokenService(issuer string, secret string) TokenService {
	if secret == "" {
		secret = "development-identity-secret"
	}
	return TokenService{issuer: issuer, secret: []byte(secret)}
}

func (s TokenService) Issue(_ context.Context, claims ports.AuthClaims) (string, error) {
	claims.Issuer = s.issuer
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	unsigned := encodeSegment(headerJSON) + "." + encodeSegment(claimsJSON)
	return unsigned + "." + s.sign(unsigned), nil
}

func (s TokenService) Validate(_ context.Context, token string) (ports.AuthClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ports.AuthClaims{}, ErrInvalidToken
	}
	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(s.sign(unsigned))) {
		return ports.AuthClaims{}, ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return ports.AuthClaims{}, ErrInvalidToken
	}
	var claims ports.AuthClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ports.AuthClaims{}, ErrInvalidToken
	}
	if claims.Issuer != s.issuer || claims.ExpiresAt <= time.Now().Unix() {
		return ports.AuthClaims{}, ErrInvalidToken
	}
	return claims, nil
}

func (s TokenService) sign(unsigned string) string {
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(unsigned))
	return encodeSegment(mac.Sum(nil))
}

func encodeSegment(value []byte) string {
	return base64.RawURLEncoding.EncodeToString(value)
}
