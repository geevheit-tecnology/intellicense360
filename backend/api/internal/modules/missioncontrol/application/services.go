package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sort"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/ports"
)

type CommandService struct {
	items       ports.CommandItemRepository
	events      ports.CommandEventRepository
	actions     ports.CommandActionRepository
	snapshots   ports.OperationalSnapshotRepository
	idempotency ports.IdempotencyRepository
}

func NewCommandService(items ports.CommandItemRepository, events ports.CommandEventRepository, actions ports.CommandActionRepository, snapshots ports.OperationalSnapshotRepository, idempotency ports.IdempotencyRepository) CommandService {
	return CommandService{items: items, events: events, actions: actions, snapshots: snapshots, idempotency: idempotency}
}

func (s CommandService) Create(ctx context.Context, item domain.CommandItem, actorID string) (domain.CommandItem, error) {
	if err := validateItem(item); err != nil {
		return domain.CommandItem{}, err
	}
	if item.IdempotencyKey != "" {
		exists, err := s.idempotency.Exists(ctx, item.TenantID, item.IdempotencyKey)
		if err != nil || exists {
			if exists {
				return domain.CommandItem{}, ErrDuplicateCommandItem
			}
			return domain.CommandItem{}, err
		}
	}
	now := time.Now().UTC()
	if item.Fingerprint == "" {
		item.Fingerprint = domain.Fingerprint(item.TenantID, item.Type, item.Category, item.SourceType, item.SourceID, item.Title)
	}
	if _, err := s.items.FindByFingerprint(ctx, item.TenantID, item.Fingerprint); err == nil {
		return domain.CommandItem{}, ErrDuplicateCommandItem
	}
	item.ID = newID("mci")
	item.Status = domain.StatusOpen
	item.SLAStatus = domain.SLA(item.DueAt, now)
	if item.DetectedAt.IsZero() {
		item.DetectedAt = now
	}
	item.CreatedAt, item.UpdatedAt, item.Version = now, now, 1
	if item.Metadata == nil {
		item.Metadata = map[string]string{}
	}
	saved, err := s.items.Create(ctx, item)
	if err != nil {
		return domain.CommandItem{}, err
	}
	_ = s.record(ctx, saved, "created", "", saved.Status, actorID, nil)
	if item.IdempotencyKey != "" {
		_ = s.idempotency.Save(ctx, item.TenantID, item.IdempotencyKey, saved.ID)
	}
	return saved, nil
}

func (s CommandService) Update(ctx context.Context, item domain.CommandItem, actorID string) (domain.CommandItem, error) {
	current, err := s.items.GetByID(ctx, item.TenantID, item.ID)
	if err != nil {
		return domain.CommandItem{}, err
	}
	if err := validateItem(item); err != nil {
		return domain.CommandItem{}, err
	}
	item.Status, item.CreatedAt, item.UpdatedAt, item.Version = current.Status, current.CreatedAt, time.Now().UTC(), current.Version+1
	item.SLAStatus = domain.SLA(item.DueAt, item.UpdatedAt)
	saved, err := s.items.Update(ctx, item)
	if err != nil {
		return domain.CommandItem{}, err
	}
	payload := map[string]string{}
	if current.Priority != saved.Priority {
		_ = s.record(ctx, saved, "priority_changed", current.Status, saved.Status, actorID, payload)
	}
	if current.Severity != saved.Severity {
		_ = s.record(ctx, saved, "severity_changed", current.Status, saved.Status, actorID, payload)
	}
	if current.AssignedTo != saved.AssignedTo {
		_ = s.record(ctx, saved, "assigned", current.Status, saved.Status, actorID, payload)
	}
	return saved, nil
}

