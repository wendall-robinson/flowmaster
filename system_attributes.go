package traceflow

import (
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

// AddCpuInfo adds CPU-related attributes like CPU count and architecture.
func (t *Trace) AddCpuInfo(cpuCount int, cpuArchitecture string) *Trace {
	t.attrs = append(t.attrs,
		attribute.Int("cpu.count", cpuCount),
		attribute.String("cpu.architecture", cpuArchitecture),
	)
	return t
}

// AddMemoryInfo adds memory-related attributes like total memory and free memory.
func (t *Trace) AddMemoryInfo(totalMemory, freeMemory int64) *Trace {
	t.attrs = append(t.attrs,
		attribute.Int64("memory.total", totalMemory),
		attribute.Int64("memory.free", freeMemory),
	)
	return t
}

// AddDiskInfo adds disk-related attributes like total disk space and free disk space.
func (t *Trace) AddDiskInfo(totalDiskSpace, freeDiskSpace int64) *Trace {
	t.attrs = append(t.attrs,
		attribute.Int64("disk.total", totalDiskSpace),
		attribute.Int64("disk.free", freeDiskSpace),
	)
	return t
}

// AddProcessInfo adds process-related attributes like process ID and command.
func (t *Trace) AddProcessInfo(processID int, command string) *Trace {
	t.attrs = append(t.attrs,
		attribute.Int("process.id", processID),
		attribute.String("process.command", command),
	)
	return t
}

// AddContainerInfo adds container-related attributes like container ID and image.
func (t *Trace) AddContainerInfo(containerID, image string) *Trace {
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
