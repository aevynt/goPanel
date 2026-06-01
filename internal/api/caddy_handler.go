package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/aevynt/goPanel/internal/caddy"
)

func isValidDomain(domain string) bool {
	if domain == "" || len(domain) > 255 {
		return false
	}
	for _, char := range domain {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '.' || char == '-' || char == '_') {
			return false
		}
	}
	return true
}

func (s *Server) ListSites(w http.ResponseWriter, r *http.Request) {
	sites, err := s.caddy.ListSites()
	if err != nil {
		http.Error(w, `{"error":"failed to list sites: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, sites)
}

func (s *Server) AddSite(w http.ResponseWriter, r *http.Request) {
	var site caddy.Site
	if err := json.NewDecoder(r.Body).Decode(&site); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if site.Domain == "" {
		http.Error(w, `{"error":"domain is required"}`, http.StatusBadRequest)
		return
	}
	if !isValidDomain(site.Domain) {
		http.Error(w, `{"error":"invalid domain format"}`, http.StatusBadRequest)
		return
	}
	if site.Type == "static" {
		dirName := site.Domain
		site.Root = filepath.Join(s.cfg.PublicSitesDir, dirName)

		if err := os.MkdirAll(site.Root, 0755); err != nil {
			http.Error(w, `{"error":"failed to create site directory"}`, http.StatusInternalServerError)
			return
		}

		indexHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>%s - Site Active</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
           display: flex; align-items: center; justify-content: center;
           min-height: 100vh; background: #141413; color: #e8e6dc; }
    .card { text-align: center; padding: 48px; max-width: 520px; }
    h1 { font-size: 2rem; font-weight: 600; margin-bottom: 12px; letter-spacing: -0.02em; }
    p { color: #9a9990; line-height: 1.6; margin-bottom: 24px; }
    .domain { color: #c96442; font-weight: 500; }
    .badge { display: inline-block; padding: 4px 16px; border-radius: 999px;
             background: #1f1f1d; font-size: 13px; color: #9a9990; border: 1px solid #2a2a28; }
  </style>
</head>
<body>
  <div class="card">
    <div class="badge">goPanel</div>
    <h1>%s</h1>
    <p>This site has been created successfully and is being served via <span class="domain">Caddy</span>.</p>
    <p style="font-size:14px;color:#73726c;">Replace this file with your own content to get started.</p>
  </div>
</body>
</html>`, site.Domain, site.Domain)

		if err := os.WriteFile(filepath.Join(site.Root, "index.html"), []byte(indexHTML), 0644); err != nil {
			http.Error(w, `{"error":"failed to write index.html"}`, http.StatusInternalServerError)
			return
		}
	} else {
		if site.ServicePort == 0 {
			http.Error(w, `{"error":"service_port is required for proxy sites"}`, http.StatusBadRequest)
			return
		}
		site.Type = "proxy"
	}
	if err := s.caddy.AddSite(site); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status": "created"})
}

func (s *Server) RemoveSite(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	if err := s.caddy.RemoveSite(domain); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

func (s *Server) CaddyHealth(w http.ResponseWriter, r *http.Request) {
	if err := s.caddy.Health(); err != nil {
		http.Error(w, `{"error":"caddy not reachable"}`, http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}
