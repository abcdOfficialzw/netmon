package collector

import (
	"fmt"
	"time"
)

// Collector manages network statistics collection and delta computation.
type Collector struct {
	lastStats map[string]InterfaceStats
}

// NewCollector creates a new network statistics collector.
func NewCollector() *Collector {
	return &Collector{
		lastStats: make(map[string]InterfaceStats),
	}
}

// Delta represents the change in network traffic over a time period.
type Delta struct {
	Interface string
	BytesIn   uint64
	BytesOut  uint64
	Timestamp int64
}

// Collect reads current interface stats and computes deltas since last collection.
// On the first call, it initializes state and returns nil (no delta yet).
func (c *Collector) Collect() ([]Delta, error) {
	currentStats, err := ReadInterfaces()
	if err != nil {
		return nil, fmt.Errorf("read interfaces: %w", err)
	}

	// First collection - just store state
	if len(c.lastStats) == 0 {
		for _, stat := range currentStats {
			c.lastStats[stat.Name] = stat
		}
		return nil, nil
	}

	// Compute deltas
	now := time.Now().Unix()
	var deltas []Delta

	for _, current := range currentStats {
		last, exists := c.lastStats[current.Name]
		if !exists {
			// New interface appeared
			c.lastStats[current.Name] = current
			continue
		}

		// Compute delta (handle counter wraparound)
		var deltaIn, deltaOut uint64

		if current.BytesIn >= last.BytesIn {
			deltaIn = current.BytesIn - last.BytesIn
		} else {
			// Counter wrapped around
			deltaIn = current.BytesIn
		}

		if current.BytesOut >= last.BytesOut {
			deltaOut = current.BytesOut - last.BytesOut
		} else {
			// Counter wrapped around
			deltaOut = current.BytesOut
		}

		deltas = append(deltas, Delta{
			Interface: current.Name,
			BytesIn:   deltaIn,
			BytesOut:  deltaOut,
			Timestamp: now,
		})

		// Update state
		c.lastStats[current.Name] = current
	}

	return deltas, nil
}

