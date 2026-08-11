package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	transactions ports.TransactionService
	categories   ports.CatalogService[domain.CostCategory]
	types        ports.CatalogService[domain.CostType]
	centers      ports.CatalogService[domain.CostCenter]
	accounts     ports.CatalogService[domain.Account]
	methods      ports.CatalogService[domain.PaymentMethod]
	periods      ports.PeriodService
	budgets      ports.BudgetService
	adjustments  ports.AdjustmentService
}

func NewHandler(transactions ports.TransactionService, categories ports.CatalogService[domain.CostCategory], types ports.CatalogService[domain.CostType], centers ports.CatalogService[domain.CostCenter], accounts ports.CatalogService[domain.Account], methods ports.CatalogService[domain.PaymentMethod], periods ports.PeriodService, budgets ports.BudgetService, adjustments ports.AdjustmentService) Handler {
	return Handler{transactions: transactions, categories: categories, types: types, centers: centers, accounts: accounts, methods: methods, periods: periods, budgets: budgets, adjustments: adjustments}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/financial")
	group.Use(middleware.RequirePermission("financial.financial.manage"))
	group.GET("/transactions", h.searchTransactions)
	group.GET("/transactions/:id", h.getTransaction)
	group.PUT("/transactions/:id", h.updateTransaction)
	group.POST("/transactions/:id/approve", h.approve)
	group.POST("/transactions/:id/pay", h.pay)
	group.POST("/transactions/:id/receive", h.receive)
	group.POST("/transactions/:id/cancel", h.cancel)
	group.POST("/transactions/:id/adjust", h.adjustByID)
	group.GET("/expenses", h.searchExpenses)
	group.POST("/expenses", h.createExpense)
	group.GET("/revenues", h.searchRevenues)
	group.POST("/revenues", h.createRevenue)
	group.GET("/categories", h.searchCategories)
	group.POST("/categories", h.createCategory)
	group.GET("/types", h.searchTypes)
	group.POST("/types", h.createType)
	group.GET("/cost-centers", h.searchCenters)
	group.POST("/cost-centers", h.createCenter)
	group.GET("/accounts", h.searchAccounts)
	group.POST("/accounts", h.createAccount)
	group.GET("/payment-methods", h.searchMethods)
	group.POST("/payment-methods", h.createMethod)
	group.GET("/periods", h.searchPeriods)
	group.POST("/periods", h.createPeriod)
	group.POST("/periods/:id/close", h.closePeriod)
	group.GET("/budgets", h.searchBudgets)
	group.POST("/budgets", h.createBudget)
	group.GET("/adjustments", h.searchAdjustments)
	group.POST("/adjustments", h.adjust)
}

func (h Handler) createExpense(c *gin.Context) {
	var req domain.FinancialTransaction
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.CreatedBy, req.UpdatedBy = tenantID(c), actorID(c), actorID(c)
	saved, err := h.transactions.CreateExpense(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) createRevenue(c *gin.Context) {
	var req domain.FinancialTransaction
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.CreatedBy, req.UpdatedBy = tenantID(c), actorID(c), actorID(c)
	saved, err := h.transactions.CreateRevenue(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateTransaction(c *gin.Context) {
	var req domain.FinancialTransaction
	if !bind(c, &req) {
		return
	}
	req.ID, req.TenantID, req.UpdatedBy = c.Param("id"), tenantID(c), actorID(c)
	saved, err := h.transactions.Update(c.Request.Context(), req)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) approve(c *gin.Context) {
	saved, err := h.transactions.Approve(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) pay(c *gin.Context) {
	saved, err := h.transactions.Pay(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) receive(c *gin.Context) {
	saved, err := h.transactions.Receive(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) cancel(c *gin.Context) {
	saved, err := h.transactions.Cancel(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) adjust(c *gin.Context) {
	var req domain.FinancialAdjustment
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.CreatedBy = tenantID(c), actorID(c)
	saved, err := h.transactions.Adjust(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) adjustByID(c *gin.Context) {
	var req domain.FinancialAdjustment
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.TransactionID, req.CreatedBy = tenantID(c), c.Param("id"), actorID(c)
	saved, err := h.transactions.Adjust(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchTransactions(c *gin.Context) {
	page, err := h.transactions.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) searchExpenses(c *gin.Context) {
	q := query(c)
	q.Filters["kind"] = "expense"
	page, err := h.transactions.Search(c.Request.Context(), tenantID(c), q)
	respond(c, http.StatusOK, page, err)
}
func (h Handler) searchRevenues(c *gin.Context) {
	q := query(c)
	q.Filters["kind"] = "revenue"
	page, err := h.transactions.Search(c.Request.Context(), tenantID(c), q)
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getTransaction(c *gin.Context) {
	item, err := h.transactions.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}

func (h Handler) createCategory(c *gin.Context) {
	var req domain.CostCategory
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.categories.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchCategories(c *gin.Context) {
	page, err := h.categories.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createType(c *gin.Context) {
	var req domain.CostType
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.types.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchTypes(c *gin.Context) {
	page, err := h.types.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createCenter(c *gin.Context) {
	var req domain.CostCenter
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.centers.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchCenters(c *gin.Context) {
	page, err := h.centers.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createAccount(c *gin.Context) {
	var req domain.Account
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.accounts.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchAccounts(c *gin.Context) {
	page, err := h.accounts.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createMethod(c *gin.Context) {
	var req domain.PaymentMethod
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.methods.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchMethods(c *gin.Context) {
	page, err := h.methods.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createPeriod(c *gin.Context) {
	var req periodRequest
	if !bind(c, &req) {
		return
	}
	item := domain.FinancialPeriod{TenantID: tenantID(c), Year: req.Year, Month: req.Month, StartDate: parseDate(req.StartDate), EndDate: parseDate(req.EndDate)}
	saved, err := h.periods.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) closePeriod(c *gin.Context) {
	saved, err := h.periods.Close(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) searchPeriods(c *gin.Context) {
	page, err := h.periods.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createBudget(c *gin.Context) {
	var req domain.Budget
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.budgets.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchBudgets(c *gin.Context) {
	page, err := h.budgets.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) searchAdjustments(c *gin.Context) {
	page, err := h.adjustments.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}

type periodRequest struct {
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func parseDate(v string) time.Time { t, _ := time.Parse("2006-01-02", v); return t }
func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return ports.Query{Search: c.Query("search"), Page: page, PageSize: size, SortBy: c.Query("sort_by"), SortOrder: c.Query("sort_order"), Filters: map[string]string{"kind": c.Query("kind"), "status": c.Query("status")}}
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
	v, _ := c.Get(string(contextkeys.TenantID))
	s, _ := v.(string)
	return s
}
func actorID(c *gin.Context) string {
	v, _ := c.Get(string(contextkeys.ActorID))
	s, _ := v.(string)
	return s
}