func (s CommandService) Get(ctx context.Context, tenantID, id string) (domain.CommandItem, error) {
	return s.items.GetByID(ctx, tenantID, id)
}
func (s CommandService) List(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.CommandItem], error) {
	page, err := s.items.List(ctx, tenantID, normalize(q))
	if err != nil {
		return page, err
	}
	SortCommandItems(page.Data)
	return page, nil
}
func (s CommandService) Acknowledge(ctx context.Context, tenantID, id, actorID string) (domain.CommandItem, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.StatusAcknowledged, "acknowledged")
}
func (s CommandService) Start(ctx context.Context, tenantID, id, actorID string) (domain.CommandItem, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.StatusInProgress, "started")
}
func (s CommandService) Resolve(ctx context.Context, tenantID, id, actorID string) (domain.CommandItem, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.StatusResolved, "resolved")
}
func (s CommandService) Dismiss(ctx context.Context, tenantID, id, actorID string) (domain.CommandItem, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.StatusDismissed, "dismissed")
}
func (s CommandService) transition(ctx context.Context, tenantID, id, actorID string, to domain.CommandStatus, eventType string) (domain.CommandItem, error) {
	item, err := s.items.GetByID(ctx, tenantID, id)
	if err != nil {
		return domain.CommandItem{}, err
	}
	if !domain.CanTransition(item.Status, to) {
		return domain.CommandItem{}, ErrInvalidStatusTransition
	}
	from := item.Status
	now := time.Now().UTC()
	item.Status, item.UpdatedAt, item.Version = to, now, item.Version+1
	if to == domain.StatusAcknowledged {
		item.AcknowledgedAt = &now
	}
	if to == domain.StatusResolved || to == domain.StatusDismissed {
		item.ResolvedAt = &now
	}
	saved, err := s.items.Update(ctx, item)
	if err != nil {
		return domain.CommandItem{}, err
	}
	_ = s.record(ctx, saved, eventType, from, to, actorID, nil)
	return saved, nil
}

func (s CommandService) CreateAction(ctx context.Context, action domain.CommandAction) (domain.CommandAction, error) {
	if action.TenantID == "" || action.CommandItemID == "" || action.Type == "" {
		return domain.CommandAction{}, ErrInvalidData
	}
	if _, err := s.items.GetByID(ctx, action.TenantID, action.CommandItemID); err != nil {
		return domain.CommandAction{}, err
	}
	now := time.Now().UTC()
	action.ID, action.CreatedAt, action.UpdatedAt = newID("mca"), now, now
	if action.Status == "" {
		action.Status = domain.ActionPending
	}
	if action.Priority == "" {
		action.Priority = domain.PriorityNormal
	}
	return s.actions.Create(ctx, action)
}
func (s CommandService) ListActions(ctx context.Context, tenantID, itemID string, q ports.Query) (ports.Page[domain.CommandAction], error) {
	return s.actions.ListByItem(ctx, tenantID, itemID, normalize(q))
}
func (s CommandService) History(ctx context.Context, tenantID, itemID string, q ports.Query) (ports.Page[domain.CommandEvent], error) {
	return s.events.ListByItem(ctx, tenantID, itemID, normalize(q))
}

func (s CommandService) Summary(ctx context.Context, tenantID string) (domain.MissionControlSummary, error) {
	items, err := s.items.AllOpen(ctx, tenantID)
	if err != nil {
		return domain.MissionControlSummary{}, err
	}
	snapshot := BuildSnapshot(tenantID, items)
	return domain.MissionControlSummary{
		TotalOpen: snapshot.OpenItems, Critical: snapshot.CriticalItems, High: countSeverity(items, domain.SeverityHigh), Medium: countSeverity(items, domain.SeverityMedium),
		Risks: countType(items, domain.TypeRisk), Alerts: countType(items, domain.TypeAlert), Incidents: countType(items, domain.TypeIncident), Opportunities: countType(items, domain.TypeOpportunity),
		Recommendations: countType(items, domain.TypeRecommendation), BreachedSLA: snapshot.BreachedSLAs, OperationalScore: snapshot.OperationalScore, RiskScore: snapshot.RiskScore, HealthScore: snapshot.HealthScore, LastUpdated: snapshot.SnapshotAt,
	}, nil
}

