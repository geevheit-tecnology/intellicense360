package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/ports"
)

type PartService struct {
	repo  ports.PartRepository
	units ports.CatalogRepository
}

func NewPartService(repo ports.PartRepository, units ports.CatalogRepository) PartService {
	return PartService{repo: repo, units: units}
}

func (s PartService) Create(ctx context.Context, part domain.Part) (domain.Part, error) {
	if err := s.validate(ctx, part, ""); err != nil {
		return domain.Part{}, err
	}
	stamp(&part.ID, &part.CreatedAt, &part.UpdatedAt, &part.Version, "prt")
	if part.Status == "" {
		part.Status = domain.StockStatusActive
	}
	return s.repo.Create(ctx, part)
}

func (s PartService) FindByID(ctx context.Context, tenantID string, id string) (domain.Part, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}

func (s PartService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.Part], error) {
	return s.repo.Search(ctx, tenantID, normalize(query))
}

func (s PartService) Update(ctx context.Context, part domain.Part) (domain.Part, error) {
	current, err := s.repo.FindByID(ctx, part.TenantID, part.ID)
	if err != nil {
		return domain.Part{}, err
	}
	if err := s.validate(ctx, part, part.ID); err != nil {
		return domain.Part{}, err
	}
	part.CreatedAt = current.CreatedAt
	part.UpdatedAt = time.Now().UTC()
	part.Version = current.Version + 1
	return s.repo.Update(ctx, part)
}

func (s PartService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func (s PartService) ValidateInventoryItemAccess(ctx context.Context, tenantID string, itemID string) error {
	_, err := s.repo.FindByID(ctx, tenantID, itemID)
	return err
}

func (s PartService) validate(ctx context.Context, part domain.Part, exceptID string) error {
	if part.TenantID == "" || strings.TrimSpace(part.SKU) == "" || strings.TrimSpace(part.InternalCode) == "" || strings.TrimSpace(part.Name) == "" || strings.TrimSpace(part.UnitID) == "" {
		return ErrValidation
	}
	if part.Status != "" && part.Status != domain.StockStatusActive && part.Status != domain.StockStatusInactive {
		return ErrValidation
	}
	if ok, err := s.repo.ExistsSKU(ctx, part.TenantID, part.SKU, exceptID); err != nil || ok {
		if err != nil {
			return err
		}
		return ErrSKUTaken
	}
	if ok, err := s.repo.ExistsInternalCode(ctx, part.TenantID, part.InternalCode, exceptID); err != nil || ok {
		if err != nil {
			return err
		}
		return ErrInternalCodeTaken
	}
	if s.units != nil {
		ok, err := s.units.ExistsCode(ctx, part.TenantID, part.UnitID, "")
		if err != nil {
			return err
		}
		if !ok {
			return ErrInvalidUnit
		}
	}
	if min, max, ok := stockRange(part.Metadata); ok && min > max {
		return ErrInvalidStockRange
	}
	return nil
}

type CatalogService struct{ repo ports.CatalogRepository }

func NewCatalogService(repo ports.CatalogRepository) CatalogService {
	return CatalogService{repo: repo}
}

func (s CatalogService) Create(ctx context.Context, item domain.Catalog) (domain.Catalog, error) {
	if item.TenantID == "" || strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Code) == "" {
		return domain.Catalog{}, ErrValidation
	}
	if ok, err := s.repo.ExistsCode(ctx, item.TenantID, item.Code, ""); err != nil || ok {
		if err != nil {
			return domain.Catalog{}, err
		}
		return domain.Catalog{}, ErrValidation
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "ica")
	return s.repo.Create(ctx, item)
}
func (s CatalogService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Catalog], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s CatalogService) Update(ctx context.Context, item domain.Catalog) (domain.Catalog, error) {
	if item.TenantID == "" || item.ID == "" || item.Name == "" || item.Code == "" {
		return domain.Catalog{}, ErrValidation
	}
	if ok, err := s.repo.ExistsCode(ctx, item.TenantID, item.Code, item.ID); err != nil || ok {
		if err != nil {
			return domain.Catalog{}, err
		}
		return domain.Catalog{}, ErrValidation
	}
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s CatalogService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type WarehouseService struct{ repo ports.WarehouseRepository }

func NewWarehouseService(repo ports.WarehouseRepository) WarehouseService {
	return WarehouseService{repo: repo}
}
func (s WarehouseService) Create(ctx context.Context, item domain.Warehouse) (domain.Warehouse, error) {
	if item.TenantID == "" || item.Name == "" || item.Code == "" {
		return domain.Warehouse{}, ErrValidation
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "whs")
	item.Active = true
	return s.repo.Create(ctx, item)
}
func (s WarehouseService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Warehouse], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s WarehouseService) Update(ctx context.Context, item domain.Warehouse) (domain.Warehouse, error) {
	if item.TenantID == "" || item.ID == "" || item.Name == "" || item.Code == "" {
		return domain.Warehouse{}, ErrValidation
	}
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s WarehouseService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type LocationService struct{ repo ports.LocationRepository }

func NewLocationService(repo ports.LocationRepository) LocationService {
	return LocationService{repo: repo}
}
func (s LocationService) Create(ctx context.Context, item domain.WarehouseLocation) (domain.WarehouseLocation, error) {
	if item.TenantID == "" || item.WarehouseID == "" || item.Name == "" || item.Code == "" {
		return domain.WarehouseLocation{}, ErrValidation
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "wlo")
	return s.repo.Create(ctx, item)
}
func (s LocationService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.WarehouseLocation], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s LocationService) Update(ctx context.Context, item domain.WarehouseLocation) (domain.WarehouseLocation, error) {
	if item.TenantID == "" || item.ID == "" || item.WarehouseID == "" || item.Name == "" || item.Code == "" {
		return domain.WarehouseLocation{}, ErrValidation
	}
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s LocationService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func normalize(q ports.Query) ports.Query {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		q.PageSize = 20
	}
	q.SortOrder = strings.ToLower(q.SortOrder)
	if q.SortOrder != "desc" {
		q.SortOrder = "asc"
	}
	return q
}

func stamp(id *string, createdAt *time.Time, updatedAt *time.Time, version *int64, prefix string) {
	now := time.Now().UTC()
	*id = newID(prefix)
	*createdAt = now
	*updatedAt = now
	*version = 1
}

func newID(prefix string) string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return prefix + "_" + time.Now().UTC().Format("20060102150405")
	}
	return prefix + "_" + hex.EncodeToString(b[:])
}

func stockRange(metadata map[string]string) (float64, float64, bool) {
	if metadata == nil {
		return 0, 0, false
	}
	minRaw, hasMin := metadata["minimum_stock"]
	maxRaw, hasMax := metadata["maximum_stock"]
	if !hasMin || !hasMax {
		return 0, 0, false
	}
	min, errMin := strconv.ParseFloat(minRaw, 64)
	max, errMax := strconv.ParseFloat(maxRaw, 64)
	return min, max, errMin == nil && errMax == nil
}
