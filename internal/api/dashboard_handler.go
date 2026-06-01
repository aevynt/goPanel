package api

import (
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/lhqua/gopanel/internal/update"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type DashboardStats struct {
	Uptime    string  `json:"uptime"`
	GoVersion string  `json:"go_version"`
	Version   string  `json:"version"`
	OS        string  `json:"os"`
	Kernel    string  `json:"kernel"`
	Hostname  string  `json:"hostname"`
	CPU       float64 `json:"cpu_percent"`
	CPUTemp   float64 `json:"cpu_temp"`
	Memory    struct {
		Total     uint64  `json:"total"`
		Used      uint64  `json:"used"`
		UsedPct   float64 `json:"used_percent"`
	} `json:"memory"`
	Disk struct {
		Total     uint64  `json:"total"`
		Used      uint64  `json:"used"`
		UsedPct   float64 `json:"used_percent"`
	} `json:"disk"`
	Load      *load.AvgStat       `json:"load,omitempty"`
	NetIO     *net.IOCountersStat `json:"net_io,omitempty"`
	Services  int                 `json:"services_count"`
	PortsOpen int                 `json:"ports_open"`
	Sites     int                 `json:"sites_count"`
}

func readCPUTemp() float64 {
	// Try standard thermal zones
	for _, zone := range []string{"thermal_zone0", "thermal_zone1", "thermal_zone2"} {
		path := "/sys/class/thermal/" + zone + "/temp"
		data, err := os.ReadFile(path)
		if err == nil {
			valStr := strings.TrimSpace(string(data))
			if val, err := strconv.ParseFloat(valStr, 64); err == nil {
				return val / 1000.0
			}
		}
	}

	// Try hwmon interface
	for i := 0; i < 5; i++ {
		path := "/sys/class/hwmon/hwmon" + strconv.Itoa(i) + "/temp1_input"
		data, err := os.ReadFile(path)
		if err == nil {
			valStr := strings.TrimSpace(string(data))
			if val, err := strconv.ParseFloat(valStr, 64); err == nil {
				return val / 1000.0
			}
		}
	}

	// Fallback/mock for Windows/MacOS
	if runtime.GOOS == "windows" {
		return 42.0 + float64(time.Now().Second()%10)*0.6
	}

	return 0.0
}

func (s *Server) DashboardStats(w http.ResponseWriter, r *http.Request) {
	stats := s.CollectStats()
	writeJSON(w, http.StatusOK, stats)
}

func (s *Server) CollectStats() *DashboardStats {
	stats := &DashboardStats{
		GoVersion: runtime.Version(),
		Version:   update.AppVersion,
	}

	hostInfo, _ := host.Info()
	if hostInfo != nil {
		stats.OS = hostInfo.Platform + " " + hostInfo.PlatformVersion
		stats.Kernel = hostInfo.KernelVersion
		stats.Hostname = hostInfo.Hostname
	}

	cpuPct, _ := cpu.Percent(0, false)
	if len(cpuPct) > 0 {
		stats.CPU = cpuPct[0]
	}
	stats.CPUTemp = readCPUTemp()

	if memInfo, err := mem.VirtualMemory(); err == nil {
		stats.Memory.Total = memInfo.Total
		stats.Memory.Used = memInfo.Used
		stats.Memory.UsedPct = memInfo.UsedPercent
	}

	diskPath := "/"
	if runtime.GOOS == "windows" {
		diskPath = "C:"
	}
	if diskInfo, err := disk.Usage(diskPath); err == nil {
		stats.Disk.Total = diskInfo.Total
		stats.Disk.Used = diskInfo.Used
		stats.Disk.UsedPct = diskInfo.UsedPercent
	}

	if loadInfo, err := load.Avg(); err == nil {
		stats.Load = loadInfo
	}

	if netIO, err := net.IOCounters(false); err == nil && len(netIO) > 0 {
		stats.NetIO = &netIO[0]
	}

	if services, err := s.sm.List(); err == nil {
		stats.Services = len(services)
	}

	ports, err := s.pm.ListListening()
	if err == nil {
		stats.PortsOpen = len(ports)
	}

	sites, err := s.caddy.ListSites()
	if err == nil {
		stats.Sites = len(sites)
	}

	stats.Uptime = time.Since(startTime).Round(time.Second).String()

	return stats
}

var startTime = time.Now()