func (s CommandService) RebuildSnapshot(ctx context.Context, tenantID string) (domain.OperationalSnapshot, error) {
	items, err := s.items.AllOpen(ctx, tenantID)
	if err != nil {
		return domain.OperationalSnapshot{}, err
	}
	return s.snapshots.Create(ctx, BuildSnapshot(tenantID, items))
}
func (s CommandService) LatestSnapshot(ctx context.Context, tenantID string) (domain.OperationalSnapshot, error) {
	return s.snapshots.Latest(ctx, tenantID)
}
func (s CommandService) EvaluateRecommendations(ctx context.Context, tenantID string) ([]domain.CommandAction, error) {
	items, err := s.items.AllOpen(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	actions := []domain.CommandAction{}
	for _, item := range items {
		if item.Severity == domain.SeverityCritical || item.SLAStatus == domain.SLABreached || (item.RiskScore >= 0.7 && item.ImpactScore >= 0.7) {
			action, err := s.CreateAction(ctx, domain.CommandAction{TenantID: tenantID, CommandItemID: item.ID, Type: domain.ActionReview, Label: "Review command item", Priority: item.Priority})
			if err != nil {
				return actions, err
			}
			actions = append(actions, action)
		}
	}
	return actions, nil
}

func (s CommandService) record(ctx context.Context, item domain.CommandItem, eventType string, from, to domain.CommandStatus, actorID string, payload map[string]string) error {
	_, err := s.events.Create(ctx, domain.CommandEvent{ID: newID("mce"), TenantID: item.TenantID, CommandItemID: item.ID, EventType: eventType, PreviousStatus: from, NewStatus: to, ActorID: actorID, Payload: payload, OccurredAt: time.Now().UTC()})
	return err
}

func BuildSnapshot(tenantID string, items []domain.CommandItem) domain.OperationalSnapshot {
	now := time.Now().UTC()
	risk := AggregateRisk(items).OverallRiskScore
	health := HealthScore(items)
	return domain.OperationalSnapshot{ID: newID("mcs"), TenantID: tenantID, SnapshotAt: now, OpenItems: len(items), CriticalItems: countSeverity(items, domain.SeverityCritical), HighPriorityItems: countPriority(items, domain.PriorityHigh) + countPriority(items, domain.PriorityUrgent) + countPriority(items, domain.PriorityCritical), ActiveRisks: countType(items, domain.TypeRisk), ActiveAlerts: countType(items, domain.TypeAlert), OpenIncidents: countType(items, domain.TypeIncident), Opportunities: countType(items, domain.TypeOpportunity), BreachedSLAs: countSLA(items, domain.SLABreached), OperationalScore: OperationalScore(items), RiskScore: risk, HealthScore: health}
}

func SortCommandItems(items []domain.CommandItem) {
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i], items[j]
		if domain.SeverityRank(a.Severity) != domain.SeverityRank(b.Severity) {
			return domain.SeverityRank(a.Severity) > domain.SeverityRank(b.Severity)
		}
		if domain.PriorityRank(a.Priority) != domain.PriorityRank(b.Priority) {
			return domain.PriorityRank(a.Priority) > domain.PriorityRank(b.Priority)
		}
		if a.RiskScore != b.RiskScore {
			return a.RiskScore > b.RiskScore
		}
		if a.ImpactScore != b.ImpactScore {
			return a.ImpactScore > b.ImpactScore
		}
		if slaRank(a.SLAStatus) != slaRank(b.SLAStatus) {
			return slaRank(a.SLAStatus) > slaRank(b.SLAStatus)
		}
		if !a.DetectedAt.Equal(b.DetectedAt) {
			return a.DetectedAt.Before(b.DetectedAt)
		}
		return a.ID < b.ID
	})
}

func OperationalScore(items []domain.CommandItem) float64 {
	score := 100.0 - (AggregateRisk(items).OverallRiskScore * 40) - (float64(countSeverity(items, domain.SeverityCritical)) * 12) - (float64(countSLA(items, domain.SLABreached)) * 8)
	return clamp(score, 0, 100)
}
func HealthScore(items []domain.CommandItem) float64 {
	return clamp(100-(float64(len(items))*2)-(float64(countSeverity(items, domain.SeverityCritical))*15)-(float64(countSLA(items, domain.SLABreached))*10), 0, 100)
}
func AggregateRisk(items []domain.CommandItem) domain.RiskAggregation {
	top := append([]domain.CommandItem{}, items...)
	SortCommandItems(top)
	if len(top) > 5 {
		top = top[:5]
	}
	total := 0.0
	weight := 0.0
	for _, item := range items {
		w := float64(domain.SeverityRank(item.Severity)) / 5
		total += item.RiskScore * item.ImpactScore * item.Confidence * w
		weight += w
	}
	score := 0.0
	if weight > 0 {
		score = clamp(total/weight, 0, 1)
	}
	return domain.RiskAggregation{OverallRiskScore: score, RiskLevel: domain.RiskLevel(score), TopRisks: top}
}

