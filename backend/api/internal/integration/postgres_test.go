package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/geevheit/intelligence360/backend/api/internal/core/config"
	"github.com/geevheit/intelligence360/backend/api/internal/core/di"
	identityapp "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/application"
	identitydomain "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
	identity "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/infrastructure/postgres"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func emptyDatabase(t *testing.T) (*sql.DB, string) {
	t.Helper()
	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("TEST_DATABASE_URL not set; PostgreSQL integration not executed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	admin, err := platform.Open(ctx, raw)
	if err != nil {
		t.Fatal("test PostgreSQL unavailable")
	}
	if _, err := admin.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public; CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public`); err != nil {
		t.Fatal(err)
	}
	schema := "test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if _, err := admin.ExecContext(ctx, `CREATE SCHEMA `+schema); err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", schema+",public")
	u.RawQuery = q.Encode()
	db, err := platform.Open(ctx, u.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
		if _, err := admin.ExecContext(context.Background(), `DROP SCHEMA `+schema+` CASCADE`); err != nil {
			t.Error(err)
		}
		admin.Close()
	})
	return db, u.String()
}

func database(t *testing.T) (*sql.DB, string) {
	t.Helper()
	db, dsn := emptyDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	files, err := platform.Migrate(ctx, db, "../../../../database/migrations")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 5 {
		t.Fatalf("expected 5 active runtime migrations, got %d", len(files))
	}
	files, err = platform.Migrate(ctx, db, "../../../../database/migrations")
	if err != nil || len(files) != 0 {
		t.Fatalf("migration replay: %v %v", files, err)
	}
	var count int
	if err := db.QueryRow(`SELECT count(*) FROM tenants`).Scan(&count); err != nil || count != 0 {
		t.Fatalf("unexpected bootstrap tenant: %d %v", count, err)
	}
	return db, dsn
}
func router(db *sql.DB) http.Handler {
	return di.NewContainer(config.Config{AppName: "test", Env: "production", StorageDriver: "postgres", AuthIssuer: "test", AuthSecret: "integration-test-secret-with-32-characters"}).WithDatabase(db).HTTPRouter()
}
func request(t *testing.T, h http.Handler, method, path, token string, body any, want int) map[string]any {
	t.Helper()
	var encoded []byte
	if body != nil {
		var err error
		encoded, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(encoded))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != want {
		t.Fatalf("%s %s: status %d want %d; %s", method, path, w.Code, want, w.Body.String())
	}
	result := map[string]any{}
	if w.Body.Len() > 0 && strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
	}
	return result
}
func login(t *testing.T, h http.Handler, tenant, email string) map[string]any {
	return request(t, h, "POST", "/api/v1/login", "", map[string]any{"tenant_id": tenant, "email": email, "password": "Local-test-12345"}, 200)
}
func TestPostgresRuntime(t *testing.T) {
	db, dsn := database(t)
	ctx := context.Background()
	first, err := identity.Provision(ctx, db, "Empresa A", "a@example.test", "Local-test-12345")
	if err != nil {
		t.Fatal(err)
	}
	second, err := identity.Provision(ctx, db, "Empresa B", "b@example.test", "Local-test-12345")
	if err != nil {
		t.Fatal(err)
	}
	h := router(db)
	a := login(t, h, first.ID, "a@example.test")
	b := login(t, h, second.ID, "b@example.test")
	token := a["access_token"].(string)
	other := b["access_token"].(string)
	request(t, h, "GET", "/api/v1/fleet/vehicles", "", nil, 401)
	request(t, h, "GET", "/api/v1/tenants", token, nil, 403)
	request(t, h, "GET", "/api/v1/fuel/transactions", token, nil, 404)
	catalogs := map[string]map[string]any{}
	for _, kind := range []string{"categories", "brands", "types", "models", "assets"} {
		body := map[string]any{"name": "Teste " + kind, "code": "truck", "status": "active"}
		if kind == "types" {
			body["category_id"] = catalogs["categories"]["id"]
		}
		if kind == "models" {
			body["brand_id"] = catalogs["brands"]["id"]
			body["category_id"] = catalogs["categories"]["id"]
			body["type_id"] = catalogs["types"]["id"]
		}
		catalogs[kind] = request(t, h, "POST", "/api/v1/fleet/"+kind, token, body, 201)
		request(t, h, "GET", "/api/v1/fleet/"+kind+"?search=Teste", token, nil, 200)
	}
	vehicle := map[string]any{"license_plate": "ABC-1D23", "renavam": "12345678901", "chassis": "TESTCHASSIS000001", "category_id": catalogs["categories"]["id"], "type_id": catalogs["types"]["id"], "brand_id": catalogs["brands"]["id"], "model_id": catalogs["models"]["id"], "asset_id": catalogs["assets"]["id"], "tenant_id": second.ID, "metadata": map[string]string{"owner": "test"}}
	saved := request(t, h, "POST", "/api/v1/fleet/vehicles", token, vehicle, 201)
	id := saved["id"].(string)
	if saved["tenant_id"] != first.ID {
		t.Fatal("request body tenant overrode authenticated tenant")
	}
	request(t, h, "POST", "/api/v1/fleet/vehicles", token, vehicle, 409)
	request(t, h, "GET", "/api/v1/fleet/vehicles/"+id, other, nil, 404)
	// Cross-company references must fail in PostgreSQL even when the caller has permissions.
	request(t, h, "POST", "/api/v1/fleet/vehicles", other, vehicle, 400)
	request(t, h, "PUT", "/api/v1/fleet/vehicles/"+id, other, saved, 404)
	request(t, h, "DELETE", "/api/v1/fleet/vehicles/"+id, other, nil, 404)
	saved["color"] = "blue"
	updated := request(t, h, "PUT", "/api/v1/fleet/vehicles/"+id, token, saved, 200)
	if updated["version"] != float64(2) {
		t.Fatalf("version: %v", updated["version"])
	}
	request(t, h, "PUT", "/api/v1/fleet/vehicles/"+id, token, saved, 409)
	// Rebuild the connection pool, stores and router: records and sessions must survive.
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db2, err := platform.Open(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db2.Close()
	h = router(db2)
	persisted := request(t, h, "GET", "/api/v1/fleet/vehicles/"+id, token, nil, 200)
	if persisted["color"] != "blue" {
		t.Fatal("vehicle update did not persist")
	}
	page := request(t, h, "GET", "/api/v1/fleet/vehicles?status=draft&page_size=1&sort_by=license_plate&sort_order=desc", token, nil, 200)
	if page["total_items"] != float64(1) {
		t.Fatalf("page: %v", page)
	}
	// Session refresh must be single-use, including concurrent requests.
	refresh := a["refresh_token"].(string)
	statuses := make(chan int, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			payload, _ := json.Marshal(map[string]string{"refresh_token": refresh})
			req := httptest.NewRequest("POST", "/api/v1/refresh", bytes.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)
			statuses <- w.Code
		}()
	}
	wg.Wait()
	close(statuses)
	good, bad := 0, 0
	for status := range statuses {
		if status == 200 {
			good++
		}
		if status == 401 {
			bad++
		}
	}
	if good != 1 || bad != 1 {
		t.Fatalf("refresh concurrency: success=%d rejected=%d", good, bad)
	}
	request(t, h, "GET", "/api/v1/me", token, nil, 401)
	current := login(t, h, first.ID, "a@example.test")
	token = current["access_token"].(string)
	request(t, h, "POST", "/api/v1/logout", token, map[string]any{"refresh_token": current["refresh_token"]}, 204)
	request(t, h, "GET", "/api/v1/me", token, nil, 401)
	current = login(t, h, first.ID, "a@example.test")
	token = current["access_token"].(string)
	request(t, h, "DELETE", "/api/v1/fleet/vehicles/"+id, token, nil, 204)
	request(t, h, "GET", "/api/v1/fleet/vehicles/"+id, token, nil, 404)
	request(t, h, "POST", "/api/v1/fleet/vehicles", token, vehicle, 201)
	for _, kind := range []string{"models", "types", "brands", "categories", "assets"} {
		item := catalogs[kind]
		id := item["id"].(string)
		request(t, h, "PUT", "/api/v1/fleet/"+kind+"/"+id, other, item, 404)
		item["name"] = "Alterado " + kind
		request(t, h, "PUT", "/api/v1/fleet/"+kind+"/"+id, token, item, 200)
		request(t, h, "PUT", "/api/v1/fleet/"+kind+"/"+id, token, item, 409)
		request(t, h, "DELETE", "/api/v1/fleet/"+kind+"/"+id, other, nil, 404)
		request(t, h, "DELETE", "/api/v1/fleet/"+kind+"/"+id, token, nil, 204)
	}
	// User editing must preserve the password hash and reject disabled accounts.
	user := current["user"].(map[string]any)
	user["name"] = "Edited"
	request(t, h, "PUT", "/api/v1/users/"+user["id"].(string), token, user, 200)
	login(t, h, first.ID, "a@example.test")
	user["status"] = "inactive"
	request(t, h, "PUT", "/api/v1/users/"+user["id"].(string), token, user, 200)
	request(t, h, "GET", "/api/v1/me", token, nil, 401)
	request(t, h, "POST", "/api/v1/login", "", map[string]any{"tenant_id": first.ID, "email": "a@example.test", "password": "Local-test-12345"}, 401)
	t.Log("PASS: migrations/replay, provisioning, catalog CRUD, persistent sessions/vehicles, tenant isolation, optimistic locking, concurrent refresh, logout and inactive users")
}

func TestPostgresTransactionRollback(t *testing.T) {
	db, _ := database(t)
	ctx := context.Background()
	store := identity.NewStore(db)
	err := (platform.Transactions{DB: db}).Run(ctx, func(ctx context.Context) error {
		if _, err := identityapp.NewTenantService(store.Tenants()).Create(ctx, identitydomain.Tenant{Name: "Rollback"}); err != nil {
			return err
		}
		return fmt.Errorf("intentional failure")
	})
	if err == nil {
		t.Fatal("expected rollback error")
	}
	var n int
	if err := db.QueryRow(`SELECT count(*) FROM tenants`).Scan(&n); err != nil || n != 0 {
		t.Fatalf("rollback failed: n=%d err=%v", n, err)
	}
}

func TestPostgresIdentityRelations(t *testing.T) {
	db, _ := database(t)
	ctx := context.Background()
	tenant, err := identity.Provision(ctx, db, "Identity", "identity@example.test", "Local-test-12345")
	if err != nil {
		t.Fatal(err)
	}
	store := identity.NewStore(db)
	users, err := store.Users().List(ctx, tenant.ID)
	if err != nil || len(users) != 1 {
		t.Fatalf("users: %v", err)
	}
	admin, err := store.Users().FindByEmail(ctx, tenant.ID, "IDENTITY@example.test")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(admin.PasswordHash, "$2") {
		t.Fatal("password is not bcrypt")
	}
	roles, err := store.Roles().List(ctx, tenant.ID)
	if err != nil || len(roles) != 1 {
		t.Fatalf("roles: %v", err)
	}
	role, err := store.Roles().FindByID(ctx, tenant.ID, roles[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	role.Description = "Edited"
	if _, err := identityapp.NewRoleService(store.Roles()).Update(ctx, role); err != nil {
		t.Fatal(err)
	}
	permissions, err := store.Roles().PermissionsForUser(ctx, tenant.ID, admin.ID)
	if err != nil || len(permissions) == 0 {
		t.Fatalf("permissions: %v", err)
	}
	for _, p := range permissions {
		if p.Key == "identity.tenants.manage" || p.Key == "identity.permissions.manage" {
			t.Fatal("company admin received platform permission")
		}
	}
	event, err := store.Audit().List(ctx, tenant.ID)
	if err != nil || len(event) != 1 || event[0].ActorID != admin.ID {
		t.Fatalf("audit did not persist: %v", err)
	}
	if _, err := store.Audit().List(ctx, uuid.NewString()); err != nil {
		t.Fatal(err)
	}
	unrelated, err := identity.Provision(ctx, db, "Other", "other@example.test", "Local-test-12345")
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Tenants().AddUser(ctx, unrelated.ID, admin.ID); err == nil {
		t.Fatal("cross-tenant user assignment allowed")
	}
	// Also enforce tenant membership at the database boundary.
	if _, err := db.Exec(`INSERT INTO tenant_users(tenant_id,user_id,role_id) VALUES($1,$2,$3)`, unrelated.ID, admin.ID, role.ID); err == nil {
		t.Fatal("cross-tenant foreign key accepted")
	}
	if err := store.Roles().Delete(ctx, tenant.ID, role.ID); err != nil {
		t.Fatal(err)
	}
	permissions, err = store.Roles().PermissionsForUser(ctx, tenant.ID, admin.ID)
	if err != nil || len(permissions) != 0 {
		t.Fatal("deleted role still authorizes user")
	}
	tenant.Active = false
	if _, err := identityapp.NewTenantService(store.Tenants()).Update(ctx, tenant); err != nil {
		t.Fatal(err)
	}
	request(t, router(db), "POST", "/api/v1/login", "", map[string]string{"tenant_id": tenant.ID, "email": admin.Email, "password": "Local-test-12345"}, 401)
}
