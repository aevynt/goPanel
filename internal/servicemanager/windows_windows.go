//go:build windows

package servicemanager

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type WindowsServiceManager struct{}

func New() *WindowsServiceManager {
	return &WindowsServiceManager{}
}

func (m *WindowsServiceManager) sc(args ...string) (string, error) {
	cmd := exec.Command("sc", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("sc %s: %w\n%s", strings.Join(args, " "), err, string(out))
	}
	return string(out), nil
}

func (m *WindowsServiceManager) List() ([]Service, error) {
	cmd := exec.Command("sc", "query", "type=", "service", "state=", "all")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("sc query: %w", err)
	}
	var services []Service
	lines := strings.Split(string(out), "\n")
	for i, line := range lines {
		if strings.Contains(line, "SERVICE_NAME:") {
			parts := strings.SplitN(line, ":", 2)
			name := strings.TrimSpace(parts[1])
			status := StatusUnknown
			for j := i + 1; j < len(lines) && j < i+5; j++ {
				if strings.Contains(lines[j], "STATE") {
					if strings.Contains(lines[j], "RUNNING") {
						status = StatusActive
					} else if strings.Contains(lines[j], "STOPPED") {
						status = StatusInactive
					}
					break
				}
			}
			services = append(services, Service{
				Name:   name,
				Status: status,
			})
		}
	}
	return services, nil
}

func (m *WindowsServiceManager) Get(name string) (*Service, error) {
	out, err := m.sc("query", name)
	if err != nil {
		return nil, err
	}
	s := &Service{Name: name}
	if strings.Contains(out, "RUNNING") {
		s.Status = StatusActive
	} else if strings.Contains(out, "STOPPED") {
		s.Status = StatusInactive
	} else {
		s.Status = StatusUnknown
	}
	return s, nil
}

func (m *WindowsServiceManager) Start(name string) error {
	_, err := m.sc("start", name)
	return err
}

func (m *WindowsServiceManager) Stop(name string) error {
	_, err := m.sc("stop", name)
	return err
}

func (m *WindowsServiceManager) Restart(name string) error {
	if err := m.Stop(name); err != nil {
		return err
	}
	return m.Start(name)
}

func (m *WindowsServiceManager) Enable(name string) error {
	_, err := m.sc("config", name, "start=", "auto")
	return err
}

func (m *WindowsServiceManager) Disable(name string) error {
	_, err := m.sc("config", name, "start=", "disabled")
	return err
}

func (m *WindowsServiceManager) Create(spec ServiceSpec) error {
	binPath := spec.BinaryPath
	if spec.Args != "" {
		binPath = fmt.Sprintf(`"%s" %s`, spec.BinaryPath, spec.Args)
	}
	_, err := m.sc("create", spec.Name, fmt.Sprintf("binPath= %s", binPath), "start= auto")
	return err
}

func (m *WindowsServiceManager) Remove(name string) error {
	_, err := m.sc("delete", name)
	return err
}

func (m *WindowsServiceManager) Logs(name string, tail int) ([]LogLine, error) {
	maxEvents := tail * 20
	if maxEvents < 200 {
		maxEvents = 200
	}
	ps := fmt.Sprintf(
		`Get-WinEvent -FilterHashtable @{LogName='System'; ProviderName='Service Control Manager'} -MaxEvents %d | Where-Object { $_.Message -like '*%s*' } | Select-Object -First %d TimeCreated, Message | Sort-Object TimeCreated | ConvertTo-Json -Compress`,
		maxEvents, strings.ReplaceAll(name, "'", "''"), tail,
	)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", ps)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return make([]LogLine, 0), nil
	}

	text := strings.TrimSpace(string(out))
	if text == "" || text == "[]" {
		return make([]LogLine, 0), nil
	}

	logs := make([]LogLine, 0)
	if strings.HasPrefix(text, "[") {
		json.Unmarshal([]byte(text), &logs)
	} else {
		var entry LogLine
		if json.Unmarshal([]byte(text), &entry) == nil {
			logs = append(logs, entry)
		}
	}
	return logs, nil
}

func (m *WindowsServiceManager) IsInstalled(name string) (bool, error) {
	_, err := m.sc("query", name)
	if err != nil {
		return false, nil
	}
	return true, nil
}