func validateItem(item domain.CommandItem) error {
	if item.TenantID == "" || item.Title == "" {
		return ErrInvalidData
	}
	if !validType(item.Type) {
		return ErrInvalidCommandItemType
	}
	if !validCategory(item.Category) {
		return ErrInvalidData
	}
	if !validSeverity(item.Severity) {
		return ErrInvalidSeverity
	}
	if !validPriority(item.Priority) {
		return ErrInvalidPriority
	}
	if item.RiskScore < 0 || item.RiskScore > 1 {
		return ErrInvalidRiskScore
	}
	if item.ImpactScore < 0 || item.ImpactScore > 1 {
		return ErrInvalidImpactScore
	}
	if item.Confidence < 0 || item.Confidence > 1 {
		return ErrInvalidConfidence
	}
	return nil
}

func validType(v domain.CommandItemType) bool {
	switch v {
	case domain.TypeAlert, domain.TypeRisk, domain.TypeIncident, domain.TypeOpportunity, domain.TypeRecommendation, domain.TypeTask, domain.TypeAnomaly, domain.TypeWarning, domain.TypeInsight:
		return true
	default:
		return false
	}
}
func validCategory(v domain.Category) bool {
	switch v {
	case domain.CategoryOperational, domain.CategoryMaintenance, domain.CategoryFleet, domain.CategoryTire, domain.CategoryFuel, domain.CategoryInventory, domain.CategoryFinancial, domain.CategoryCompliance, domain.CategoryDriver, domain.CategoryDocument, domain.CategoryCIOT, domain.CategorySafety, domain.CategoryPerformance, domain.CategoryCost:
		return true
	default:
		return false
	}
}
func validSeverity(v domain.Severity) bool {
	return v == domain.SeverityInfo || v == domain.SeverityLow || v == domain.SeverityMedium || v == domain.SeverityHigh || v == domain.SeverityCritical
}
func validPriority(v domain.Priority) bool {
	return v == domain.PriorityLow || v == domain.PriorityNormal || v == domain.PriorityHigh || v == domain.PriorityUrgent || v == domain.PriorityCritical
}

func countSeverity(items []domain.CommandItem, v domain.Severity) int {
	c := 0
	for _, i := range items {
		if i.Severity == v {
			c++
		}
	}
	return c
}
func countPriority(items []domain.CommandItem, v domain.Priority) int {
	c := 0
	for _, i := range items {
		if i.Priority == v {
			c++
		}
	}
	return c
}
func countType(items []domain.CommandItem, v domain.CommandItemType) int {
	c := 0
	for _, i := range items {
		if i.Type == v {
			c++
		}
	}
	return c
}
func countSLA(items []domain.CommandItem, v domain.SLAStatus) int {
	c := 0
	for _, i := range items {
		if i.SLAStatus == v {
			c++
		}
	}
	return c
}
func slaRank(v domain.SLAStatus) int {
	if v == domain.SLABreached {
		return 4
	}
	if v == domain.SLAAtRisk {
		return 3
	}
	if v == domain.SLAWithin {
		return 2
	}
	return 1
}
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
func normalize(q ports.Query) ports.Query {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PerPage <= 0 {
		q.PerPage = 25
	}
	if q.PerPage > 100 {
		q.PerPage = 100
	}
	if q.Filters == nil {
		q.Filters = map[string]string{}
	}
	return q
}
func newID(prefix string) string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return prefix + "_" + time.Now().UTC().Format("20060102150405")
	}
	return prefix + "_" + hex.EncodeToString(b[:])
}
