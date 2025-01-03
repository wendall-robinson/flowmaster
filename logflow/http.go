package logflow

import "net/http"

// HTTP represents HTTP request/response metadata
type HTTP struct {
	Method  string            `json:"http_method,omitempty"`      // HTTP method (GET, POST, etc.)
	URL     string            `json:"http_url,omitempty"`         // URL of the request
	Headers map[string]string `json:"http_headers,omitempty"`     // HTTP headers
	Host    string            `json:"http_host,omitempty"`        // Hostname
	Status  int               `json:"http_status_code,omitempty"` // HTTP status code
	Remote  string            `json:"http_remote_addr,omitempty"` // Remote address
}

// extractHTTPMetadata extracts HTTP metadata from an HTTP request
func extractHTTPMetadata(req *http.Request, statusCode int) HTTP {
	if req == nil {
		return HTTP{}
	}

	return HTTP{
		Method:  req.Method,
		URL:     req.URL.String(),
		Headers: sanitizeHeaders(req.Header),
		Host:    req.Host,
		Status:  statusCode,
		Remote:  req.RemoteAddr,
	}
}

// sanitizeHeaders removes sensitive information from headers
func sanitizeHeaders(headers http.Header) map[string]string {
	sanitized := make(map[string]string)
	for key, values := range headers {
		if key == "Authorization" || key == "Cookie" {
			sanitized[key] = "[REDACTED]" // Redact sensitive headers
		} else {
			sanitized[key] = values[0] // Use the first header value
		}
	}

	return sanitized
}
