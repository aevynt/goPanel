package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lhqua/gopanel/internal/api"
	"github.com/lhqua/gopanel/internal/config"
)

func StartMonitor(cfg *config.Config, srv *api.Server) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	// Cooldown maps to avoid spamming alerts (30 min cooldown)
	lastAlertTimes := make(map[string]time.Time)
	cooldownDuration := 30 * time.Minute

	shouldAlert := func(alertKey string) bool {
		lastAlert, exists := lastAlertTimes[alertKey]
		if !exists || time.Since(lastAlert) > cooldownDuration {
			lastAlertTimes[alertKey] = time.Now()
			return true
		}
		return false
	}

	for {
		select {
		case <-ticker.C:
			stats := srv.CollectStats()

			if stats.CPU > cfg.AlertCPUThreshold {
				if shouldAlert("cpu") {
					SendNotification(cfg, fmt.Sprintf("⚠️ [goPanel Alert] High CPU usage detected: %.2f%% (Threshold: %.2f%%)", stats.CPU, cfg.AlertCPUThreshold))
				}
			}

			if stats.Memory.UsedPct > cfg.AlertRAMThreshold {
				if shouldAlert("ram") {
					SendNotification(cfg, fmt.Sprintf("⚠️ [goPanel Alert] High RAM usage detected: %.2f%% (Threshold: %.2f%%)", stats.Memory.UsedPct, cfg.AlertRAMThreshold))
				}
			}

			if stats.Disk.UsedPct > cfg.AlertDiskThreshold {
				if shouldAlert("disk") {
					SendNotification(cfg, fmt.Sprintf("⚠️ [goPanel Alert] High Disk usage detected: %.2f%% (Threshold: %.2f%%)", stats.Disk.UsedPct, cfg.AlertDiskThreshold))
				}
			}

			if stats.CPUTemp > 0 && stats.CPUTemp > cfg.AlertTempThreshold {
				if shouldAlert("temp") {
					SendNotification(cfg, fmt.Sprintf("🔥 [goPanel Alert] High CPU Temperature detected: %.1f°C (Threshold: %.1f°C)", stats.CPUTemp, cfg.AlertTempThreshold))
				}
			}
		}
	}
}

func SendNotification(cfg *config.Config, message string) {
	if cfg.DiscordWebhook != "" {
		go func() {
			payload, _ := json.Marshal(map[string]string{"content": message})
			resp, err := http.Post(cfg.DiscordWebhook, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				log.Printf("failed to send discord alert: %v", err)
			} else {
				resp.Body.Close()
			}
		}()
	}

	if cfg.TelegramToken != "" && cfg.TelegramChatID != "" {
		go func() {
			url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", cfg.TelegramToken)
			payload, _ := json.Marshal(map[string]string{
				"chat_id": cfg.TelegramChatID,
				"text":    message,
			})
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				log.Printf("failed to send telegram alert: %v", err)
			} else {
				resp.Body.Close()
			}
		}()
	}
}
