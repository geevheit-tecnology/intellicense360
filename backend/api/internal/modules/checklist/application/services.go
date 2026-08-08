package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/ports"
)

type ChecklistService struct{ repo ports.ChecklistRepository }

func NewChecklistService(repo ports.ChecklistRepository) ChecklistService {
	return ChecklistService{repo: repo}
}

func (s ChecklistService) Create(ctx context.Context, checklist domain.Checklist) (domain.Checklist, error) {
	if strings.TrimSpace(checklist.TenantID) == "" || strings.TrimSpace(checklist.VehicleID) == "" || strings.TrimSpace(checklist.Title) == "" {
		return domain.Checklist{}, ErrValidation
	}
	now := time.Now().UTC()
	checklist.ID = newID("chk")
	if checklist.Status == "" {
		checklist.Status = domain.ChecklistStatusDraft
	}
	checklist.CreatedAt = now
	checklist.UpdatedAt = now
	return s.repo.Create(ctx, checklist)
}

func (s ChecklistService) Update(ctx context.Context, checklist domain.Checklist) (domain.Checklist, error) {
	if strings.TrimSpace(checklist.ID) == "" || strings.TrimSpace(checklist.TenantID) == "" || strings.TrimSpace(checklist.Title) == "" {
		return domain.Checklist{}, ErrValidation
	}
	current, err := s.repo.FindByID(ctx, checklist.TenantID, checklist.ID)
	if err != nil {
		return domain.Checklist{}, err
	}
	checklist.CreatedAt = current.CreatedAt
	checklist.StartedAt = current.StartedAt
	checklist.FinishedAt = current.FinishedAt
	if checklist.Status == "" {
		checklist.Status = current.Status
	}
	checklist.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, checklist)
}

func (s ChecklistService) Delete(ctx context.Context, tenantID string, checklistID string) error {
	return s.repo.Delete(ctx, tenantID, checklistID)
}

func (s ChecklistService) FindByID(ctx context.Context, tenantID string, checklistID string) (domain.Checklist, error) {
	return s.repo.FindByID(ctx, tenantID, checklistID)
}

func (s ChecklistService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.Checklist], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(query))
}

func (s ChecklistService) Start(ctx context.Context, tenantID string, checklistID string, actorID string) (domain.Checklist, error) {
	checklist, err := s.repo.FindByID(ctx, tenantID, checklistID)
	if err != nil {
		return domain.Checklist{}, err
	}
	if checklist.Status != domain.ChecklistStatusDraft {
		return domain.Checklist{}, ErrInvalidStatus
	}
	now := time.Now().UTC()
	checklist.Status = domain.ChecklistStatusInProgress
	checklist.StartedAt = &now
	checklist.UpdatedAt = now
	checklist.UpdatedBy = actorID
	return s.repo.Update(ctx, checklist)
}

func (s ChecklistService) Finish(ctx context.Context, tenantID string, checklistID string, actorID string) (domain.Checklist, error) {
	checklist, err := s.repo.FindByID(ctx, tenantID, checklistID)
	if err != nil {
		return domain.Checklist{}, err
	}
	if checklist.Status != domain.ChecklistStatusInProgress {
		return domain.Checklist{}, ErrInvalidStatus
	}
	now := time.Now().UTC()
	checklist.Status = domain.ChecklistStatusCompleted
	checklist.FinishedAt = &now
	checklist.UpdatedAt = now
	checklist.UpdatedBy = actorID
	return s.repo.Update(ctx, checklist)
}

func (s ChecklistService) Cancel(ctx context.Context, tenantID string, checklistID string, actorID string) (domain.Checklist, error) {
	checklist, err := s.repo.FindByID(ctx, tenantID, checklistID)
	if err != nil {
		return domain.Checklist{}, err
	}
	if checklist.Status == domain.ChecklistStatusCompleted {
		return domain.Checklist{}, ErrInvalidStatus
	}
	now := time.Now().UTC()
	checklist.Status = domain.ChecklistStatusCancelled
	checklist.UpdatedAt = now
	checklist.UpdatedBy = actorID
	return s.repo.Update(ctx, checklist)
}

