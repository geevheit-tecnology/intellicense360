package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/transport/dto"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/transport/mapper"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	suppliers  ports.SupplierService
	categories ports.CategoryService
	types      ports.TypeService
	contacts   ports.ContactService
	addresses  ports.AddressService
	documents  ports.DocumentService
	contracts  ports.ContractService
	ratings    ports.RatingService
}

func NewHandler(suppliers ports.SupplierService, categories ports.CategoryService, types ports.TypeService, contacts ports.ContactService, addresses ports.AddressService, documents ports.DocumentService, contracts ports.ContractService, ratings ports.RatingService) Handler {
	return Handler{suppliers: suppliers, categories: categories, types: types, contacts: contacts, addresses: addresses, documents: documents, contracts: contracts, ratings: ratings}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/suppliers")
	group.Use(middleware.RequirePermission("suppliers.suppliers.manage"))
	group.GET("", h.searchSuppliers)
	group.POST("", h.createSupplier)
	group.GET("/:id", h.getSupplier)
	group.PUT("/:id", h.updateSupplier)
	group.DELETE("/:id", h.deleteSupplier)
	group.GET("/categories", h.searchCategories)
	group.POST("/categories", h.createCategory)
	group.PUT("/categories/:id", h.updateCategory)
	group.DELETE("/categories/:id", h.deleteCategory)
	group.GET("/types", h.searchTypes)
	group.POST("/types", h.createType)
	group.PUT("/types/:id", h.updateType)
	group.DELETE("/types/:id", h.deleteType)
	group.GET("/contacts", h.searchContacts)
	group.POST("/contacts", h.createContact)
	group.PUT("/contacts/:id", h.updateContact)
	group.DELETE("/contacts/:id", h.deleteContact)
	group.GET("/addresses", h.searchAddresses)
	group.POST("/addresses", h.createAddress)
	group.PUT("/addresses/:id", h.updateAddress)
	group.DELETE("/addresses/:id", h.deleteAddress)
	group.GET("/documents", h.searchDocuments)
	group.POST("/documents", h.createDocument)
	group.PUT("/documents/:id", h.updateDocument)
	group.DELETE("/documents/:id", h.deleteDocument)
	group.GET("/contracts", h.searchContracts)
	group.POST("/contracts", h.createContract)
	group.PUT("/contracts/:id", h.updateContract)
	group.DELETE("/contracts/:id", h.deleteContract)
	group.GET("/ratings", h.searchRatings)
	group.POST("/ratings", h.createRating)
	group.PUT("/ratings/:id", h.updateRating)
	group.DELETE("/ratings/:id", h.deleteRating)
}

func (h Handler) createSupplier(c *gin.Context) {
	var req dto.SupplierRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.SupplierFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.suppliers.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchSuppliers(c *gin.Context) {
	page, err := h.suppliers.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getSupplier(c *gin.Context) {
	item, err := h.suppliers.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) updateSupplier(c *gin.Context) {
	var req dto.SupplierRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.SupplierFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.suppliers.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteSupplier(c *gin.Context) {
	err := h.suppliers.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) createCategory(c *gin.Context) {
	var req dto.CategoryRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.CategoryFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.categories.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchCategories(c *gin.Context) {
	page, err := h.categories.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateCategory(c *gin.Context) {
	var req dto.CategoryRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.CategoryFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.categories.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteCategory(c *gin.Context) {
	_ = h.categories.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createType(c *gin.Context) {
	var req dto.TypeRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TypeFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.types.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchTypes(c *gin.Context) {
	page, err := h.types.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateType(c *gin.Context) {
	var req dto.TypeRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TypeFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.types.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteType(c *gin.Context) {
	_ = h.types.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createContact(c *gin.Context) {
	var req dto.ContactRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ContactFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.contacts.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchContacts(c *gin.Context) {
	page, err := h.contacts.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateContact(c *gin.Context) {
	var req dto.ContactRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ContactFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.contacts.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteContact(c *gin.Context) {
	_ = h.contacts.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createAddress(c *gin.Context) {
	var req dto.AddressRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.AddressFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.addresses.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchAddresses(c *gin.Context) {
	page, err := h.addresses.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateAddress(c *gin.Context) {
	var req dto.AddressRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.AddressFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.addresses.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteAddress(c *gin.Context) {
	_ = h.addresses.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createDocument(c *gin.Context) {
	var req dto.DocumentRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.DocumentFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.documents.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchDocuments(c *gin.Context) {
	page, err := h.documents.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateDocument(c *gin.Context) {
	var req dto.DocumentRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.DocumentFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.documents.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteDocument(c *gin.Context) {
	_ = h.documents.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createContract(c *gin.Context) {
	var req dto.ContractRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ContractFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.contracts.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchContracts(c *gin.Context) {
	page, err := h.contracts.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateContract(c *gin.Context) {
	var req dto.ContractRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ContractFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.contracts.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteContract(c *gin.Context) {
	_ = h.contracts.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createRating(c *gin.Context) {
	var req dto.RatingRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.RatingFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.ratings.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchRatings(c *gin.Context) {
	page, err := h.ratings.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateRating(c *gin.Context) {
	var req dto.RatingRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.RatingFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.ratings.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteRating(c *gin.Context) {
	_ = h.ratings.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return ports.Query{Search: c.Query("search"), Page: page, PageSize: pageSize, SortBy: c.Query("sort_by"), SortOrder: c.Query("sort_order"), Filters: map[string]string{"status": c.Query("status"), "type": c.Query("type"), "category_id": c.Query("category_id"), "supplier_id": c.Query("supplier_id")}}
}
func bind(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return false
	}
	return true
}
func respond(c *gin.Context, status int, payload any, err error) {
	if err != nil {
		code := http.StatusBadRequest
		if errors.Is(err, application.ErrNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, payload)
}
func tenantID(c *gin.Context) string {
	value, _ := c.Get(string(contextkeys.TenantID))
	tenant, _ := value.(string)
	return tenant
}
