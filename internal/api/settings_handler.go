package api

import (
	"encoding/json"
	"net/http"

	"github.com/lhqua/gopanel/internal/caddy"
)

type SettingsResponse struct {
	PanelDomain  string `json:"panel_domain"`
	Port         int    `json:"port"`
	LogLevel     string `json:"log_level"`
	PublicDomain string `json:"public_domain"`
	PublicPort   int    `json:"public_port"`
}

type UpdateSettingsRequest struct {
	PanelDomain  string `json:"panel_domain"`
	PublicDomain string `json:"public_domain"`
}

func (s *Server) GetSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, SettingsResponse{
		PanelDomain:  s.cfg.PanelDomain,
		Port:         s.cfg.Port,
		LogLevel:     s.cfg.LogLevel,
		PublicDomain: s.cfg.PublicDomain,
		PublicPort:   s.cfg.PublicPort,
	})
}

func (s *Server) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req UpdateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	oldPublic := s.cfg.PublicDomain
	s.cfg.PanelDomain = req.PanelDomain
	s.cfg.PublicDomain = req.PublicDomain

	if oldPublic != "" && req.PublicDomain == "" {
		s.caddy.RemoveSite(oldPublic)
	}

	if req.PublicDomain != "" && req.PublicDomain != oldPublic {
		s.caddy.AddSite(caddy.Site{
			Domain:      req.PublicDomain,
			ServicePort: s.cfg.PublicPort,
			TLSEnabled:  false,
			Type:        "proxy",
		})
	}

	if err := s.cfg.Save(); err != nil {
		http.Error(w, `{"error":"failed to save config"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, SettingsResponse{
		PanelDomain:  s.cfg.PanelDomain,
		Port:         s.cfg.Port,
		LogLevel:     s.cfg.LogLevel,
		PublicDomain: s.cfg.PublicDomain,
		PublicPort:   s.cfg.PublicPort,
	})
}
