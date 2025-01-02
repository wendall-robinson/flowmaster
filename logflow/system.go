package logflow

import (
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

type SystemAttributes struct {
	OperatingSystem  string           `json:"operating_system,omitempty"`
	MemoryUsage      Memory           `json:"memory_usage,omitempty"`
	DiskUsage        Disk             `json:"disk_usage,omitempty"`
	CPUUsagePercent  CPU              `json:"cpu_usage_percent,omitempty"`
	RunningProcesses RunningProcesses `json:"running_processes,omitempty"`
}

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

type Disk struct {
	Path             string  `json:"path,omitempty"`
	TotalDiskSpace   uint64  `json:"total_disk_space,omitempty"`
	UsedDiskSpace    uint64  `json:"used_disk_space,omitempty"`
	FreeDiskSpace    uint64  `json:"free_disk_space,omitempty"`
	DiskUsagePercent float64 `json:"disk_usage_percent,omitempty"`
}

type CPU struct {
	UsagePercent interface{} `json:"usage_percent,omitempty"`
}

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
