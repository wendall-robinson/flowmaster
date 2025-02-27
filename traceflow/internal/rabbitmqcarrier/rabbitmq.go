package rabbitmqcarrier

import (
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/propagation"
)

// AMQPTableCarrier is a TextMapCarrier for AMQP headers.
type AMQPTableCarrier struct {
	Headers amqp.Table
}

// New creates a new AMQPTableCarrier.
func New(headers amqp.Table) *AMQPTableCarrier {
	return &AMQPTableCarrier{
		Headers: headers,
	}
}

// Get returns the value associated with the key.
func (c *AMQPTableCarrier) Get(key string) string {
	if val, ok := c.Headers[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}

	return ""
}

// Set stores the key-value pair.
func (c *AMQPTableCarrier) Set(key, value string) {
	c.Headers[key] = value
}

// Keys returns the keys of the AMQPTableCarrier.
func (c *AMQPTableCarrier) Keys() []string {
	keys := make([]string, 0, len(c.Headers))
	for k := range c.Headers {
		keys = append(keys, k)
	}

	return keys
}

// Ensure AMQPTableCarrier implements TextMapCarrier
var _ propagation.TextMapCarrier = (*AMQPTableCarrier)(nil)
