package api

import (
	"encoding/json"
	"net/http"

	"github.com/aevynt/goPanel/internal/caddy"
)

type SettingsResponse struct {
	PanelDomain  string `json:"panel_domain"`
	Port         int    `json:"port"`
	LogLevel     string `json:"log_level"`
	PublicDomain string `json:"public_domain"`
	PublicPort   int    `json:"public_port"`

	DiscordWebhook     string  `json:"discord_webhook"`
	TelegramToken      string  `json:"telegram_token"`
	TelegramChatID     string  `json:"telegram_chat_id"`
	AlertTempThreshold float64 `json:"alert_temp_threshold"`
	AlertCPUThreshold  float64 `json:"alert_cpu_threshold"`
	AlertRAMThreshold  float64 `json:"alert_ram_threshold"`
	AlertDiskThreshold float64 `json:"alert_disk_threshold"`
}

type UpdateSettingsRequest struct {
	PanelDomain  string `json:"panel_domain"`
	PublicDomain string `json:"public_domain"`

	DiscordWebhook     string  `json:"discord_webhook"`
	TelegramToken      string  `json:"telegram_token"`
	TelegramChatID     string  `json:"telegram_chat_id"`
	AlertTempThreshold float64 `json:"alert_temp_threshold"`
	AlertCPUThreshold  float64 `json:"alert_cpu_threshold"`
	AlertRAMThreshold  float64 `json:"alert_ram_threshold"`
	AlertDiskThreshold float64 `json:"alert_disk_threshold"`
}

func (s *Server) GetSettings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, SettingsResponse{
		PanelDomain:        s.cfg.PanelDomain,
		Port:               s.cfg.Port,
		LogLevel:           s.cfg.LogLevel,
		PublicDomain:       s.cfg.PublicDomain,
		PublicPort:         s.cfg.PublicPort,
		DiscordWebhook:     s.cfg.DiscordWebhook,
		TelegramToken:      s.cfg.TelegramToken,
		TelegramChatID:     s.cfg.TelegramChatID,
		AlertTempThreshold: s.cfg.AlertTempThreshold,
		AlertCPUThreshold:  s.cfg.AlertCPUThreshold,
		AlertRAMThreshold:  s.cfg.AlertRAMThreshold,
		AlertDiskThreshold: s.cfg.AlertDiskThreshold,
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

	s.cfg.DiscordWebhook = req.DiscordWebhook
	s.cfg.TelegramToken = req.TelegramToken
	s.cfg.TelegramChatID = req.TelegramChatID
	s.cfg.AlertTempThreshold = req.AlertTempThreshold
	s.cfg.AlertCPUThreshold = req.AlertCPUThreshold
	s.cfg.AlertRAMThreshold = req.AlertRAMThreshold
	s.cfg.AlertDiskThreshold = req.AlertDiskThreshold

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
		PanelDomain:        s.cfg.PanelDomain,
		Port:               s.cfg.Port,
		LogLevel:           s.cfg.LogLevel,
		PublicDomain:       s.cfg.PublicDomain,
		PublicPort:         s.cfg.PublicPort,
		DiscordWebhook:     s.cfg.DiscordWebhook,
		TelegramToken:      s.cfg.TelegramToken,
		TelegramChatID:     s.cfg.TelegramChatID,
		AlertTempThreshold: s.cfg.AlertTempThreshold,
		AlertCPUThreshold:  s.cfg.AlertCPUThreshold,
		AlertRAMThreshold:  s.cfg.AlertRAMThreshold,
		AlertDiskThreshold: s.cfg.AlertDiskThreshold,
	})
}
