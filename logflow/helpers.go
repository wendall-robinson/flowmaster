package logflow

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
)

// Helper function to handle logging errors
func processLog(entry *LogEntry) error {
	// Example: Simulate sending the log entry to stdout or a sink
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	fmt.Println(string(data)) // Output the log entry
	return nil
}

// Helper function to generate a unique LogID
func generateLogID() string {
	// Generate a random LogID (replace with a more robust solution if needed)
	return fmt.Sprintf("%d", rand.Int63())
}

// Helper function to handle logging errors
func handleLoggingError(err error) {
	// Log the error to a fallback (e.g., stderr)
	fmt.Fprintf(os.Stderr, "Logging error: %v\n", err)
}

// getCaller identifies the calling module or file
func getCaller() string {
	_, file, _, ok := runtime.Caller(2) // Adjust depth as needed
	if !ok {
		return "unknown"
	}

	// Extract the file name and optional package/module
	parts := strings.Split(file, "/")
	return parts[len(parts)-1]
}

// addKVPairsToContext processes variadic string key/value pairs and adds them to the context
func addKVPairsToContext(entry *LogEntry, keyValuePairs []string) {
	for i := 0; i < len(keyValuePairs); i += 2 {
		key := keyValuePairs[i]

		// Check if there's a matching value
		if i+1 < len(keyValuePairs) {
			value := keyValuePairs[i+1]
			entry.Context[key] = value
		} else {
			// If no matching value, add key with an empty string
			entry.Context[key] = ""
		}
	}
}
