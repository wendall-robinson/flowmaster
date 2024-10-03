package traceflow

import (
	"math"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// AddSystemInfo adds system-related attributes like hostname and environment.
func (t *Trace) AddSystemInfo(hostname, ipAddress, environment string) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("system.hostname", hostname),
		attribute.String("system.ip_address", ipAddress),
		attribute.String("system.environment", environment),
	)

	return t
}

// AddCPUInfo adds CPU count and CPU architecture attributes to the trace.
// This information is gathered automatically using Go's runtime package.
//
// Example usage:
//
//	trace.AddCpuInfo()
func (t *Trace) AddCPUInfo() *Trace {
	cpuCount := runtime.NumCPU()
	cpuArchitecture := runtime.GOARCH

	t.attrs = append(t.attrs,
		attribute.Int("cpu.count", cpuCount),
		attribute.String("cpu.architecture", cpuArchitecture),
	)

	return t
}

// AddMemoryInfo automatically adds memory-related attributes to the trace.
// This includes details such as total memory allocation, system memory, and heap memory.
// The memory information is retrieved using Go's runtime package.
//
// Attributes added:
// - memory.total_alloc: Total bytes allocated.
// - memory.sys: System memory in use.
// - memory.heap_alloc: Heap memory allocated.
// - memory.heap_idle: Heap memory currently idle.
//
// Example usage:
//
//	trace.AddMemoryInfo()
//
// Notes:
// - The memory statistics are collected automatically, and no manual input is required.
func (t *Trace) AddMemoryInfo() *Trace {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)
	t.attrs = append(t.attrs,
		attribute.Int64("memory.total_alloc", safeUint64ToInt64(memStats.TotalAlloc)),
		attribute.Int64("memory.sys", safeUint64ToInt64(memStats.Sys)),
		attribute.Int64("memory.heap_alloc", safeUint64ToInt64(memStats.HeapAlloc)),
		attribute.Int64("memory.heap_idle", safeUint64ToInt64(memStats.HeapIdle)),
	)

	return t
}

// AddDiskInfo automatically adds disk-related attributes to the trace.
// This includes details such as total disk space and free disk space. The information
// is retrieved using system calls to gather disk statistics.
//
// Attributes added:
// - disk.total: Total disk space in bytes.
// - disk.free: Free disk space in bytes.
//
// Example usage:
//
//	trace.AddDiskInfo()
//
// Notes:
//   - Disk statistics are collected automatically for the root filesystem ("/").
//   - This implementation uses syscall for Unix-like systems. Adjustments may be
//     required for other operating systems.
func (t *Trace) AddDiskInfo() *Trace {
	var stat syscall.Statfs_t

	err := syscall.Statfs("/", &stat)
	if err != nil {
		return t
	}

	totalDiskSpace := stat.Blocks * uint64(stat.Bsize)
	freeDiskSpace := stat.Bfree * uint64(stat.Bsize)

	t.attrs = append(t.attrs,
		attribute.Int64("disk.total", safeUint64ToInt64(totalDiskSpace)),
		attribute.Int64("disk.free", safeUint64ToInt64(freeDiskSpace)),
	)

	return t
}

// AddProcessInfo automatically adds process-related attributes such as the process ID
// and the command being executed to the trace. The process ID is retrieved using Go's
// os package, and the command is determined from the current executable.
//
// Attributes added:
// - process.id: The current process ID.
// - process.command: The command or path of the executable.
//
// Example usage:
//
//	trace.AddProcessInfo()
//
// Notes:
// - If the command cannot be determined, it defaults to "unknown".
func (t *Trace) AddProcessInfo() *Trace {
	processID := os.Getpid() // Gets the current process ID

	command, err := exec.LookPath(os.Args[0]) // Gets the command (path) being executed
	if err != nil {
		command = "unknown"
	}

	t.attrs = append(t.attrs,
		attribute.Int("process.id", processID),
		attribute.String("process.command", command),
	)

	return t
}

// AddContainerInfo automatically adds container-related attributes such as the container ID
// and the image to the trace. The container ID is retrieved from the cgroup file, and the
// container image is fetched from the CONTAINER_IMAGE environment variable.
//
// Attributes added:
// - container.id: The container ID, retrieved from the cgroup file.
// - container.image: The container image, retrieved from the environment or set to "unknown".
//
// Example usage:
//
//	trace.AddContainerInfo()
//
// Notes:
// - This method is designed to work in Docker or Kubernetes environments.
func (t *Trace) AddContainerInfo() *Trace {
	// Get the container ID from the cgroup (works in Docker/Kubernetes)
	containerID := "unknown"

	if data, err := os.ReadFile("/proc/self/cgroup"); err == nil {
		// Extract the container ID from the cgroup file
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.Contains(line, "docker") || strings.Contains(line, "kubepods") {
				parts := strings.Split(line, "/")
				if len(parts) > 1 {
					containerID = parts[len(parts)-1]
					break
				}
			}
		}
	}

	// The container image might be passed as an environment variable
	image := os.Getenv("CONTAINER_IMAGE")
	if image == "" {
		image = "unknown"
	}

	t.attrs = append(t.attrs,
		attribute.String("container.id", containerID),
		attribute.String("container.image", image),
	)

	return t
}

// AddKubernetesInfo adds Kubernetes-related attributes like pod name and namespace.
func (t *Trace) AddKubernetesInfo(podName, namespace string) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("kubernetes.pod_name", podName),
		attribute.String("kubernetes.namespace", namespace),
	)

	return t
}

// AddNetworkInfo adds network-related attributes to the trace.
func (t *Trace) AddNetworkInfo(protocol string, latency time.Duration) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("network.protocol", protocol),
		attribute.Int64("network.latency_ms", latency.Milliseconds()),
	)

	return t
}

// safeUint64ToInt64 safely converts a uint64 value to an int64, ensuring that no overflow occurs.
// Go's OpenTelemetry attributes only support int64 for numeric values, so this conversion is
// necessary when dealing with system values that are represented as uint64 (such as memory and
// disk sizes). If the uint64 value exceeds the maximum value that an int64 can hold, this function
// returns math.MaxInt64 to prevent overflow.
//
// This approach ensures compatibility with OpenTelemetry while avoiding errors that could occur
// from direct uint64 to int64 conversion.
//
// Example:
//
//	value := safeUint64ToInt64(9876543210)
//	// This will convert the uint64 value safely, or return MaxInt64 if it's too large.
func safeUint64ToInt64(u uint64) int64 {
	if u > math.MaxInt64 {
		return math.MaxInt64 // Handle overflow case (use max int64 value)
	}

	return int64(u)
}
