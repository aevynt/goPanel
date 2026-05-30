package servicemanager

import "fmt"

type MockManager struct {
	services map[string]*Service
}

func NewMockManager() *MockManager {
	return &MockManager{
		services: make(map[string]*Service),
	}
}

func (m *MockManager) List() ([]Service, error) {
	var list []Service
	for _, s := range m.services {
		list = append(list, *s)
	}
	return list, nil
}

func (m *MockManager) Get(name string) (*Service, error) {
	s, ok := m.services[name]
	if !ok {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return s, nil
}

func (m *MockManager) Start(name string) error {
	s, ok := m.services[name]
	if !ok {
		return fmt.Errorf("service %s not found", name)
	}
	s.Status = StatusActive
	return nil
}

func (m *MockManager) Stop(name string) error {
	s, ok := m.services[name]
	if !ok {
		return fmt.Errorf("service %s not found", name)
	}
	s.Status = StatusInactive
	return nil
}

func (m *MockManager) Restart(name string) error {
	if err := m.Stop(name); err != nil {
		return err
	}
	return m.Start(name)
}

func (m *MockManager) Enable(name string) error {
	return nil
}

func (m *MockManager) Disable(name string) error {
	return nil
}

func (m *MockManager) Create(spec ServiceSpec) error {
	if _, ok := m.services[spec.Name]; ok {
		return fmt.Errorf("service %s already exists", spec.Name)
	}
	m.services[spec.Name] = &Service{
		Name:        spec.Name,
		Description: spec.Description,
		Status:      StatusInactive,
		Port:        spec.Port,
		BinaryPath:  spec.BinaryPath,
	}
	return nil
}

func (m *MockManager) Remove(name string) error {
	delete(m.services, name)
	return nil
}

func (m *MockManager) Logs(name string, tail int) ([]LogLine, error) {
	return []LogLine{
		{Timestamp: "2026-01-01T00:00:00Z", Message: "mock log line 1"},
		{Timestamp: "2026-01-01T00:00:01Z", Message: "mock log line 2"},
	}, nil
}

func (m *MockManager) IsInstalled(name string) (bool, error) {
	_, ok := m.services[name]
	return ok, nil
}
