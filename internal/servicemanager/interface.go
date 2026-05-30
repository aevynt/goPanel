package servicemanager

type ServiceStatus string

const (
	StatusActive    ServiceStatus = "active"
	StatusInactive  ServiceStatus = "inactive"
	StatusFailed    ServiceStatus = "failed"
	StatusUnknown   ServiceStatus = "unknown"
)

type Service struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Status       ServiceStatus `json:"status"`
	Port         int           `json:"port,omitempty"`
	BinaryPath   string        `json:"binary_path,omitempty"`
	PanelManaged bool          `json:"panel_managed"`
}

type ServiceSpec struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BinaryPath  string `json:"binary_path"`
	WorkingDir  string `json:"working_dir"`
	Port        int    `json:"port"`
	EnvVars     string `json:"env_vars"`
	Args        string `json:"args"`
	AutoStart   bool   `json:"auto_start"`
	RunAs       string `json:"run_as"`
}

type LogLine struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type Manager interface {
	List() ([]Service, error)
	Get(name string) (*Service, error)
	Start(name string) error
	Stop(name string) error
	Restart(name string) error
	Enable(name string) error
	Disable(name string) error
	Create(spec ServiceSpec) error
	Remove(name string) error
	Logs(name string, tail int) ([]LogLine, error)
	IsInstalled(name string) (bool, error)
}