func (s ChecklistService) ValidateChecklistAccess(ctx context.Context, tenantID string, checklistID string) error {
	exists, err := s.repo.Exists(ctx, tenantID, checklistID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}

type ChecklistItemService struct {
	checklists ports.ChecklistRepository
	items      ports.ChecklistItemRepository
}

func NewChecklistItemService(checklists ports.ChecklistRepository, items ports.ChecklistItemRepository) ChecklistItemService {
	return ChecklistItemService{checklists: checklists, items: items}
}

func (s ChecklistItemService) Add(ctx context.Context, item domain.ChecklistItem) (domain.ChecklistItem, error) {
	if err := s.validateItem(ctx, item); err != nil {
		return domain.ChecklistItem{}, err
	}
	now := time.Now().UTC()
	item.ID = newID("chi")
	if item.AnswerType == "" {
		item.AnswerType = domain.AnswerTypeBoolean
	}
	item.CreatedAt = now
	item.UpdatedAt = now
	return s.items.Create(ctx, item)
}

func (s ChecklistItemService) Update(ctx context.Context, item domain.ChecklistItem) (domain.ChecklistItem, error) {
	if err := s.validateItem(ctx, item); err != nil {
		return domain.ChecklistItem{}, err
	}
	current, err := s.items.FindByID(ctx, item.TenantID, item.ChecklistID, item.ID)
	if err != nil {
		return domain.ChecklistItem{}, err
	}
	item.CreatedAt = current.CreatedAt
	item.UpdatedAt = time.Now().UTC()
	return s.items.Update(ctx, item)
}

func (s ChecklistItemService) Delete(ctx context.Context, tenantID string, checklistID string, itemID string) error {
	return s.items.Delete(ctx, tenantID, checklistID, itemID)
}

func (s ChecklistItemService) List(ctx context.Context, tenantID string, checklistID string, query ports.Query) (ports.Page[domain.ChecklistItem], error) {
	return s.items.List(ctx, tenantID, checklistID, normalizeQuery(query))
}

func (s ChecklistItemService) validateItem(ctx context.Context, item domain.ChecklistItem) error {
	if strings.TrimSpace(item.TenantID) == "" || strings.TrimSpace(item.ChecklistID) == "" || strings.TrimSpace(item.Title) == "" {
		return ErrValidation
	}
	exists, err := s.checklists.Exists(ctx, item.TenantID, item.ChecklistID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}

type ChecklistAnswerService struct {
	items   ports.ChecklistItemRepository
	answers ports.ChecklistAnswerRepository
}

func NewChecklistAnswerService(items ports.ChecklistItemRepository, answers ports.ChecklistAnswerRepository) ChecklistAnswerService {
	return ChecklistAnswerService{items: items, answers: answers}
}

func (s ChecklistAnswerService) AnswerItem(ctx context.Context, answer domain.ChecklistAnswer) (domain.ChecklistAnswer, error) {
	if strings.TrimSpace(answer.TenantID) == "" || strings.TrimSpace(answer.ChecklistID) == "" || strings.TrimSpace(answer.ChecklistItemID) == "" {
		return domain.ChecklistAnswer{}, ErrValidation
	}
	item, err := s.items.FindByID(ctx, answer.TenantID, answer.ChecklistID, answer.ChecklistItemID)
	if err != nil {
		return domain.ChecklistAnswer{}, err
	}
	if answer.Answer == "" && item.AnswerType != domain.AnswerTypePhoto && item.AnswerType != domain.AnswerTypeSignature {
		return domain.ChecklistAnswer{}, ErrInvalidAnswer
	}
	now := time.Now().UTC()
	answer.ID = newID("cha")
	answer.AnsweredAt = now
	answer.CreatedAt = now
	answer.UpdatedAt = now
	return s.answers.Create(ctx, answer)
}

func (s ChecklistAnswerService) List(ctx context.Context, tenantID string, checklistID string, query ports.Query) (ports.Page[domain.ChecklistAnswer], error) {
	return s.answers.List(ctx, tenantID, checklistID, normalizeQuery(query))
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
