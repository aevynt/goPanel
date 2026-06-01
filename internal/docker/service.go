package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Container struct {
	ID     string `json:"id"`
	Names  string `json:"names"`
	Image  string `json:"image"`
	State  string `json:"state"`
	Status string `json:"status"`
	Ports  string `json:"ports"`
}

type Service struct {
	dataDir string
}

func NewService(dataDir string) *Service {
	return &Service{
		dataDir: dataDir,
	}
}

func (s *Service) ListContainers() ([]Container, error) {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{json .}}")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run docker command: %s", strings.TrimSpace(stderr.String()))
	}

	containers := make([]Container, 0)
	lines := strings.Split(stdout.String(), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var c struct {
			ID     string `json:"ID"`
			Names  string `json:"Names"`
			Image  string `json:"Image"`
			State  string `json:"State"`
			Status string `json:"Status"`
			Ports  string `json:"Ports"`
		}
		if err := json.Unmarshal([]byte(line), &c); err == nil {
			containers = append(containers, Container{
				ID:     c.ID,
				Names:  c.Names,
				Image:  c.Image,
				State:  c.State,
				Status: c.Status,
				Ports:  c.Ports,
			})
		}
	}
	return containers, nil
}

func (s *Service) StartContainer(id string) error {
	cmd := exec.Command("docker", "start", id)
	return cmd.Run()
}

func (s *Service) StopContainer(id string) error {
	cmd := exec.Command("docker", "stop", id)
	return cmd.Run()
}

func (s *Service) RestartContainer(id string) error {
	cmd := exec.Command("docker", "restart", id)
	return cmd.Run()
}

func (s *Service) GetContainerLogs(id string) (string, error) {
	cmd := exec.Command("docker", "logs", "--tail", "100", id)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	_ = cmd.Run() // ignore exit errors (sometimes stderr is not empty)

	logs := stdout.String()
	if logs == "" && stderr.Len() > 0 {
		logs = stderr.String()
	}
	return logs, nil
}

func (s *Service) DeployCompose(name, content string) error {
	// Validate project name to prevent directory traversal
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			return fmt.Errorf("invalid name: only alphanumeric, dashes, and underscores allowed")
		}
	}

	appDir := filepath.Join(s.dataDir, "compose", name)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create deployment directory: %w", err)
	}

	composeFile := filepath.Join(appDir, "docker-compose.yml")
	if err := os.WriteFile(composeFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write docker-compose.yml: %w", err)
	}

	// Deploy compose in the background
	go func() {
		cmd := exec.Command("docker", "compose", "up", "-d")
		cmd.Dir = appDir
		_ = cmd.Run()
	}()

	return nil
}
