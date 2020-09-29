package cmd

import (
	"sync"
)

// countMap is a map of counters
type countMap struct {
	m sync.Mutex
	c map[string]int
}

// add adds count to the counter identified by name
func (c *countMap) add(name string, count int) {
	c.m.Lock()
	defer c.m.Unlock()

	c.c[name] += count
}

// get returns the current count of the counter identified by name
func (c *countMap) get(name string) int {
	c.m.Lock()
	defer c.m.Unlock()

	return c.c[name]
}

// getAll returns a map of all counters in the countMap
func (c *countMap) getAll() map[string]int {
	c.m.Lock()
	defer c.m.Unlock()

	m := make(map[string]int)
	for k, v := range c.c {
		m[k] = v
	}
	return m
}

// newCountMap creates a new countMap
func newCountMap() *countMap {
	m := countMap{
		c: make(map[string]int),
	}
	return &m
}
