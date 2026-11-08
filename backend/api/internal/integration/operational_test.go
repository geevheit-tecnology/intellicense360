package integration_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	identityapp "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/application"
	identitydomain "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
	identity "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/infrastructure/postgres"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type operationalFixture struct {
	db          *sql.DB
	dsn         string
	handler     http.Handler
	tenant      string
	otherTenant string
	token       string
	otherToken  string
	vehicle     map[string]any
}

func operational(t *testing.T) operationalFixture {
	t.Helper()
	db, dsn := database(t)
	ctx := context.Background()
	tenant, err := identity.Provision(ctx, db, "Operação A", "op-a@example.test", "Local-test-12345")
	if err != nil {
		t.Fatal(err)
	}
	other, err := identity.Provision(ctx, db, "Operação B", "op-b@example.test", "Local-test-12345")
	if err != nil {
		t.Fatal(err)
	}
	h := router(db)
	token := login(t, h, tenant.ID, "op-a@example.test")["access_token"].(string)
	otherToken := login(t, h, other.ID, "op-b@example.test")["access_token"].(string)
	category := request(t, h, "POST", "/api/v1/fleet/categories", token, map[string]any{"name": "Caminhão", "code": "truck"}, 201)
	brand := request(t, h, "POST", "/api/v1/fleet/brands", token, map[string]any{"name": "Teste"}, 201)
	typ := request(t, h, "POST", "/api/v1/fleet/types", token, map[string]any{"name": "Teste", "category_id": category["id"]}, 201)
	model := request(t, h, "POST", "/api/v1/fleet/models", token, map[string]any{"name": "Teste", "brand_id": brand["id"], "category_id": category["id"], "type_id": typ["id"]}, 201)
	vehicle := request(t, h, "POST", "/api/v1/fleet/vehicles", token, map[string]any{"category_id": category["id"], "brand_id": brand["id"], "type_id": typ["id"], "model_id": model["id"], "license_plate": "OPR1T23", "renavam": "00000000000", "chassis": "SYNTHETIC-TEST-01"}, 201)
	return operationalFixture{db, dsn, h, tenant.ID, other.ID, token, otherToken, vehicle}
}
func driverBody(code string) map[string]any {
	return map[string]any{"code": code, "name": "Motorista de teste", "email": "motorista@example.test", "status": "active"}
}
func documentBody(driver any, number string) map[string]any {
	return map[string]any{"driver_id": driver, "kind": "cnh", "title": "CNH cadastrada", "number": number, "issued_on": "2020-01-01"}
}
func TestDriversDocumentsLifecycle(t *testing.T) {
	f := operational(t)
	h := f.handler
	request(t, h, "GET", "/api/v1/drivers", "", nil, 401)
	request(t, h, "GET", "/api/v1/documents", "", nil, 401)
	body := driverBody(" mot-001 ")
	body["vehicle_id"] = f.vehicle["id"]
	body["tenant_id"] = f.otherTenant
	body["tax_id"] = "000.000.000-00"
	driver := request(t, h, "POST", "/api/v1/drivers", f.token, body, 201)
	id := driver["id"].(string)
	if driver["code"] != "MOT-001" || driver["tenant_id"] != f.tenant || driver["tax_id"] != "00000000000" {
		t.Fatal("driver normalization or tenant failed")
	}
	request(t, h, "POST", "/api/v1/drivers", f.token, body, 409)
	request(t, h, "POST", "/api/v1/drivers", f.otherToken, body, 404)
	request(t, h, "GET", "/api/v1/drivers/"+id, f.otherToken, nil, 404)
	request(t, h, "PUT", "/api/v1/drivers/"+id, f.otherToken, driver, 404)
	request(t, h, "DELETE", "/api/v1/drivers/"+id+"?version=1", f.otherToken, nil, 404)
	page := request(t, h, "GET", "/api/v1/drivers?status=active&search=mot-001&vehicle_id="+f.vehicle["id"].(string), f.token, nil, 200)
	if page["total_items"] != float64(1) {
		t.Fatal("driver filters failed")
	}
	driver["name"] = "Nome atualizado"
	driver["vehicle_id"] = ""
	updated := request(t, h, "PUT", "/api/v1/drivers/"+id, f.token, driver, 200)
	if updated["version"] != float64(2) {
		t.Fatal("driver version failed")
	}
	request(t, h, "PUT", "/api/v1/drivers/"+id, f.token, driver, 409)
	docBody := documentBody(id, "CNH-001")
	docBody["tenant_id"] = f.otherTenant
	docBody["expires_on"] = time.Now().UTC().AddDate(0, 0, -1).Format(time.DateOnly)
	document := request(t, h, "POST", "/api/v1/documents", f.token, docBody, 201)
	docID := document["id"].(string)
	if document["expiry_status"] != "expired" || document["tenant_id"] != f.tenant {
		t.Fatal("document expiry or tenant failed")
	}
	request(t, h, "POST", "/api/v1/documents", f.token, docBody, 409)
	request(t, h, "POST", "/api/v1/documents", f.otherToken, docBody, 404)
	request(t, h, "GET", "/api/v1/documents/"+docID, f.otherToken, nil, 404)
	request(t, h, "PUT", "/api/v1/documents/"+docID, f.otherToken, document, 404)
	request(t, h, "DELETE", "/api/v1/documents/"+docID+"?version=1", f.otherToken, nil, 404)
	page = request(t, h, "GET", "/api/v1/documents?driver_id="+id+"&kind=cnh&expiry_status=expired", f.token, nil, 200)
	if page["total_items"] != float64(1) {
		t.Fatal("document expiry filter failed")
	}
	document["expires_on"] = time.Now().UTC().AddDate(0, 0, 90).Format(time.DateOnly)
	renewed := request(t, h, "PUT", "/api/v1/documents/"+docID, f.token, document, 200)
	if renewed["expiry_status"] != "valid" || renewed["version"] != float64(2) {
		t.Fatal("renewal failed")
	}
	request(t, h, "PUT", "/api/v1/documents/"+docID, f.token, document, 409)
	vehicleDoc := request(t, h, "POST", "/api/v1/documents", f.token, map[string]any{"vehicle_id": f.vehicle["id"], "kind": "crlv", "title": "CRLV de teste", "number": "CRLV-001"}, 201)
	if vehicleDoc["expiry_status"] != "no_expiry" {
		t.Fatal("optional expiry failed")
	}
	request(t, h, "POST", "/api/v1/documents", f.otherToken, map[string]any{"vehicle_id": f.vehicle["id"], "kind": "crlv", "title": "CRLV de teste", "number": "CROSS"}, 404)
	soon := documentBody(id, "CNH-002")
	soon["expires_on"] = time.Now().UTC().AddDate(0, 0, 30).Format(time.DateOnly)
	request(t, h, "POST", "/api/v1/documents", f.token, soon, 201)
	page = request(t, h, "GET", "/api/v1/documents?expiry_status=due_soon&page_size=1", f.token, nil, 200)
	if page["total_items"] != float64(1) {
		t.Fatal("inclusive expiry window failed")
	}
	// Recreate all runtime objects, keeping the same session and database contents.
	if err := f.db.Close(); err != nil {
		t.Fatal(err)
	}
	db2, err := platform.Open(context.Background(), f.dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db2.Close()
	h = router(db2)
	persisted := request(t, h, "GET", "/api/v1/drivers/"+id, f.token, nil, 200)
	if persisted["name"] != "Nome atualizado" {
		t.Fatal("driver not durable")
	}
	request(t, h, "GET", "/api/v1/documents/"+docID, f.token, nil, 200)
	// Ownership is immutable even when the replacement owner is in the same company.
	replacement := request(t, h, "POST", "/api/v1/drivers", f.token, driverBody("mot-002"), 201)
	changed := map[string]any{}
	for k, v := range renewed {
		changed[k] = v
	}
	changed["driver_id"] = replacement["id"]
	request(t, h, "PUT", "/api/v1/documents/"+docID, f.token, changed, 400)
	renewed["status"] = "archived"
	archived := request(t, h, "PUT", "/api/v1/documents/"+docID, f.token, renewed, 200)
	if archived["expiry_status"] != "archived" {
		t.Fatal("archive status failed")
	}
	// Archived records preserve history; the same number may be registered again.
	historical := request(t, h, "POST", "/api/v1/documents", f.token, docBody, 201)
	request(t, h, "DELETE", "/api/v1/documents/"+docID+"?version=1", f.token, nil, 409)
	request(t, h, "DELETE", "/api/v1/documents/"+docID+"?version=3", f.token, nil, 204)
	request(t, h, "GET", "/api/v1/documents/"+docID, f.token, nil, 404)
	request(t, h, "DELETE", "/api/v1/drivers/"+id+"?version=1", f.token, nil, 409)
	request(t, h, "DELETE", "/api/v1/drivers/"+id+"?version=2", f.token, nil, 204)
	request(t, h, "GET", "/api/v1/drivers/"+id, f.token, nil, 404)
	request(t, h, "POST", "/api/v1/documents", f.token, documentBody(id, "AFTER-DELETE"), 404)
	// A historical document remains available and can be archived after its owner is soft-deleted.
	historicalID := historical["id"].(string)
	request(t, h, "GET", "/api/v1/documents/"+historicalID, f.token, nil, 200)
	historical["status"] = "archived"
	request(t, h, "PUT", "/api/v1/documents/"+historicalID, f.token, historical, 200)
	request(t, h, "GET", "/api/v1/documents/"+vehicleDoc["id"].(string), f.token, nil, 200)
	var previousDate string
	if err := db2.QueryRow(`SELECT metadata->>'previous_expires_on' FROM audit_logs WHERE resource_id=$1 AND action='document.updated' ORDER BY created_at LIMIT 1`, docID).Scan(&previousDate); err != nil || previousDate != docBody["expires_on"] {
		t.Fatalf("renewal history lost: %q %v", previousDate, err)
	}
	var audited int
	if err := db2.QueryRow(`SELECT count(*) FROM audit_logs WHERE tenant_id=$1 AND resource_type IN ('driver','document')`, f.tenant).Scan(&audited); err != nil || audited < 10 {
		t.Fatalf("durable audit: %d %v", audited, err)
	}
	t.Log("PASS: driver/document CRUD, allocations, expiry, ownership, tenant isolation, version conflicts, soft deletion, durable records and audit")
}

func TestOperationalValidationAndPermissions(t *testing.T) {
	f := operational(t)
	h := f.handler
	for _, body := range []map[string]any{{"code": "", "name": "Teste"}, {"code": "X", "name": "Teste", "status": "unknown"}, {"code": "X", "name": "Teste", "tax_id": "bad"}, {"code": "X", "name": "Teste", "email": "invalid"}, {"code": "X", "name": "Teste", "vehicle_id": "bad"}} {
		request(t, h, "POST", "/api/v1/drivers", f.token, body, 400)
	}
	driver := request(t, h, "POST", "/api/v1/drivers", f.token, driverBody("VALID"), 201)
	id := driver["id"].(string)
	for _, body := range []map[string]any{
		{"kind": "other", "title": "Teste", "number": "1"},
		{"driver_id": id, "vehicle_id": f.vehicle["id"], "kind": "other", "title": "Teste", "number": "1"},
		{"driver_id": id, "kind": "crlv", "title": "Teste", "number": "1"},
		{"driver_id": id, "kind": "cnh", "title": "Teste", "number": "1", "expires_on": "2026-02-30"},
		{"driver_id": id, "kind": "cnh", "title": "Teste", "number": "1", "issued_on": "2026-11-07", "expires_on": "2026-11-06"},
	} {
		request(t, h, "POST", "/api/v1/documents", f.token, body, 400)
	}
	for _, path := range []string{"/api/v1/drivers?page=bad", "/api/v1/drivers?page_size=101", "/api/v1/drivers?status=bad", "/api/v1/documents?expiry_status=bad", "/api/v1/documents?driver_id=bad", "/api/v1/documents?kind=bad"} {
		request(t, h, "GET", path, f.token, nil, 400)
	}
	request(t, h, "DELETE", "/api/v1/drivers/"+id, f.token, nil, 400)
	request(t, h, "GET", "/api/v1/drivers/not-an-id", f.token, nil, 404)
	request(t, h, "GET", "/api/v1/documents/not-an-id", f.token, nil, 404)
	store := identity.NewStore(f.db)
	ctx := context.Background()
	user, err := identityapp.NewUserService(store.Users()).Create(ctx, identitydomain.User{TenantID: f.tenant, Name: "Leitura", Email: "reader@example.test"}, "Local-test-12345")
	if err != nil {
		t.Fatal(err)
	}
	role, err := identityapp.NewRoleService(store.Roles()).Create(ctx, identitydomain.Role{TenantID: f.tenant, Name: "reader"})
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"drivers.read", "documents.read"} {
		if err := store.Roles().AssignPermission(ctx, f.tenant, role.ID, key); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := f.db.Exec(`INSERT INTO tenant_users(tenant_id,user_id,role_id) VALUES($1,$2,$3)`, f.tenant, user.ID, role.ID); err != nil {
		t.Fatal(err)
	}
	reader := login(t, h, f.tenant, user.Email)["access_token"].(string)
	request(t, h, "GET", "/api/v1/drivers", reader, nil, 200)
	request(t, h, "GET", "/api/v1/documents", reader, nil, 200)
	request(t, h, "POST", "/api/v1/drivers", reader, driverBody("BLOCKED"), 403)
	request(t, h, "POST", "/api/v1/documents", reader, documentBody(id, "BLOCKED"), 403)
	request(t, h, "PUT", "/api/v1/drivers/"+id, reader, driver, 403)
	request(t, h, "DELETE", "/api/v1/drivers/"+id+"?version=1", reader, nil, 403)
	if _, err := f.db.Exec(`DELETE FROM role_permissions WHERE role_id=$1`, role.ID); err != nil {
		t.Fatal(err)
	}
	request(t, h, "GET", "/api/v1/drivers", reader, nil, 403)
	request(t, h, "GET", "/api/v1/documents", reader, nil, 403)
	// Composite foreign keys protect data even when bypassing application validation.
	if _, err := f.db.Exec(`INSERT INTO drivers(id,tenant_id,code,name,vehicle_id,created_at,updated_at) VALUES($1,$2,'CROSS','Teste',$3,now(),now())`, uuid.NewString(), f.otherTenant, f.vehicle["id"]); err == nil {
		t.Fatal("cross-tenant vehicle FK accepted")
	}
	if _, err := f.db.Exec(`INSERT INTO documents(id,tenant_id,driver_id,kind,title,number,created_at,updated_at) VALUES($1,$2,$3,'cnh','Teste','CROSS',now(),now())`, uuid.NewString(), f.otherTenant, id); err == nil {
		t.Fatal("cross-tenant owner FK accepted")
	}
}

func TestOperationalAuditRollback(t *testing.T) {
	f := operational(t)
	h := f.handler
	driver := request(t, h, "POST", "/api/v1/drivers", f.token, driverBody("BASE"), 201)
	id := driver["id"].(string)
	doc := request(t, h, "POST", "/api/v1/documents", f.token, documentBody(id, "BASE"), 201)
	docID := doc["id"].(string)
	_, err := f.db.Exec(`CREATE FUNCTION reject_operational_audit() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF NEW.resource_type IN ('driver','document') THEN RAISE EXCEPTION 'synthetic audit failure'; END IF; RETURN NEW; END $$; CREATE TRIGGER reject_operational_audit BEFORE INSERT ON audit_logs FOR EACH ROW EXECUTE FUNCTION reject_operational_audit()`)
	if err != nil {
		t.Fatal(err)
	}
	request(t, h, "POST", "/api/v1/drivers", f.token, driverBody("ROLLBACK"), 500)
	request(t, h, "POST", "/api/v1/documents", f.token, documentBody(id, "ROLLBACK"), 500)
	driver["name"] = "Must not persist"
	request(t, h, "PUT", "/api/v1/drivers/"+id, f.token, driver, 500)
	doc["title"] = "Must not persist"
	request(t, h, "PUT", "/api/v1/documents/"+docID, f.token, doc, 500)
	request(t, h, "DELETE", "/api/v1/drivers/"+id+"?version=1", f.token, nil, 500)
	request(t, h, "DELETE", "/api/v1/documents/"+docID+"?version=1", f.token, nil, 500)
	persisted := request(t, h, "GET", "/api/v1/drivers/"+id, f.token, nil, 200)
	if persisted["name"] == "Must not persist" || persisted["version"] != float64(1) {
		t.Fatal("driver audit rollback failed")
	}
	persisted = request(t, h, "GET", "/api/v1/documents/"+docID, f.token, nil, 200)
	if persisted["title"] == "Must not persist" || persisted["version"] != float64(1) {
		t.Fatal("document audit rollback failed")
	}
	var n int
	if err := f.db.QueryRow(`SELECT (SELECT count(*) FROM drivers WHERE code='ROLLBACK')+(SELECT count(*) FROM documents WHERE number='ROLLBACK')`).Scan(&n); err != nil || n != 0 {
		t.Fatal("creation audit rollback failed")
	}
}

func TestOperationalConcurrentCreate(t *testing.T) {
	f := operational(t)
	for _, path := range []string{"/api/v1/drivers", "/api/v1/documents"} {
		var body any = driverBody("RACE")
		if strings.HasSuffix(path, "documents") {
			driver := request(t, f.handler, "POST", "/api/v1/drivers", f.token, driverBody("OWNER"), 201)
			body = documentBody(driver["id"], "RACE")
		}
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		var wg sync.WaitGroup
		statuses := make(chan int, 2)
		for range 2 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				req := httptest.NewRequest("POST", path, bytes.NewReader(payload))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+f.token)
				w := httptest.NewRecorder()
				f.handler.ServeHTTP(w, req)
				statuses <- w.Code
			}()
		}
		wg.Wait()
		close(statuses)
		created, conflict := 0, 0
		for status := range statuses {
			if status == 201 {
				created++
			}
			if status == 409 {
				conflict++
			}
		}
		if created != 1 || conflict != 1 {
			t.Fatalf("%s created=%d conflict=%d", path, created, conflict)
		}
	}
}

