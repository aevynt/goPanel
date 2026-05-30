package ports

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	psnet "github.com/shirou/gopsutil/v3/net"
)

type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	PID      int    `json:"pid,omitempty"`
	Process  string `json:"process,omitempty"`
}

type Manager struct {
	minPort int
	maxPort int
}

func NewManager() *Manager {
	return &Manager{
		minPort: 1024,
		maxPort: 65535,
	}
}

var commonPorts = []int{
	21, 22, 23, 25, 53, 80, 110, 135, 139, 143, 443, 445, 465, 587, 993, 995,
	1080, 1433, 1521, 3000, 3306, 3389, 3636, 3637, 5000, 5432, 5672, 6379,
	8000, 8080, 8443, 8888, 9000, 9200, 11211, 27017,
}

func (m *Manager) ListListening() ([]PortInfo, error) {
	conns, err := psnet.Connections("tcp")
	if err != nil {
		return m.ScanCommon()
	}
	ports := make([]PortInfo, 0)
	seen := make(map[int]bool)
	for _, c := range conns {
		if c.Status == "LISTEN" {
			port := int(c.Laddr.Port)
			if seen[port] {
				continue
			}
			seen[port] = true
			p := PortInfo{
				Port:     port,
				Protocol: "tcp",
				State:    "listening",
				PID:      int(c.Pid),
			}
			ports = append(ports, p)
		}
	}
	return ports, nil
}

func (m *Manager) ScanCommon() ([]PortInfo, error) {
	ports := make([]PortInfo, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 50)

	for _, port := range commonPorts {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()
			addr := fmt.Sprintf("localhost:%d", p)
			conn, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
			if err == nil {
				conn.Close()
				mu.Lock()
				ports = append(ports, PortInfo{
					Port:     p,
					Protocol: "tcp",
					State:    "listening",
				})
				mu.Unlock()
			}
		}(port)
	}
	wg.Wait()
	return ports, nil
}

func (m *Manager) Scan() ([]PortInfo, error) {
	ports := make([]PortInfo, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 200)

	for port := m.minPort; port <= m.maxPort; port++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()
			addr := fmt.Sprintf("localhost:%d", p)
			conn, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
			if err == nil {
				conn.Close()
				mu.Lock()
				ports = append(ports, PortInfo{
					Port:     p,
					Protocol: "tcp",
					State:    "listening",
				})
				mu.Unlock()
			}
		}(port)
	}
	wg.Wait()
	return ports, nil
}

func (m *Manager) ScanRange(start, end int) ([]PortInfo, error) {
	if start < 1 {
		start = 1
	}
	if end > 65535 {
		end = 65535
	}
	// Limit scan range to max 1000 ports to prevent system overload
	if end-start > 1000 {
		end = start + 1000
	}
	ports := make([]PortInfo, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 200)

	for port := start; port <= end; port++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			defer wg.Done()
			defer func() { <-sem }()
			addr := fmt.Sprintf("localhost:%d", p)
			conn, err := net.DialTimeout("tcp", addr, 50*time.Millisecond)
			if err == nil {
				conn.Close()
				mu.Lock()
				ports = append(ports, PortInfo{
					Port:     p,
					Protocol: "tcp",
					State:    "listening",
				})
				mu.Unlock()
			}
		}(port)
	}
	wg.Wait()
	return ports, nil
}

func (m *Manager) IsPortAvailable(port int) bool {
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
	if err == nil {
		conn.Close()
		return false
	}
	return true
}

func (m *Manager) FindAvailablePort(preferred int) (int, error) {
	if preferred >= m.minPort && preferred <= m.maxPort && m.IsPortAvailable(preferred) {
		return preferred, nil
	}
	// Dynamically assign an available port via TCP listener on port 0
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("no available ports found: %w", err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	return port, nil
}

func (m *Manager) ParsePortRange(s string) (int, int, error) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) == 1 {
		port, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return 0, 0, fmt.Errorf("invalid port: %s", parts[0])
		}
		return port, port, nil
	}
	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start port: %s", parts[0])
	}
	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end port: %s", parts[1])
	}
	if start > end {
		start, end = end, start
	}
	return start, end, nil
}
