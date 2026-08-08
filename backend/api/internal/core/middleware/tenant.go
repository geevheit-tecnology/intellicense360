package middleware

import (
	coreerrors "github.com/geevheit/intelligence360/backend/api/internal/core/errors"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

const tenantHeader = "X-Tenant-ID"

func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader(tenantHeader)
		if current, ok := c.Get(string(contextkeys.TenantID)); ok {
			if value, ok := current.(string); ok && value != "" {
				tenantID = value
			}
		}
		if tenantID == "" {
			tenantID = "bootstrap-tenant"
		}

		if tenantID == "" {
			_ = c.Error(coreerrors.TenantMissing())
			c.Abort()
			return
		}

		c.Set(string(contextkeys.TenantID), tenantID)
		c.Request = c.Request.WithContext(contextkeys.WithTenantID(c.Request.Context(), tenantID))
		c.Next()
	}
}
