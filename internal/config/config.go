package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type Config struct {
	mu sync.Mutex `json:"-"`

	Port           int    `json:"port"`
	DataDir        string `json:"data_dir"`
	BinariesDir    string `json:"binaries_dir"`
	JWTSecret      string `json:"jwt_secret"`
	JWTExpiry      int    `json:"jwt_expiry"`
	CaddyAdminURL  string `json:"caddy_admin_url"`
	PanelDomain    string `json:"panel_domain"`
	LogLevel       string `json:"log_level"`
	PublicDir      string `json:"public_dir"`
	PublicSitesDir string `json:"public_sites_dir"`
	PublicPort     int    `json:"public_port"`
	PublicDomain   string `json:"public_domain"`

	configPath string `json:"-"`
}

func Default() *Config {
	dataDir := "/var/lib/gopanel"
	binariesDir := "/var/lib/gopanel/binaries"
	publicDir := "/var/lib/gopanel/public"
	publicSitesDir := "/var/lib/gopanel/public-sites"

	if runtime.GOOS == "windows" {
		dataDir = "./data"
		binariesDir = "./data/binaries"
		publicDir = "./data/public"
		publicSitesDir = "./data/public-sites"
	}

	return &Config{
		Port:           3636,
		DataDir:        dataDir,
		BinariesDir:    binariesDir,
		JWTSecret:      "change-me-in-production",
		JWTExpiry:      24,
		CaddyAdminURL:  "http://localhost:2019",
		PanelDomain:    "",
		LogLevel:       "info",
		PublicDir:      publicDir,
		PublicSitesDir: publicSitesDir,
		PublicPort:     3637,
		PublicDomain:   "",
	}
}

func (c *Config) DataPath(sub ...string) string {
	return filepath.Join(append([]string{c.DataDir}, sub...)...)
}

func (c *Config) PublicSitesPath(sub ...string) string {
	return filepath.Join(append([]string{c.PublicSitesDir}, sub...)...)
}

func Load(path string) (*Config, error) {
	cfg := Default()
	cfg.configPath = path
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, cfg.Save()
		}
		return nil, fmt.Errorf("read config: %w", err)
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func (c *Config) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(c.configPath), 0755); err != nil {
		return fmt.Errorf("mkdir config dir: %w", err)
	}
	return os.WriteFile(c.configPath, data, 0644)
}
