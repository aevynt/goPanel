package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/lhqua/gopanel/internal/auth"
	"github.com/lhqua/gopanel/internal/config"
	"github.com/lhqua/gopanel/internal/database"
)

func Initialize(cfg *config.Config, db *database.DB) error {
	// Ensure required directories exist
	dirs := []string{
		cfg.DataDir,
		cfg.BinariesDir,
		cfg.PublicDir,
		cfg.PublicSitesDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Setup default admin if no users exist
	ensureDefaultAdmin(db)

	return nil
}

func ensureDefaultAdmin(db *database.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		hash, _ := auth.HashPassword("admin")
		_, err := db.Exec("INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)",
			"admin", hash, "admin")
		if err != nil {
			log.Printf("warning: failed to create default admin user: %v", err)
			return
		}
		log.Println("default admin user created (username: admin, password: admin)")
	}
}
