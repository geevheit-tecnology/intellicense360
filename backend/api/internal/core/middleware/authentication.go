package middleware

import (
	"net/http"
	"strings"

	identityports "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

func Authentication(auth identityports.AuthenticationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		claims, err := auth.ValidateAccessToken(c.Request.Context(), strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid bearer token"})
			return
		}

		c.Set(string(contextkeys.ActorID), claims.UserID)
		c.Set(string(contextkeys.TenantID), claims.TenantID)
		c.Set(string(contextkeys.SessionID), claims.SessionID)
		c.Set(string(contextkeys.Permissions), claims.Permissions)
		c.Request = c.Request.WithContext(contextkeys.WithTenantID(c.Request.Context(), claims.TenantID))
		c.Request = c.Request.WithContext(contextkeys.WithActorID(c.Request.Context(), claims.UserID))
		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get(string(contextkeys.Permissions))
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}
		permissions, _ := raw.([]string)
		for _, item := range permissions {
			if item == permission {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
	}
}
