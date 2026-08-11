package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/infrastructure"
)

func TestPublishedTemplateVersionIsImmutableForItems(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	templates := application.NewChecklistTemplateService(store.Templates())
	versions := application.NewChecklistTemplateVersionService(store.Templates(), store.Versions())
	items := application.NewChecklistEngineItemService(store.Versions(), store.EngineItems())

	template, err := templates.Create(ctx, domain.ChecklistTemplate{TenantID: "tenant-a", Name: "Daily"})
	if err != nil {
		t.Fatal(err)
	}
	version, err := versions.Create(ctx, domain.ChecklistTemplateVersion{TenantID: "tenant-a", TemplateID: template.ID})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := versions.Publish(ctx, "tenant-a", version.ID, "user-a"); err != nil {
		t.Fatal(err)
	}
	_, err = items.Create(ctx, domain.ChecklistEngineItem{TenantID: "tenant-a", TemplateVersionID: version.ID, Question: "Freio ok?", ItemType: "yes_no"})
	if !errors.Is(err, application.ErrPublishedVersionImmutable) {
		t.Fatalf("expected immutable published version error, got %v", err)
	}
}

func TestExecutionCannotCompleteWithRequiredItemUnanswered(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	templates := application.NewChecklistTemplateService(store.Templates())
	versions := application.NewChecklistTemplateVersionService(store.Templates(), store.Versions())
	items := application.NewChecklistEngineItemService(store.Versions(), store.EngineItems())
	executions := application.NewChecklistExecutionService(store.Versions(), store.EngineItems(), store.Executions(), store.Responses(), store.Evidence(), store.Signatures(), store.History())

	template, _ := templates.Create(ctx, domain.ChecklistTemplate{TenantID: "tenant-a", Name: "Daily"})
	version, _ := versions.Create(ctx, domain.ChecklistTemplateVersion{TenantID: "tenant-a", TemplateID: template.ID})
	item, err := items.Create(ctx, domain.ChecklistEngineItem{TenantID: "tenant-a", TemplateVersionID: version.ID, Question: "Luzes ok?", ItemType: "yes_no", Required: true})
	if err != nil || item.ID == "" {
		t.Fatalf("create item: %v", err)
	}
	published, _ := versions.Publish(ctx, "tenant-a", version.ID, "user-a")
	execution, err := executions.Start(ctx, domain.ChecklistExecution{TenantID: "tenant-a", TemplateVersionID: published.ID, PerformedBy: "user-a"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = executions.Complete(ctx, "tenant-a", execution.ID, "user-a")
	if !errors.Is(err, application.ErrRequiredItemsUnanswered) {
		t.Fatalf("expected required item error, got %v", err)
	}
}

func TestExecutionCompletesWithRequiredResponseAndEvidence(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	templates := application.NewChecklistTemplateService(store.Templates())
	versions := application.NewChecklistTemplateVersionService(store.Templates(), store.Versions())
	items := application.NewChecklistEngineItemService(store.Versions(), store.EngineItems())
	executions := application.NewChecklistExecutionService(store.Versions(), store.EngineItems(), store.Executions(), store.Responses(), store.Evidence(), store.Signatures(), store.History())
	responses := application.NewChecklistEngineResponseService(store.Executions(), store.EngineItems(), store.Responses(), store.History())
	evidence := application.NewChecklistEvidenceService(store.Evidence(), store.History())

	template, _ := templates.Create(ctx, domain.ChecklistTemplate{TenantID: "tenant-a", Name: "Safety"})
	version, _ := versions.Create(ctx, domain.ChecklistTemplateVersion{TenantID: "tenant-a", TemplateID: template.ID})
	item, _ := items.Create(ctx, domain.ChecklistEngineItem{TenantID: "tenant-a", TemplateVersionID: version.ID, Question: "Extintor ok?", ItemType: "ok_not_ok", Required: true, EvidenceRequired: true})
	published, _ := versions.Publish(ctx, "tenant-a", version.ID, "user-a")
	execution, _ := executions.Start(ctx, domain.ChecklistExecution{TenantID: "tenant-a", TemplateVersionID: published.ID, PerformedBy: "user-a"})
	_, err := responses.Record(ctx, domain.ChecklistResponse{TenantID: "tenant-a", ExecutionID: execution.ID, ItemID: item.ID, Value: "ok", Responder: "user-a"})
	if err != nil {
		t.Fatalf("record response: %v", err)
	}
	_, err = evidence.Add(ctx, domain.ChecklistEvidence{TenantID: "tenant-a", ExecutionID: execution.ID, ResponseID: item.ID, EvidenceType: "photo", Reference: "photo://extintor"})
	if err != nil {
		t.Fatalf("add evidence: %v", err)
	}
	completed, err := executions.Complete(ctx, "tenant-a", execution.ID, "user-a")
	if err != nil {
		t.Fatalf("complete execution: %v", err)
	}
	if completed.Status != domain.ExecutionStatusCompleted {
		t.Fatalf("expected completed, got %s", completed.Status)
	}
}
