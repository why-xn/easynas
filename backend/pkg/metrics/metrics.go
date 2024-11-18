package metrics

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemMetrics represents detailed system metrics.
type SystemMetrics struct {
	TotalCPUs       int     `json:"totalCpus"`       // Total number of CPUs
	CPUUsagePercent float64 `json:"cpuUsagePercent"` // CPU usage percentage
	TotalMemory     uint64  `json:"totalMemory"`     // Total memory in bytes
	MemoryUsage     uint64  `json:"memoryUsed"`      // Used memory in bytes
	MemoryPercent   float64 `json:"memoryPercent"`   // Memory usage percentage
	TotalDisk       uint64  `json:"totalDisk"`       // Total disk space in bytes
	DiskUsage       uint64  `json:"diskUsed"`        // Used disk space in bytes
	DiskPercent     float64 `json:"diskPercent"`     // Disk usage percentage
	Uptime          uint64  `json:"uptime"`          // Uptime in seconds
}

// GetSystemMetrics fetches and returns detailed system metrics.
func GetSystemMetrics() (*SystemMetrics, error) {
	// Get total CPU count
	totalCPUs, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU count: %w", err)
	}

	// Get CPU usage percentage
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %w", err)
	}

	// Get memory stats
	memoryStats, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory stats: %w", err)
	}

	// Get disk stats for the root filesystem
	diskStats, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("failed to get disk stats: %w", err)
	}

	// Get system uptime
	uptime, err := host.Uptime()
	if err != nil {
		return nil, fmt.Errorf("failed to get uptime: %w", err)
	}

	// Assemble metrics into the struct
	metrics := &SystemMetrics{
		TotalCPUs:       totalCPUs,
		CPUUsagePercent: cpuUsage[0], // Overall CPU usage percentage
		TotalMemory:     memoryStats.Total,
		MemoryUsage:     memoryStats.Used,
		MemoryPercent:   memoryStats.UsedPercent,
		TotalDisk:       diskStats.Total,
		DiskUsage:       diskStats.Used,
		DiskPercent:     diskStats.UsedPercent,
		Uptime:          uptime,
	}

	return metrics, nil
}
