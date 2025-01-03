package logflow

// Option defines a function that configures the Logger
type Option func(*LoggerConfig)

// WithLogLevel sets the log level
func WithLogLevel(level string) Option {
	return func(config *LoggerConfig) {
		config.LogLevel = level
	}
}

// WithStdout enables or disables logging to stdout
func WithStdout(enabled bool) Option {
	return func(config *LoggerConfig) {
		config.OutputToStdout = enabled
	}
}

// WithTimestampFormat sets the timestamp format
func WithTimestampFormat(format string) Option {
	return func(config *LoggerConfig) {
		config.TimestampFmt = format
	}
}

// WithLogCollector sets a custom log collector
func WithLogCollector(collector LogCollector) Option {
	return func(config *LoggerConfig) {
		config.LogCollector = collector
	}
}

// WithEnvHostname configures the Logger to include the hostname in log entries
func WithEnvHostname() Option {
	return func(config *LoggerConfig) {
		config.IncludeHostname = true
	}
}

// WithEnvironment configures the Logger to include the environment in log entries
func WithEnvironment() Option {
	return func(config *LoggerConfig) {
		config.IncludeEnvironment = true
	}
}

// WithK8sMetadata adds Kubernetes metadata to the Logger
func WithK8sMetadata(metadata map[string]string) Option {
	return func(config *LoggerConfig) {
		config.IncludeK8s = true
		config.K8sConfig = metadata
	}
}

// WithYangConfig adds YANG-derived configuration to the Logger
func WithYangConfig(yangData map[string]string) Option {
	return func(config *LoggerConfig) {
		config.IncludeYANG = true
		config.YANGConfig = yangData
	}
}

// WithSystemInfo informs that all logs should include system metadata
func WithSystemInfo() Option {
	return func(config *LoggerConfig) {
		config.IncludeSystem = true
	}
}
