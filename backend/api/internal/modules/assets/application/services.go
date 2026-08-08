package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/ports"
)

type AssetService struct {
	repo  ports.AssetRepository
	audit ports.AuditRecorder
}

func NewAssetService(repo ports.AssetRepository, audit ports.AuditRecorder) AssetService {
	return AssetService{repo: repo, audit: audit}
}

func (s AssetService) Create(ctx context.Context, asset domain.Asset) (domain.Asset, error) {
	if err := s.validate(ctx, asset, ""); err != nil {
		return domain.Asset{}, err
	}
	now := time.Now().UTC()
	asset.ID = newID("ast")
	if asset.Status == "" {
		asset.Status = domain.AssetStatusDraft
	}
	if asset.Ownership == "" {
		asset.Ownership = domain.OwnershipOwned
	}
	asset.CreatedAt = now
	asset.UpdatedAt = now
	asset.Version = 1
	saved, err := s.repo.Create(ctx, asset)
	if err == nil && s.audit != nil {
		_ = s.audit.RecordAssetEvent(ctx, asset.TenantID, "", "assets.create", saved.ID)
	}
	return saved, err
}

func (s AssetService) FindByID(ctx context.Context, tenantID string, assetID string) (domain.Asset, error) {
	return s.repo.FindByID(ctx, tenantID, assetID)
}

func (s AssetService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.Asset], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}

func (s AssetService) Update(ctx context.Context, asset domain.Asset) (domain.Asset, error) {
	if err := s.validate(ctx, asset, asset.ID); err != nil {
		return domain.Asset{}, err
	}
	current, err := s.repo.FindByID(ctx, asset.TenantID, asset.ID)
	if err != nil {
		return domain.Asset{}, err
	}
	asset.CreatedAt = current.CreatedAt
	asset.Version = current.Version + 1
	asset.UpdatedAt = time.Now().UTC()
	saved, err := s.repo.Update(ctx, asset)
	if err == nil && s.audit != nil {
		_ = s.audit.RecordAssetEvent(ctx, asset.TenantID, "", "assets.update", saved.ID)
	}
	return saved, err
}

func (s AssetService) Delete(ctx context.Context, tenantID string, assetID string) error {
	err := s.repo.Delete(ctx, tenantID, assetID)
	if err == nil && s.audit != nil {
		_ = s.audit.RecordAssetEvent(ctx, tenantID, "", "assets.delete", assetID)
	}
	return err
}

func (s AssetService) validate(ctx context.Context, asset domain.Asset, exceptID string) error {
	if asset.TenantID == "" || strings.TrimSpace(asset.InternalCode) == "" || strings.TrimSpace(asset.AssetTag) == "" || strings.TrimSpace(asset.Name) == "" {
		return ErrValidation
	}
	if asset.Status != "" && !validStatus(asset.Status) {
		return ErrInvalidStatus
	}
	if ok, err := s.repo.ExistsInternalCode(ctx, asset.TenantID, asset.InternalCode, exceptID); err != nil || ok {
		if err != nil {
			return err
		}
		return ErrInternalCodeTaken
	}
	if strings.TrimSpace(asset.SerialNumber) != "" {
		if ok, err := s.repo.ExistsSerialNumber(ctx, asset.TenantID, asset.SerialNumber, exceptID); err != nil || ok {
			if err != nil {
				return err
			}
			return ErrSerialNumberTaken
		}
	}
	if ok, err := s.repo.ExistsAssetTag(ctx, asset.TenantID, asset.AssetTag, exceptID); err != nil || ok {
		if err != nil {
			return err
		}
		return ErrAssetTagTaken
	}
	return nil
}

func validStatus(status domain.AssetStatus) bool {
	switch status {
	case domain.AssetStatusDraft, domain.AssetStatusAvailable, domain.AssetStatusAssigned, domain.AssetStatusInOperation, domain.AssetStatusMaintenance, domain.AssetStatusInactive, domain.AssetStatusSold, domain.AssetStatusDisposed:
		return true
	default:
		return false
	}
}

type CategoryService struct{ repo ports.CategoryRepository }

func NewCategoryService(repo ports.CategoryRepository) CategoryService {
	return CategoryService{repo: repo}
}
func (s CategoryService) Create(ctx context.Context, item domain.AssetCategory) (domain.AssetCategory, error) {
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "aca")
	return s.repo.Create(ctx, item)
}
func (s CategoryService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.AssetCategory], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}
func (s CategoryService) Update(ctx context.Context, item domain.AssetCategory) (domain.AssetCategory, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s CategoryService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type TypeService struct{ repo ports.TypeRepository }

func NewTypeService(repo ports.TypeRepository) TypeService { return TypeService{repo: repo} }
func (s TypeService) Create(ctx context.Context, item domain.AssetType) (domain.AssetType, error) {
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "aty")
	return s.repo.Create(ctx, item)
}
func (s TypeService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.AssetType], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}
func (s TypeService) Update(ctx context.Context, item domain.AssetType) (domain.AssetType, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s TypeService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type ManufacturerService struct{ repo ports.ManufacturerRepository }

func NewManufacturerService(repo ports.ManufacturerRepository) ManufacturerService {
	return ManufacturerService{repo: repo}
}
func (s ManufacturerService) Create(ctx context.Context, item domain.Manufacturer) (domain.Manufacturer, error) {
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "man")
	return s.repo.Create(ctx, item)
}
func (s ManufacturerService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.Manufacturer], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}
func (s ManufacturerService) Update(ctx context.Context, item domain.Manufacturer) (domain.Manufacturer, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s ManufacturerService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type ModelService struct{ repo ports.ModelRepository }

func NewModelService(repo ports.ModelRepository) ModelService { return ModelService{repo: repo} }
func (s ModelService) Create(ctx context.Context, item domain.Model) (domain.Model, error) {
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "mod")
	return s.repo.Create(ctx, item)
}
func (s ModelService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.Model], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}
func (s ModelService) Update(ctx context.Context, item domain.Model) (domain.Model, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s ModelService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type EquipmentService struct{ repo ports.EquipmentRepository }

func NewEquipmentService(repo ports.EquipmentRepository) EquipmentService {
	return EquipmentService{repo: repo}
}
func (s EquipmentService) Create(ctx context.Context, item domain.Equipment) (domain.Equipment, error) {
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "eqp")
	return s.repo.Create(ctx, item)
}
func (s EquipmentService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.Equipment], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}
func (s EquipmentService) Update(ctx context.Context, item domain.Equipment) (domain.Equipment, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s EquipmentService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func normalizeQuery(query ports.Query) ports.Query {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 20
	}
	if query.SortOrder != "desc" {
		query.SortOrder = "asc"
	}
	return query
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
	_, _ = rand.Read(b[:])
	return prefix + "_" + hex.EncodeToString(b[:])
}
