package rabbitmqcarrier

import (
	"github.com/streadway/amqp"
)

// RabbitMQHeaderCarrier adapts amqp.Table to satisfy the TextMapCarrier interface.
type RabbitMQHeaderCarrier struct {
	headers amqp.Table
}

// Get retrieves the value of the header with the given key.
func (c *RabbitMQHeaderCarrier) Get(key string) string {
	if val, ok := c.headers[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}

	return ""
}

// Set sets the value of the header with the given key.
func (c *RabbitMQHeaderCarrier) Set(key, value string) {
	c.headers[key] = value
}

// Keys returns the keys of the headers.
func (c *RabbitMQHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.headers))
	for key := range c.headers {
		keys = append(keys, key)
	}

	return keys
}