func TestOperationalMigrationUpgrade(t *testing.T) {
	db, _ := emptyDatabase(t)
	ctx := context.Background()
	if _, err := db.Exec(`CREATE TABLE schema_migrations(name text PRIMARY KEY,checksum text NOT NULL,applied_at timestamptz NOT NULL DEFAULT now())`); err != nil {
		t.Fatal(err)
	}
	// Reconstruct the previously delivered runtime exactly; the historical SQL files remain unchanged.
	for _, name := range []string{"000001_init.sql", "000002_identity_platform.sql", "000003_fleet_core.sql", "000020_identity_fleet_tenant_constraints.sql"} {
		body, err := os.ReadFile(filepath.Join("../../../../database/migrations", name))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(string(body)); err != nil {
			t.Fatal(err)
		}
		if _, err := db.Exec(`INSERT INTO schema_migrations(name,checksum) VALUES($1,$2)`, name, fmt.Sprintf("%x", sha256.Sum256(body))); err != nil {
			t.Fatal(err)
		}
	}
	tenant, err := identity.Provision(ctx, db, "Existing customer", "existing@example.test", "Local-test-12345")
	if err != nil {
		t.Fatal(err)
	}
	// Match the previous company's permissions, which predate the new modules.
	if _, err := db.Exec(`DELETE FROM permissions WHERE key IN ('drivers.read','drivers.manage','documents.read','documents.manage')`); err != nil {
		t.Fatal(err)
	}
	h := router(db)
	session := login(t, h, tenant.ID, "existing@example.test")
	token := session["access_token"].(string)
	files, err := platform.Migrate(ctx, db, "../../../../database/migrations")
	if err != nil || len(files) != 1 || files[0] != "000021_drivers_documents_core.sql" {
		t.Fatalf("upgrade: %v %v", files, err)
	}
	request(t, h, "GET", "/api/v1/drivers", token, nil, 200)
	request(t, h, "POST", "/api/v1/drivers", token, driverBody("UPGRADED"), 201)
	request(t, h, "GET", "/api/v1/documents", token, nil, 200)
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM users WHERE tenant_id=$1`, tenant.ID).Scan(&n); err != nil || n != 1 {
		t.Fatal("upgrade modified existing identity")
	}
}
