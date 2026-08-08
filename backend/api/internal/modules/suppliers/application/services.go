package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/ports"
)

var emailPattern = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type SupplierService struct {
	repo       ports.SupplierRepository
	categories ports.CategoryRepository
	audit      ports.AuditRecorder
}

func NewSupplierService(repo ports.SupplierRepository, categories ports.CategoryRepository, audit ports.AuditRecorder) SupplierService {
	return SupplierService{repo: repo, categories: categories, audit: audit}
}

func (s SupplierService) Create(ctx context.Context, supplier domain.Supplier) (domain.Supplier, error) {
	if supplier.Status == "" {
		supplier.Status = domain.StatusDraft
	}
	if supplier.Type == "" {
		supplier.Type = domain.TypeOther
	}
	if err := s.validate(ctx, supplier, ""); err != nil {
		return domain.Supplier{}, err
	}
	stamp(&supplier.ID, &supplier.CreatedAt, &supplier.UpdatedAt, &supplier.Version, "sup")
	saved, err := s.repo.Create(ctx, supplier)
	if err == nil && s.audit != nil {
		_ = s.audit.RecordSupplierEvent(ctx, supplier.TenantID, "", "supplier.created", saved.ID)
	}
	return saved, err
}
func (s SupplierService) FindByID(ctx context.Context, tenantID string, id string) (domain.Supplier, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}
func (s SupplierService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Supplier], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s SupplierService) Update(ctx context.Context, supplier domain.Supplier) (domain.Supplier, error) {
	current, err := s.repo.FindByID(ctx, supplier.TenantID, supplier.ID)
	if err != nil {
		return domain.Supplier{}, err
	}
	if current.Status == domain.StatusArchived && supplier.Status != domain.StatusArchived {
		return domain.Supplier{}, ErrInvalidTransition
	}
	if supplier.Status == "" {
		supplier.Status = current.Status
	}
	if supplier.Type == "" {
		supplier.Type = current.Type
	}
	if err := s.validate(ctx, supplier, supplier.ID); err != nil {
		return domain.Supplier{}, err
	}
	supplier.CreatedAt = current.CreatedAt
	supplier.UpdatedAt = time.Now().UTC()
	supplier.Version = current.Version + 1
	saved, err := s.repo.Update(ctx, supplier)
	if err == nil && s.audit != nil {
		action := "supplier.updated"
		if current.Status != saved.Status {
			action = "supplier.status_changed"
		}
		_ = s.audit.RecordSupplierEvent(ctx, saved.TenantID, "", action, saved.ID)
	}
	return saved, err
}
func (s SupplierService) Delete(ctx context.Context, tenantID string, id string) error {
	err := s.repo.Delete(ctx, tenantID, id)
	if err == nil && s.audit != nil {
		_ = s.audit.RecordSupplierEvent(ctx, tenantID, "", "supplier.deleted", id)
	}
	return err
}
func (s SupplierService) validate(ctx context.Context, supplier domain.Supplier, exceptID string) error {
	if supplier.TenantID == "" || strings.TrimSpace(supplier.LegalName) == "" {
		return ErrValidation
	}
	if !validStatus(supplier.Status) {
		return ErrInvalidStatus
	}
	if !validType(supplier.Type) {
		return ErrInvalidType
	}
	if supplier.Email != "" && !emailPattern.MatchString(supplier.Email) {
		return ErrInvalidEmail
	}
	if supplier.StateRegistration != "" && !validUF(supplier.StateRegistration) {
		return ErrInvalidState
	}
	if supplier.CategoryID != "" && s.categories != nil {
		ok, err := s.categories.Exists(ctx, supplier.TenantID, supplier.CategoryID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrValidation
		}
	}
	cnpj := onlyDigits(supplier.CNPJ)
	if supplier.CNPJ != "" {
		if !validCNPJ(cnpj) {
			return ErrInvalidCNPJ
		}
		if ok, err := s.repo.ExistsCNPJ(ctx, supplier.TenantID, cnpj, exceptID); err != nil || ok {
			if err != nil {
				return err
			}
			return ErrCNPJTaken
		}
		supplier.CNPJ = cnpj
	}
	return nil
}

