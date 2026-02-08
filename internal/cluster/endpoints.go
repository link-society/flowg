package cluster

import (
	"iter"
	"net/url"
	"sync"
)

type endpointCache struct {
	index map[string]*url.URL
	mu    sync.RWMutex
}

func (c *endpointCache) Get(nodeID string) (*url.URL, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	endpoint, ok := c.index[nodeID]
	return endpoint, ok
}

func (c *endpointCache) Set(nodeID string, endpoint *url.URL) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.index[nodeID] = endpoint
}

func (c *endpointCache) Delete(nodeID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.index, nodeID)
}

func (c *endpointCache) All() iter.Seq2[string, *url.URL] {
	return func(yield func(string, *url.URL) bool) {
		c.mu.RLock()
		defer c.mu.RUnlock()

		for nodeID, endpoint := range c.index {
			if !yield(nodeID, endpoint) {
				return
			}
		}
	}
}

func newEndpointCache() *endpointCache {
	return &endpointCache{
		index: make(map[string]*url.URL),
	}
}
