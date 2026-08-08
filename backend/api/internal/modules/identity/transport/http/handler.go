package http

import (
	"net/http"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/identity/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	auth        ports.AuthenticationService
	users       ports.UserService
	roles       ports.RoleService
	permissions ports.PermissionService
	tenants     ports.TenantService
	audit       ports.AuditService
}

func NewHandler(auth ports.AuthenticationService, users ports.UserService, roles ports.RoleService, permissions ports.PermissionService, tenants ports.TenantService, audit ports.AuditService) Handler {
	return Handler{auth: auth, users: users, roles: roles, permissions: permissions, tenants: tenants, audit: audit}
}

func (h Handler) RegisterPublicRoutes(router gin.IRoutes) {
	router.POST("/login", h.login)
	router.POST("/refresh", h.refresh)
}

func (h Handler) RegisterProtectedRoutes(router *gin.RouterGroup) {
	router.POST("/logout", h.logout)
	router.GET("/me", h.me)

	users := router.Group("/users", middleware.RequirePermission("identity.users.manage"))
	users.POST("", h.createUser)
	users.GET("", h.listUsers)
	users.GET("/:id", h.getUser)
	users.PUT("/:id", h.updateUser)
	users.DELETE("/:id", h.deleteUser)

	roles := router.Group("/roles", middleware.RequirePermission("identity.roles.manage"))
	roles.POST("", h.createRole)
	roles.GET("", h.listRoles)
	roles.PUT("/:id", h.updateRole)
	roles.DELETE("/:id", h.deleteRole)

	permissions := router.Group("/permissions", middleware.RequirePermission("identity.permissions.manage"))
	permissions.POST("", h.createPermission)
	permissions.GET("", h.listPermissions)
	permissions.PUT("/:key", h.updatePermission)
	permissions.DELETE("/:key", h.deletePermission)

	tenants := router.Group("/tenants", middleware.RequirePermission("identity.tenants.manage"))
	tenants.POST("", h.createTenant)
	tenants.GET("", h.listTenants)
	tenants.PUT("/:id", h.updateTenant)
	tenants.DELETE("/:id", h.deleteTenant)

	router.GET("/audit-logs", h.listAudit)
}

func (h Handler) login(c *gin.Context) {
	var req struct {
		TenantID string `json:"tenant_id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if !bind(c, &req) {
		return
	}
	if req.TenantID == "" {
		req.TenantID = c.GetHeader("X-Tenant-ID")
	}
	if req.TenantID == "" {
		req.TenantID = "bootstrap-tenant"
	}
	result, err := h.auth.Login(c.Request.Context(), req.TenantID, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h Handler) refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if !bind(c, &req) {
		return
	}
	result, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h Handler) logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if !bind(c, &req) {
		return
	}
	if err := h.auth.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) me(c *gin.Context) {
	user, err := h.auth.Me(c.Request.Context(), tenantID(c), actorID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h Handler) createUser(c *gin.Context) {
	var req struct{ Name, Email, Password string }
	if !bind(c, &req) {
		return
	}
	user, err := h.users.Create(c.Request.Context(), domain.User{TenantID: tenantID(c), Name: req.Name, Email: req.Email}, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h Handler) listUsers(c *gin.Context) {
	users, _ := h.users.List(c.Request.Context(), tenantID(c))
	c.JSON(http.StatusOK, users)
}
func (h Handler) getUser(c *gin.Context) {
	user, err := h.users.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
func (h Handler) updateUser(c *gin.Context) {
	var user domain.User
	if !bind(c, &user) {
		return
	}
	user.ID = c.Param("id")
	user.TenantID = tenantID(c)
	saved, err := h.users.Update(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, saved)
}
func (h Handler) deleteUser(c *gin.Context) {
	_ = h.users.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createRole(c *gin.Context) {
	var role domain.Role
	if !bind(c, &role) {
		return
	}
	role.TenantID = tenantID(c)
	saved, err := h.roles.Create(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, saved)
}
func (h Handler) listRoles(c *gin.Context) {
	roles, _ := h.roles.List(c.Request.Context(), tenantID(c))
	c.JSON(http.StatusOK, roles)
}
func (h Handler) updateRole(c *gin.Context) {
	var role domain.Role
	if !bind(c, &role) {
		return
	}
	role.ID = c.Param("id")
	role.TenantID = tenantID(c)
	saved, err := h.roles.Update(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, saved)
}
func (h Handler) deleteRole(c *gin.Context) {
	_ = h.roles.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createPermission(c *gin.Context) {
	var permission domain.Permission
	if !bind(c, &permission) {
		return
	}
	saved, err := h.permissions.Create(c.Request.Context(), permission)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, saved)
}
func (h Handler) listPermissions(c *gin.Context) {
	permissions, _ := h.permissions.List(c.Request.Context())
	c.JSON(http.StatusOK, permissions)
}
func (h Handler) updatePermission(c *gin.Context) {
	var permission domain.Permission
	if !bind(c, &permission) {
		return
	}
	permission.Key = c.Param("key")
	saved, err := h.permissions.Update(c.Request.Context(), permission)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, saved)
}
func (h Handler) deletePermission(c *gin.Context) {
	_ = h.permissions.Delete(c.Request.Context(), c.Param("key"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createTenant(c *gin.Context) {
	var tenant domain.Tenant
	if !bind(c, &tenant) {
		return
	}
	saved, err := h.tenants.Create(c.Request.Context(), tenant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, saved)
}
func (h Handler) listTenants(c *gin.Context) {
	tenants, _ := h.tenants.List(c.Request.Context())
	c.JSON(http.StatusOK, tenants)
}
func (h Handler) updateTenant(c *gin.Context) {
	var tenant domain.Tenant
	if !bind(c, &tenant) {
		return
	}
	tenant.ID = c.Param("id")
	saved, err := h.tenants.Update(c.Request.Context(), tenant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, saved)
}
func (h Handler) deleteTenant(c *gin.Context) {
	_ = h.tenants.Delete(c.Request.Context(), c.Param("id"))
	c.Status(http.StatusNoContent)
}
func (h Handler) listAudit(c *gin.Context) {
	logs, _ := h.audit.List(c.Request.Context(), tenantID(c))
	c.JSON(http.StatusOK, logs)
}

func bind(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return false
	}
	return true
}

func tenantID(c *gin.Context) string {
	value, _ := c.Get(string(contextkeys.TenantID))
	tenantID, _ := value.(string)
	if tenantID == "" {
		tenantID = "bootstrap-tenant"
	}
	return tenantID
}

func actorID(c *gin.Context) string {
	value, _ := c.Get(string(contextkeys.ActorID))
	actorID, _ := value.(string)
	return actorID
}
