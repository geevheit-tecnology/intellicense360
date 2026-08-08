package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
)

type VehicleService struct{ repo ports.VehicleRepository }

func NewVehicleService(repo ports.VehicleRepository) VehicleService {
	return VehicleService{repo: repo}
}

func (s VehicleService) Create(ctx context.Context, vehicle domain.Vehicle) (domain.Vehicle, error) {
	if err := s.validateUnique(ctx, vehicle, ""); err != nil {
		return domain.Vehicle{}, err
	}
	now := time.Now().UTC()
	vehicle.ID = newID("veh")
	if vehicle.Status == "" {
		vehicle.Status = domain.VehicleStatusDraft
	}
	if vehicle.OwnershipType == "" {
		vehicle.OwnershipType = domain.OwnershipOwned
	}
	vehicle.LicensePlate = domain.LicensePlate(vehicle.LicensePlate.Normalized())
	vehicle.CreatedAt = now
	vehicle.UpdatedAt = now
	vehicle.Version = 1
	return s.repo.Create(ctx, vehicle)
}

func (s VehicleService) FindByID(ctx context.Context, tenantID string, vehicleID string) (domain.Vehicle, error) {
	return s.repo.FindByID(ctx, tenantID, vehicleID)
}

func (s VehicleService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.Vehicle], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}

func (s VehicleService) Update(ctx context.Context, vehicle domain.Vehicle) (domain.Vehicle, error) {
	if err := s.validateUnique(ctx, vehicle, vehicle.ID); err != nil {
		return domain.Vehicle{}, err
	}
	vehicle.LicensePlate = domain.LicensePlate(vehicle.LicensePlate.Normalized())
	vehicle.UpdatedAt = time.Now().UTC()
	vehicle.Version++
	return s.repo.Update(ctx, vehicle)
}

func (s VehicleService) Delete(ctx context.Context, tenantID string, vehicleID string) error {
	return s.repo.Delete(ctx, tenantID, vehicleID)
}

