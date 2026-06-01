package apps

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aevynt/goPanel/internal/caddy"
	"github.com/aevynt/goPanel/internal/config"
	"github.com/aevynt/goPanel/internal/docker"
	"github.com/aevynt/goPanel/internal/ports"
)

type Service struct {
	cfg   *config.Config
	pm    *ports.Manager
	caddy *caddy.Client
}

func NewService(cfg *config.Config, pm *ports.Manager, cc *caddy.Client) *Service {
	return &Service{
		cfg:   cfg,
		pm:    pm,
		caddy: cc,
	}
}

func (s *Service) ListCatalog() []AppInfo {
	return Catalog
}

func (s *Service) DeployApp(key string) (int, error) {
	var app *AppInfo
	for _, a := range Catalog {
		if a.Key == key {
			app = &a
			break
		}
	}
	if app == nil {
		return 0, fmt.Errorf("app not found")
	}

	// Find an available port
	port, err := s.pm.FindAvailablePort(app.DefaultPort)
	if err != nil {
		return 0, fmt.Errorf("no available ports: %w", err)
	}

	// Create apps directory
	appDir := s.cfg.DataPath("apps", app.Key)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create app folder: %w", err)
	}

	// Replace PORT in template
	composeContent := strings.ReplaceAll(app.Template, "${PORT}", fmt.Sprintf("%d", port))

	// Write docker-compose.yml
	composePath := filepath.Join(appDir, "docker-compose.yml")
	if err := os.WriteFile(composePath, []byte(composeContent), 0644); err != nil {
		return 0, fmt.Errorf("failed to write compose file: %w", err)
	}

	// Deploy using docker compose
	composeCmd := docker.GetDockerComposeCmd()
	args := append(composeCmd, "-f", composePath, "up", "-d")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = appDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to deploy container: %s", string(output))
	}

	// Register caddy reverse proxy if Domain is configured
	if s.cfg.PanelDomain != "" {
		appDomain := fmt.Sprintf("%s.%s", app.Key, s.cfg.PanelDomain)
		err = s.caddy.AddSite(caddy.Site{
			Domain:      appDomain,
			ServicePort: port,
			TLSEnabled:  true, // Automated Let's Encrypt!
			Type:        "proxy",
		})
		if err != nil {
			return port, fmt.Errorf("Caddy reverse proxy mapping failed: %w", err)
		}
	}

	return port, nil
}
