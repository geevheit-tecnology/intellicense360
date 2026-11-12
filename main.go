// Geevheit Intelligence 360 — hosted server (Vercel Go preset entrypoint).
//
// Vercel detects this root main.go and runs it as a long-running server. It
// serves the REST API (/api/*, /health) through the modular Go monolith and the
// embedded Mission Control SPA for every other path. For Docker/VPS deploys use
// backend/api/cmd/api plus the Nginx frontend instead (docker-compose.prod.yml).
package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	api "github.com/geevheit/intelligence360/backend/api"
)

//go:embed apps/mission-control/index.html apps/mission-control/src
var webFS embed.FS

func newHandler() (http.Handler, error) {
	web, err := fs.Sub(webFS, "apps/mission-control")
	if err != nil {
		return nil, err
	}
	indexHTML, err := fs.ReadFile(web, "index.html")
	if err != nil {
		return nil, err
	}
	static := http.FileServer(http.FS(web))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasPrefix(path, "/api/") || path == "/health":
			api.ServeHTTP(w, r)
		case strings.HasPrefix(path, "/src/"):
			static.ServeHTTP(w, r)
		case path == "/favicon.ico":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-cache")
			_, _ = w.Write(indexHTML)
		}
	}), nil
}

func main() {
	handler, err := newHandler()
	if err != nil {
		log.Fatalf("frontend embutido inválido: %v", err)
	}

	// Warm the API (database connection) in background without blocking startup.
	go func() {
		if _, err := api.Handler(); err != nil {
			log.Printf("aquecimento da API falhou (será repetido na próxima requisição): %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("Geevheit Intelligence 360 ouvindo em :%s", port)
	log.Fatal(server.ListenAndServe())
}

