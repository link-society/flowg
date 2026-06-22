package cluster

import "sync"

type watermarkCache struct {
	mu sync.Mutex
	m  map[string]map[string]uint64
}

func newWatermarkCache() *watermarkCache {
	return &watermarkCache{
		m: make(map[string]map[string]uint64),
	}
}

func (c *watermarkCache) observe(nodeID, namespace string, since uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ns := c.m[nodeID]
	if ns == nil {
		ns = make(map[string]uint64)
		c.m[nodeID] = ns
	}
	if since > ns[namespace] {
		ns[namespace] = since
	}
}

func (c *watermarkCache) get(nodeID, namespace string) uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.m[nodeID][namespace]
}

func (c *watermarkCache) forget(nodeID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.m, nodeID)
}