type CategoryService struct{ repo ports.CategoryRepository }

func NewCategoryService(repo ports.CategoryRepository) CategoryService {
	return CategoryService{repo: repo}
}
func (s CategoryService) Create(ctx context.Context, item domain.SupplierCategory) (domain.SupplierCategory, error) {
	if item.TenantID == "" || item.Name == "" || item.Code == "" {
		return domain.SupplierCategory{}, ErrValidation
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "sca")
	return s.repo.Create(ctx, item)
}
func (s CategoryService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierCategory], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s CategoryService) Update(ctx context.Context, item domain.SupplierCategory) (domain.SupplierCategory, error) {
	if item.TenantID == "" || item.ID == "" || item.Name == "" || item.Code == "" {
		return domain.SupplierCategory{}, ErrValidation
	}
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s CategoryService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type TypeService struct{ repo ports.TypeRepository }

func NewTypeService(repo ports.TypeRepository) TypeService { return TypeService{repo: repo} }
func (s TypeService) Create(ctx context.Context, item domain.SupplierType) (domain.SupplierType, error) {
	if item.TenantID == "" || item.Name == "" || item.Code == "" {
		return domain.SupplierType{}, ErrValidation
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "sty")
	return s.repo.Create(ctx, item)
}
func (s TypeService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierType], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s TypeService) Update(ctx context.Context, item domain.SupplierType) (domain.SupplierType, error) {
	if item.TenantID == "" || item.ID == "" || item.Name == "" || item.Code == "" {
		return domain.SupplierType{}, ErrValidation
	}
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s TypeService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type ContactService struct {
	repo      ports.ContactRepository
	suppliers ports.SupplierRepository
}

func NewContactService(repo ports.ContactRepository, suppliers ports.SupplierRepository) ContactService {
	return ContactService{repo: repo, suppliers: suppliers}
}
func (s ContactService) Create(ctx context.Context, item domain.SupplierContact) (domain.SupplierContact, error) {
	if item.TenantID == "" || item.SupplierID == "" || item.Name == "" {
		return domain.SupplierContact{}, ErrValidation
	}
	if item.Email != "" && !emailPattern.MatchString(item.Email) {
		return domain.SupplierContact{}, ErrInvalidEmail
	}
	if ok, _ := s.suppliers.Exists(ctx, item.TenantID, item.SupplierID); !ok {
		return domain.SupplierContact{}, ErrNotFound
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "sco")
	return s.repo.Create(ctx, item)
}
func (s ContactService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierContact], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s ContactService) Update(ctx context.Context, item domain.SupplierContact) (domain.SupplierContact, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s ContactService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type AddressService struct {
	repo      ports.AddressRepository
	suppliers ports.SupplierRepository
}

func NewAddressService(repo ports.AddressRepository, suppliers ports.SupplierRepository) AddressService {
	return AddressService{repo: repo, suppliers: suppliers}
}
func (s AddressService) Create(ctx context.Context, item domain.SupplierAddress) (domain.SupplierAddress, error) {
	if item.TenantID == "" || item.SupplierID == "" || item.Street == "" || item.City == "" {
		return domain.SupplierAddress{}, ErrValidation
	}
	if item.State != "" && !validUF(item.State) {
		return domain.SupplierAddress{}, ErrInvalidState
	}
	if item.Country == "" {
		item.Country = "BR"
	}
	if ok, _ := s.suppliers.Exists(ctx, item.TenantID, item.SupplierID); !ok {
		return domain.SupplierAddress{}, ErrNotFound
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "sad")
	return s.repo.Create(ctx, item)
}
func (s AddressService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierAddress], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s AddressService) Update(ctx context.Context, item domain.SupplierAddress) (domain.SupplierAddress, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s AddressService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type DocumentService struct {
	repo      ports.DocumentRepository
	suppliers ports.SupplierRepository
	audit     ports.AuditRecorder
}

func NewDocumentService(repo ports.DocumentRepository, suppliers ports.SupplierRepository, audit ports.AuditRecorder) DocumentService {
	return DocumentService{repo: repo, suppliers: suppliers, audit: audit}
}
func (s DocumentService) Create(ctx context.Context, item domain.SupplierDocument) (domain.SupplierDocument, error) {
	if item.TenantID == "" || item.SupplierID == "" || item.DocumentType == "" {
		return domain.SupplierDocument{}, ErrValidation
	}
	if ok, _ := s.suppliers.Exists(ctx, item.TenantID, item.SupplierID); !ok {
		return domain.SupplierDocument{}, ErrNotFound
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "sdo")
	saved, err := s.repo.Create(ctx, item)
	if err == nil && s.audit != nil {
		_ = s.audit.RecordSupplierEvent(ctx, item.TenantID, "", "supplier.document_added", item.SupplierID)
	}
	return saved, err
}
func (s DocumentService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierDocument], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s DocumentService) Update(ctx context.Context, item domain.SupplierDocument) (domain.SupplierDocument, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s DocumentService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type ContractService struct {
	repo      ports.ContractRepository
	suppliers ports.SupplierRepository
	audit     ports.AuditRecorder
}

func NewContractService(repo ports.ContractRepository, suppliers ports.SupplierRepository, audit ports.AuditRecorder) ContractService {
	return ContractService{repo: repo, suppliers: suppliers, audit: audit}
}
func (s ContractService) Create(ctx context.Context, item domain.SupplierContract) (domain.SupplierContract, error) {
	if item.TenantID == "" || item.SupplierID == "" || item.ContractNumber == "" {
		return domain.SupplierContract{}, ErrValidation
	}
	if ok, _ := s.suppliers.Exists(ctx, item.TenantID, item.SupplierID); !ok {
		return domain.SupplierContract{}, ErrNotFound
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "sct")
	saved, err := s.repo.Create(ctx, item)
	if err == nil && s.audit != nil {
		_ = s.audit.RecordSupplierEvent(ctx, item.TenantID, "", "supplier.contract_added", item.SupplierID)
	}
	return saved, err
}
func (s ContractService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierContract], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s ContractService) Update(ctx context.Context, item domain.SupplierContract) (domain.SupplierContract, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s ContractService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type RatingService struct {
	repo      ports.RatingRepository
	suppliers ports.SupplierRepository
}

func NewRatingService(repo ports.RatingRepository, suppliers ports.SupplierRepository) RatingService {
	return RatingService{repo: repo, suppliers: suppliers}
}
func (s RatingService) Create(ctx context.Context, item domain.SupplierRating) (domain.SupplierRating, error) {
	if item.TenantID == "" || item.SupplierID == "" {
		return domain.SupplierRating{}, ErrValidation
	}
	if ok, _ := s.suppliers.Exists(ctx, item.TenantID, item.SupplierID); !ok {
		return domain.SupplierRating{}, ErrNotFound
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "sra")
	return s.repo.Create(ctx, item)
}
func (s RatingService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierRating], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s RatingService) Update(ctx context.Context, item domain.SupplierRating) (domain.SupplierRating, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s RatingService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func normalize(q ports.Query) ports.Query {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		q.PageSize = 20
	}
	if strings.ToLower(q.SortOrder) != "desc" {
		q.SortOrder = "asc"
	} else {
		q.SortOrder = "desc"
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
func validStatus(status domain.SupplierStatus) bool {
	switch status {
	case domain.StatusDraft, domain.StatusActive, domain.StatusInactive, domain.StatusBlocked, domain.StatusArchived:
		return true
	default:
		return false
	}
}
func validType(kind domain.SupplierTypeCode) bool {
	switch kind {
	case domain.TypePartsSupplier, domain.TypeTireSupplier, domain.TypeFuelSupplier, domain.TypeWorkshop, domain.TypeServiceProvider, domain.TypeTechnologyProvider, domain.TypeInsuranceProvider, domain.TypeTransportProvider, domain.TypeOther:
		return true
	default:
		return false
	}
}
func onlyDigits(value string) string {
	b := strings.Builder{}
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
func validCNPJ(value string) bool { return len(value) == 14 }
func validUF(value string) bool   { uf := strings.ToUpper(strings.TrimSpace(value)); return len(uf) == 2 }
