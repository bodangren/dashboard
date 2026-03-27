package api

import "net/http"

// RegisterRoutes registers all API endpoints on mux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/projects", projectsHandler)
	mux.HandleFunc("/api/diff", diffHandler)
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("[]"))
}

func diffHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
