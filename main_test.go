package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandlerServesFrontendAndAPI(t *testing.T) {
	os.Setenv("STORAGE_DRIVER", "memory")
	os.Unsetenv("DATABASE_URL")

	handler, err := newHandler()
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		path, contains string
		status         int
	}{
		{"/", "<!DOCTYPE html", http.StatusOK},
		{"/dashboard", "<!DOCTYPE html", http.StatusOK},
		{"/src/app.js", "class App", http.StatusOK},
		{"/src/api/client.js", "class ApiClient", http.StatusOK},
		{"/src/styles/theme.css", "", http.StatusOK},
		{"/health", `"status":"ok"`, http.StatusOK},
		{"/favicon.ico", "", http.StatusNoContent},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if rec.Code != tc.status || !strings.Contains(rec.Body.String(), tc.contains) {
			t.Errorf("%s -> %d, esperado %d contendo %q; corpo: %.120s", tc.path, rec.Code, tc.status, tc.contains, rec.Body.String())
		}
	}
}

