package main

import (
	"crypto/rand"
	"embed"
	"encoding/hex"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"


	"github.com/lhqua/gopanel/internal/api"
	"github.com/lhqua/gopanel/internal/auth"
	"github.com/lhqua/gopanel/internal/caddy"
	"github.com/lhqua/gopanel/internal/config"
	"github.com/lhqua/gopanel/internal/database"
	"github.com/lhqua/gopanel/internal/filemanager"
	"github.com/lhqua/gopanel/internal/ports"
	"github.com/lhqua/gopanel/internal/servicemanager"
)

//go:embed web/dist web/dist/assets/_*
var webFS embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	defaultConfig := "/etc/gopanel.json"
	if runtime.GOOS == "windows" {
		defaultConfig = "./gopanel.json"
	}
	configPath := flag.String("config", defaultConfig, "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.JWTSecret == "change-me-in-production" {
		log.Println("warning: default insecure JWT secret detected, generating a cryptographically secure random one...")
		secret := make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			return fmt.Errorf("failed to generate random JWT secret: %w", err)
		}
		cfg.JWTSecret = hex.EncodeToString(secret)
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save generated JWT secret: %w", err)
		}
		log.Println("secure JWT secret generated and saved successfully to config")
	}

	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data dir: %w", err)
	}
	if err := os.MkdirAll(cfg.BinariesDir, 0755); err != nil {
		return fmt.Errorf("failed to create binaries dir: %w", err)
	}
	if err := os.MkdirAll(cfg.PublicDir, 0755); err != nil {
		return fmt.Errorf("failed to create public dir: %w", err)
	}
	if err := os.MkdirAll(cfg.PublicSitesDir, 0755); err != nil {
		return fmt.Errorf("failed to create public-sites dir: %w", err)
	}

	dbPath := filepath.Join(cfg.DataDir, "gopanel.db")
	db, err := database.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	ensureDefaultAdmin(db)

	sm := servicemanager.New()
	fm := filemanager.New(cfg.DataDir)
	cc := caddy.NewClient(cfg.CaddyAdminURL)
	pm := ports.NewManager()

	subFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Printf("warning: embedded frontend not found, serving API only")
		subFS = nil
	}

	srv := api.NewServer(cfg, db, sm, fm, cc, pm, subFS)
	srv.StartWShub()

	// Start public file server on public port
	publicAddr := fmt.Sprintf(":%d", cfg.PublicPort)
	go func() {
		log.Printf("public server starting on %s", publicAddr)
		if err := http.ListenAndServe(publicAddr, srv.PublicServer()); err != nil {
			log.Printf("public server error: %v", err)
		}
	}()

	// Set up Caddy reverse proxy for public domain if configured
	if cfg.PublicDomain != "" {
		cc.AddSite(caddy.Site{
			Domain:      cfg.PublicDomain,
			ServicePort: cfg.PublicPort,
			TLSEnabled:  false,
			Type:        "proxy",
		})
		log.Printf("public domain %s -> :%d", cfg.PublicDomain, cfg.PublicPort)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("goPanel starting on %s", addr)
	log.Printf("data directory: %s", cfg.DataDir)
	log.Printf("binaries directory: %s", cfg.BinariesDir)

	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func ensureDefaultAdmin(db *database.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		hash, _ := auth.HashPassword("admin")
		db.Exec("INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)",
			"admin", hash, "admin")
		log.Println("default admin user created (username: admin, password: admin)")
	}
}
