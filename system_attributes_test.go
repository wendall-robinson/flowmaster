package traceflow

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

// TestAddSystemInfo verifies that system attributes like hostname, IP, and environment are added.
func TestAddSystemInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")

	trace.AddSystemInfo("test-hostname", "192.168.1.1", "production")

	assert.Len(t, trace.attrs, 3)
	assert.Contains(t, trace.attrs, attribute.String("system.hostname", "test-hostname"))
	assert.Contains(t, trace.attrs, attribute.String("system.ip_address", "192.168.1.1"))
	assert.Contains(t, trace.attrs, attribute.String("system.environment", "production"))
}

// TestAddCpuInfo verifies that CPU-related attributes are added correctly.
func TestAddCpuInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")

	trace.AddCpuInfo()

	assert.Len(t, trace.attrs, 2)
	assert.Equal(t, attribute.Int("cpu.count", runtime.NumCPU()), trace.attrs[0])
	assert.Equal(t, attribute.String("cpu.architecture", runtime.GOARCH), trace.attrs[1])
}

// TestAddMemoryInfo verifies that memory attributes are added correctly.
func TestAddMemoryInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")

	trace.AddMemoryInfo()

	// Verify that 4 memory attributes are added
	assert.Len(t, trace.attrs, 4)
}

// TestAddProcessInfo verifies that process attributes are added correctly.
func TestAddProcessInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")

	trace.AddProcessInfo()

	processID := os.Getpid()
	command, err := exec.LookPath(os.Args[0])
	if err != nil {
		command = "unknown"
	}

	assert.Contains(t, trace.attrs, attribute.Int("process.id", processID))
	assert.Contains(t, trace.attrs, attribute.String("process.command", command))
}
