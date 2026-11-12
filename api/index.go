package handler

import (
	"net/http"

	api "github.com/geevheit/intelligence360/backend/api"
)

// Handler is the official Vercel Go Serverless Function entrypoint.
func Handler(w http.ResponseWriter, r *http.Request) {
	api.ServeHTTP(w, r)
}
