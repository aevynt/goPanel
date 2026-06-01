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

	"github.com/aevynt/goPanel/internal/alerts"
	"github.com/aevynt/goPanel/internal/api"
	"github.com/aevynt/goPanel/internal/apps"
	"github.com/aevynt/goPanel/internal/bootstrap"
	"github.com/aevynt/goPanel/internal/caddy"
	"github.com/aevynt/goPanel/internal/config"
	"github.com/aevynt/goPanel/internal/database"
	"github.com/aevynt/goPanel/internal/docker"
	"github.com/aevynt/goPanel/internal/filemanager"
	"github.com/aevynt/goPanel/internal/ports"
	"github.com/aevynt/goPanel/internal/servicemanager"
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
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(cfg.DataDir, "gopanel.db")
	db, err := database.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Perform bootstrap initialization
	if err := bootstrap.Initialize(cfg, db); err != nil {
		return fmt.Errorf("bootstrap failure: %w", err)
	}

	sm := servicemanager.New()
	fm := filemanager.New(cfg.DataDir)
	cc := caddy.NewClient(cfg.CaddyAdminURL)
	pm := ports.NewManager()

	dockerSvc := docker.NewService(cfg.DataDir)
	appsSvc := apps.NewService(cfg, pm, cc)

	subFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Printf("warning: embedded frontend not found, serving API only")
		subFS = nil
	}

	srv := api.NewServer(cfg, db, sm, fm, cc, pm, dockerSvc, appsSvc, subFS)
	srv.StartWShub()
	go alerts.StartMonitor(cfg, srv)

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
