package natscarrier

import "github.com/nats-io/nats.go"

type NATSHeaderCarrier struct {
	headers nats.Header
}

func (c *NATSHeaderCarrier) Get(key string) string {
	return c.headers.Get(key)
}

func (c *NATSHeaderCarrier) Set(key, value string) {
	c.headers.Set(key, value)
}

// Keys returns the keys of the headers.
func (c *NATSHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c.headers))
	for key := range c.headers {
		keys = append(keys, key)
	}

	return keys
}