func (s VehicleService) ValidateVehicleAccess(ctx context.Context, tenantID string, vehicleID string) error {
	exists, err := s.repo.Exists(ctx, tenantID, vehicleID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}

func (s VehicleService) validateUnique(ctx context.Context, vehicle domain.Vehicle, exceptID string) error {
	if strings.TrimSpace(string(vehicle.LicensePlate)) == "" || strings.TrimSpace(string(vehicle.Chassis)) == "" || strings.TrimSpace(string(vehicle.Renavam)) == "" {
		return ErrValidation
	}
	if ok, err := s.repo.ExistsLicensePlate(ctx, vehicle.TenantID, vehicle.LicensePlate, exceptID); err != nil || ok {
		if err != nil {
			return err
		}
		return ErrLicensePlateTaken
	}
	if ok, err := s.repo.ExistsChassis(ctx, vehicle.TenantID, vehicle.Chassis, exceptID); err != nil || ok {
		if err != nil {
			return err
		}
		return ErrChassisTaken
	}
	if ok, err := s.repo.ExistsRenavam(ctx, vehicle.TenantID, vehicle.Renavam, exceptID); err != nil || ok {
		if err != nil {
			return err
		}
		return ErrRenavamTaken
	}
	return nil
}

type VehicleBrandService struct{ repo ports.VehicleBrandRepository }

func NewVehicleBrandService(repo ports.VehicleBrandRepository) VehicleBrandService {
	return VehicleBrandService{repo: repo}
}
func (s VehicleBrandService) Create(ctx context.Context, brand domain.VehicleBrand) (domain.VehicleBrand, error) {
	stampBase(&brand.ID, &brand.CreatedAt, &brand.UpdatedAt, &brand.Version, "vbr")
	return s.repo.Create(ctx, brand)
}
func (s VehicleBrandService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.VehicleBrand], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}
func (s VehicleBrandService) Update(ctx context.Context, brand domain.VehicleBrand) (domain.VehicleBrand, error) {
	brand.UpdatedAt = time.Now().UTC()
	brand.Version++
	return s.repo.Update(ctx, brand)
}
func (s VehicleBrandService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type VehicleModelService struct{ repo ports.VehicleModelRepository }

func NewVehicleModelService(repo ports.VehicleModelRepository) VehicleModelService {
	return VehicleModelService{repo: repo}
}
func (s VehicleModelService) Create(ctx context.Context, model domain.VehicleModel) (domain.VehicleModel, error) {
	stampBase(&model.ID, &model.CreatedAt, &model.UpdatedAt, &model.Version, "vmo")
	return s.repo.Create(ctx, model)
}
func (s VehicleModelService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.VehicleModel], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}
func (s VehicleModelService) Update(ctx context.Context, model domain.VehicleModel) (domain.VehicleModel, error) {
	model.UpdatedAt = time.Now().UTC()
	model.Version++
	return s.repo.Update(ctx, model)
}
func (s VehicleModelService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type VehicleCategoryService struct {
	repo ports.VehicleCategoryRepository
}

func NewVehicleCategoryService(repo ports.VehicleCategoryRepository) VehicleCategoryService {
	return VehicleCategoryService{repo: repo}
}
func (s VehicleCategoryService) Create(ctx context.Context, category domain.VehicleCategoryEntity) (domain.VehicleCategoryEntity, error) {
	stampBase(&category.ID, &category.CreatedAt, &category.UpdatedAt, &category.Version, "vca")
	return s.repo.Create(ctx, category)
}
func (s VehicleCategoryService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.VehicleCategoryEntity], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}
func (s VehicleCategoryService) Update(ctx context.Context, category domain.VehicleCategoryEntity) (domain.VehicleCategoryEntity, error) {
	category.UpdatedAt = time.Now().UTC()
	category.Version++
	return s.repo.Update(ctx, category)
}
func (s VehicleCategoryService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type VehicleTypeService struct{ repo ports.VehicleTypeRepository }

func NewVehicleTypeService(repo ports.VehicleTypeRepository) VehicleTypeService {
	return VehicleTypeService{repo: repo}
}
func (s VehicleTypeService) Create(ctx context.Context, vehicleType domain.VehicleTypeEntity) (domain.VehicleTypeEntity, error) {
	stampBase(&vehicleType.ID, &vehicleType.CreatedAt, &vehicleType.UpdatedAt, &vehicleType.Version, "vty")
	return s.repo.Create(ctx, vehicleType)
}
func (s VehicleTypeService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.VehicleTypeEntity], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}
func (s VehicleTypeService) Update(ctx context.Context, vehicleType domain.VehicleTypeEntity) (domain.VehicleTypeEntity, error) {
	vehicleType.UpdatedAt = time.Now().UTC()
	vehicleType.Version++
	return s.repo.Update(ctx, vehicleType)
}
func (s VehicleTypeService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type AssetService struct{ repo ports.AssetRepository }

func NewAssetService(repo ports.AssetRepository) AssetService { return AssetService{repo: repo} }
func (s AssetService) Create(ctx context.Context, asset domain.Asset) (domain.Asset, error) {
	stampBase(&asset.ID, &asset.CreatedAt, &asset.UpdatedAt, &asset.Version, "ast")
	if asset.Status == "" {
		asset.Status = "active"
	}
	return s.repo.Create(ctx, asset)
}
func (s AssetService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.Asset], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}
func (s AssetService) Update(ctx context.Context, asset domain.Asset) (domain.Asset, error) {
	asset.UpdatedAt = time.Now().UTC()
	asset.Version++
	return s.repo.Update(ctx, asset)
}
func (s AssetService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func stampBase(id *string, createdAt *time.Time, updatedAt *time.Time, version *int64, prefix string) {
	now := time.Now().UTC()
	*id = newID(prefix)
	*createdAt = now
	*updatedAt = now
	*version = 1
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

func newID(prefix string) string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return prefix + "_" + hex.EncodeToString(b[:])
}
