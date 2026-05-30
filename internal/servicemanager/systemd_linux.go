//go:build linux

package servicemanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type SystemdManager struct {
	unitDir string
}

func New() *SystemdManager {
	return &SystemdManager{
		unitDir: "/etc/systemd/system",
	}
}

func (m *SystemdManager) systemctl(args ...string) (string, error) {
	cmd := exec.Command("systemctl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("systemctl %s: %w\n%s", strings.Join(args, " "), err, string(out))
	}
	return string(out), nil
}

func (m *SystemdManager) List() ([]Service, error) {
	out, err := m.systemctl("list-units", "--type=service", "--all", "--no-legend", "--no-pager")
	if err != nil {
		return nil, err
	}
	var services []Service
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		services = append(services, Service{
			Name:   strings.TrimSuffix(fields[0], ".service"),
			Status: ServiceStatus(fields[2]),
		})
	}
	return services, nil
}

func (m *SystemdManager) Get(name string) (*Service, error) {
	svcName := name + ".service"
	out, err := m.systemctl("show", svcName, "--property=ActiveState,Description")
	if err != nil {
		return nil, fmt.Errorf("service %s not found", name)
	}
	s := &Service{Name: name}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if strings.HasPrefix(line, "ActiveState=") {
			s.Status = ServiceStatus(strings.TrimPrefix(line, "ActiveState="))
		}
		if strings.HasPrefix(line, "Description=") {
			s.Description = strings.TrimPrefix(line, "Description=")
		}
	}
	return s, nil
}

func (m *SystemdManager) Start(name string) error {
	_, err := m.systemctl("start", name+".service")
	return err
}

func (m *SystemdManager) Stop(name string) error {
	_, err := m.systemctl("stop", name+".service")
	return err
}

func (m *SystemdManager) Restart(name string) error {
	_, err := m.systemctl("restart", name+".service")
	return err
}

func (m *SystemdManager) Enable(name string) error {
	_, err := m.systemctl("enable", name+".service")
	return err
}

func (m *SystemdManager) Disable(name string) error {
	_, err := m.systemctl("disable", name+".service")
	return err
}

const unitTemplate = `[Unit]
Description={{.Description}}
After=network.target

[Service]
Type=simple
User={{.RunAs}}
WorkingDirectory={{.WorkingDir}}
ExecStart={{.BinaryPath}} {{.Args}}
{{if .EnvVars}}Environment={{.EnvVars}}{{end}}
{{if .Port}}Environment=PORT={{.Port}}{{end}}
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
`

func sanitizeSystemdValue(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return strings.TrimSpace(s)
}

func (m *SystemdManager) Create(spec ServiceSpec) error {
	spec.Name = sanitizeSystemdValue(spec.Name)
	spec.Description = sanitizeSystemdValue(spec.Description)
	spec.BinaryPath = sanitizeSystemdValue(spec.BinaryPath)
	spec.WorkingDir = sanitizeSystemdValue(spec.WorkingDir)
	spec.EnvVars = sanitizeSystemdValue(spec.EnvVars)
	spec.Args = sanitizeSystemdValue(spec.Args)
	spec.RunAs = sanitizeSystemdValue(spec.RunAs)

	unitPath := filepath.Join(m.unitDir, spec.Name+".service")
	tmpl, err := template.New("unit").Parse(unitTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, spec); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	if err := os.WriteFile(unitPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write unit file: %w", err)
	}
	if _, err := m.systemctl("daemon-reload"); err != nil {
		return err
	}
	if spec.AutoStart {
		_, err = m.systemctl("enable", spec.Name+".service")
		return err
	}
	return nil
}

func (m *SystemdManager) Remove(name string) error {
	if _, err := m.systemctl("stop", name+".service"); err != nil {
		// ignore stop errors
	}
	if _, err := m.systemctl("disable", name+".service"); err != nil {
		// ignore disable errors
	}
	unitPath := filepath.Join(m.unitDir, name+".service")
	if err := os.Remove(unitPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove unit file: %w", err)
	}
	_, err := m.systemctl("daemon-reload")
	return err
}

func (m *SystemdManager) Logs(name string, tail int) ([]LogLine, error) {
	tailStr := strconv.Itoa(tail)
	cmd := exec.Command("journalctl", "-u", name+".service", "--no-pager", "-n", tailStr, "-o", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		// journalctl may fail due to permissions — try fallback with sudo
		if !bytes.HasPrefix(out, []byte("{")) {
			cmd2 := exec.Command("sudo", "journalctl", "-u", name+".service", "--no-pager", "-n", tailStr, "-o", "json")
			out, err = cmd2.CombinedOutput()
			if err != nil {
				return nil, fmt.Errorf("journalctl: %w\n%s", err, string(out))
			}
		} else {
			return nil, fmt.Errorf("journalctl: %w\n%s", err, string(out))
		}
	}
	var lines []LogLine
	for _, raw := range bytes.Split(bytes.TrimSpace(out), []byte("\n")) {
		if len(raw) == 0 {
			continue
		}
		var entry struct {
			Message        string `json:"MESSAGE"`
			Realtime       string `json:"__REALTIME_TIMESTAMP"`
			Monotonic      string `json:"__MONOTONIC_TIMESTAMP"`
		}
		if err := json.Unmarshal(raw, &entry); err != nil {
			continue
		}
		line := LogLine{Message: entry.Message}
		if ts, err := strconv.ParseInt(entry.Realtime, 10, 64); err == nil && ts > 0 {
			line.Timestamp = time.UnixMicro(ts).Format(time.RFC3339)
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func (m *SystemdManager) IsInstalled(name string) (bool, error) {
	unitPath := filepath.Join(m.unitDir, name+".service")
	_, err := os.Stat(unitPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
