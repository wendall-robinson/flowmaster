package kafkacarrier

import (
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/propagation"
)

// KafkaHeadersCarrier is a custom carrier for Kafka headers.
type KafkaHeadersCarrier struct {
	Headers *[]kafka.Header
}

// New creates a new KafkaHeadersCarrier.
func New(headers *[]kafka.Header) *KafkaHeadersCarrier {
	return &KafkaHeadersCarrier{Headers: headers}
}

// Get retrieves the value of the header with the given key.
func (c *KafkaHeadersCarrier) Get(key string) string {
	for _, h := range *c.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}

	return ""
}

// Set sets the value of the header with the given key.
func (c *KafkaHeadersCarrier) Set(key, value string) {
	// Remove existing header with the same key to prevent duplicates
	for i, h := range *c.Headers {
		if h.Key == key {
			*c.Headers = append((*c.Headers)[:i], (*c.Headers)[i+1:]...)
			break
		}
	}

	*c.Headers = append(*c.Headers, kafka.Header{Key: key, Value: []byte(value)})
}

// Keys returns the keys of the headers.
func (c *KafkaHeadersCarrier) Keys() []string {
	keys := make([]string, 0, len(*c.Headers))
	for _, h := range *c.Headers {
		keys = append(keys, h.Key)
	}

	return keys
}

// Ensure KafkaHeadersCarrier implements TextMapCarrier
var _ propagation.TextMapCarrier = (*KafkaHeadersCarrier)(nil)
