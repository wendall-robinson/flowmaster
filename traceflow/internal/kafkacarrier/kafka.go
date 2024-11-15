package kafkacarrier

import (
	"github.com/segmentio/kafka-go"
)

// KafkaHeaderCarrier adapts kafka.Header to satisfy the TextMapCarrier interface.
type KafkaHeaderCarrier struct {
	headers []kafka.Header
}

// Get retrieves the value of the header with the given key.
func (c *KafkaHeaderCarrier) Get(key string) string {
	for _, h := range c.headers {
		if h.Key == key {
			return string(h.Value)
		}
	}

	return ""
}

// Set sets the value of the header with the given key.
func (c *KafkaHeaderCarrier) Set(key, value string) {
	// Remove existing header with the same key to prevent duplicates
	for i, h := range c.headers {
		if h.Key == key {
			c.headers = append(c.headers[:i], c.headers[i+1:]...)
			break
		}
	}

	c.headers = append(c.headers, kafka.Header{Key: key, Value: []byte(value)})
}

// Keys returns the keys of the headers.
func (c *KafkaHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.headers))
	for _, h := range c.headers {
		keys = append(keys, h.Key)
	}

	return keys
}

// To use the modified headers after injection/extraction
func (c *KafkaHeaderCarrier) Headers() []kafka.Header {
	return c.headers
}
