package logflow

import (
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

// SystemAttributes represents system attributes for a log entry
type SystemAttributes struct {
	OperatingSystem  string           `json:"operating_system,omitempty"`
	MemoryUsage      Memory           `json:"memory_usage,omitempty"`
	DiskUsage        Disk             `json:"disk_usage,omitempty"`
	CPUUsagePercent  CPU              `json:"cpu_usage_percent,omitempty"`
	RunningProcesses RunningProcesses `json:"running_processes,omitempty"`
}

// Memory represents memory usage statistics
type Memory struct {
	TotalMemory  uint64 `json:"total_memory,omitempty"`
	UsedMemory   uint64 `json:"used_memory,omitempty"`
	FreeMemory   uint64 `json:"free_memory,omitempty"`
	HeapAlloc    uint64 `json:"heap_alloc,omitempty"`
	HeapSys      uint64 `json:"heap_sys,omitempty"`
	StackInUse   uint64 `json:"stack_in_use,omitempty"`
	StackSys     uint64 `json:"stack_sys,omitempty"`
	GCCount      uint32 `json:"gc_count,omitempty"`
	GCPauseTotal uint64 `json:"gc_pause_total,omitempty"`
}

// Disk represents disk usage statistics
type Disk struct {
	Path             string  `json:"path,omitempty"`
	TotalDiskSpace   uint64  `json:"total_disk_space,omitempty"`
	UsedDiskSpace    uint64  `json:"used_disk_space,omitempty"`
	FreeDiskSpace    uint64  `json:"free_disk_space,omitempty"`
	DiskUsagePercent float64 `json:"disk_usage_percent,omitempty"`
}

// CPU represents CPU usage statistics
type CPU struct {
	UsagePercent interface{} `json:"usage_percent,omitempty"`
}

// RunningProcesses represents a list of running processes
type RunningProcesses struct {
	Processes interface{} `json:"processes,omitempty"`
}

// findSystemAttributes will apply system attributes to the log entry
func (l *Logger) findSystemAttributes(entry *LogEntry) *LogEntry {
	attributes := SystemAttributes{
		OperatingSystem:  l.getOperatingSystem(),
		MemoryUsage:      l.getMemoryUsage(),
		DiskUsage:        l.getDiskUsage("/"),
		CPUUsagePercent:  l.getCPUUsage(),
		RunningProcesses: l.getRunningProcesses(),
	}

	entry.Context["system_attributes"] = attributes

	return entry
}

// systemAttribuesExists checks if the system attributes field exists in the context.
//
//   - If it exists, it returns the existing SystemAttributes pointer and true.
//   - If it doesn't it
func systemAttribuesExists(context map[string]interface{}) (*SystemAttributes, bool) {
	// Check if "system_attributes" already exists
	if attributes, ok := context["system_attributes"].(*SystemAttributes); ok {
		return attributes, true
	}

	return &SystemAttributes{}, false
}

// systemAttribuesFieldExists checks if a specific field exists in the SystemAttributes struct.
func systemAttribuesFieldExists(attributes *SystemAttributes, field string) bool {
	// check if field is in attributes
	switch field {
	case "OperatingSystem":
		return attributes.OperatingSystem != ""
	case "MemoryUsage":
		return (attributes.MemoryUsage != Memory{})
	case "DiskUsage":
		return (attributes.DiskUsage != Disk{})
	case "CPUUsagePercent":
		return (attributes.CPUUsagePercent != CPU{})
	case "RunningProcesses":
		return (attributes.RunningProcesses != RunningProcesses{})
	default:
		return false
	}
}

// getOperatingSystem returns the operating system name
func (l *Logger) getOperatingSystem() string {
	return runtime.GOOS
}

// getMemoryUsage returns memory usage statistics
func (l *Logger) getMemoryUsage() Memory {
	vmStats, err := mem.VirtualMemory()

	var memory Memory
	if err == nil {
		memory.TotalMemory = vmStats.Total
		memory.UsedMemory = vmStats.Used
		memory.FreeMemory = vmStats.Free
	} else {
		memory.TotalMemory = 0
		memory.UsedMemory = 0
		memory.FreeMemory = 0
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memory.HeapAlloc = memStats.HeapAlloc
	memory.HeapSys = memStats.HeapSys
	memory.StackInUse = memStats.StackInuse
	memory.StackSys = memStats.StackSys
	memory.GCCount = memStats.NumGC
	memory.GCPauseTotal = memStats.PauseTotalNs

	return memory
}

// getDiskStats returns disk usage statistics for the specified path
func (l *Logger) getDiskUsage(path string) Disk {
	var diskUsage Disk
	diskUsage.Path = path

	usageStats, err := disk.Usage(path)
	if err == nil {
		diskUsage.TotalDiskSpace = usageStats.Total
		diskUsage.UsedDiskSpace = usageStats.Used
		diskUsage.FreeDiskSpace = usageStats.Free
		diskUsage.DiskUsagePercent = usageStats.UsedPercent
	}

	return diskUsage
}

// getCPUUsage returns the overall CPU usage percentage
func (l *Logger) getCPUUsage() CPU {
	cpuUsage := CPU{}
	cpuPercentages, err := cpu.Percent(0, false) // Overall CPU usage

	if err == nil && len(cpuPercentages) > 0 {
		cpuUsage.UsagePercent = cpuPercentages[0]
	} else {
		cpuUsage.UsagePercent = nil
	}

	return cpuUsage
}

// getRunningProcesses returns a list of running processes
func (l *Logger) getRunningProcesses() RunningProcesses {
	processes, err := process.Processes()
	runningProcesses := RunningProcesses{}

	if err != nil {
		runningProcesses.Processes = nil

		return runningProcesses
	}

	// Simplify output to reduce log size
	processList := []string{}
	for _, p := range processes {
		name, err := p.Name()
		if err == nil {
			processList = append(processList, name)
		}

		if len(processList) >= 10 { // Limit to top 10 processes
			break
		}
	}

	runningProcesses.Processes = processList

	return runningProcesses
}
