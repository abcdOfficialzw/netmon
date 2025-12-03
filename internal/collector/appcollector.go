package collector

import (
	"fmt"
)

// AppCollector manages application-level network statistics collection.
type AppCollector struct {
	interfaceCollector *Collector
	connectionMapper   *ConnectionMapper
	lastTotalBytes     uint64
}

// NewAppCollector creates a new application network statistics collector.
func NewAppCollector() *AppCollector {
	return &AppCollector{
		interfaceCollector: NewCollector(),
		connectionMapper:   NewConnectionMapper(),
		lastTotalBytes:     0,
	}
}

// AppDelta represents the change in network traffic for an application.
type AppDelta struct {
	AppName   string
	BytesIn   uint64
	BytesOut  uint64
	Timestamp int64
}

// Collect reads current network stats and distributes traffic among active applications.
// This uses a heuristic approach: traffic is distributed proportionally among apps
// with active network connections.
func (ac *AppCollector) Collect() ([]AppDelta, error) {
	// Update connection mapping
	if err := ac.connectionMapper.Update(); err != nil {
		return nil, fmt.Errorf("update connections: %w", err)
	}

	// Get interface-level deltas
	interfaceDeltas, err := ac.interfaceCollector.Collect()
	if err != nil {
		return nil, fmt.Errorf("collect interfaces: %w", err)
	}

	// First collection - no deltas yet
	if interfaceDeltas == nil {
		return nil, nil
	}

	// Get active apps
	activeApps := ac.connectionMapper.GetActiveApps()
	if len(activeApps) == 0 {
		// No active apps detected
		return []AppDelta{}, nil
	}

	// Sum total traffic across all interfaces
	var totalBytesIn, totalBytesOut uint64
	var timestamp int64

	for _, delta := range interfaceDeltas {
		totalBytesIn += delta.BytesIn
		totalBytesOut += delta.BytesOut
		timestamp = delta.Timestamp
	}

	// If there's no traffic, return empty
	if totalBytesIn == 0 && totalBytesOut == 0 {
		return []AppDelta{}, nil
	}

	// Distribute traffic evenly among active apps
	// Note: This is a heuristic. For accurate per-app tracking, you'd need
	// a kernel extension or network extension with proper entitlements.
	bytesInPerApp := totalBytesIn / uint64(len(activeApps))
	bytesOutPerApp := totalBytesOut / uint64(len(activeApps))

	appDeltas := make([]AppDelta, 0, len(activeApps))
	for _, appName := range activeApps {
		appDeltas = append(appDeltas, AppDelta{
			AppName:   appName,
			BytesIn:   bytesInPerApp,
			BytesOut:  bytesOutPerApp,
			Timestamp: timestamp,
		})
	}

	return appDeltas, nil
}

// CollectWithWeighting collects app traffic using connection-count weighting.
// Apps with more connections get proportionally more traffic attributed.
func (ac *AppCollector) CollectWithWeighting() ([]AppDelta, error) {
	// Update connection mapping
	if err := ac.connectionMapper.Update(); err != nil {
		return nil, fmt.Errorf("update connections: %w", err)
	}

	// Get interface-level deltas
	interfaceDeltas, err := ac.interfaceCollector.Collect()
	if err != nil {
		return nil, fmt.Errorf("collect interfaces: %w", err)
	}

	// First collection - no deltas yet
	if interfaceDeltas == nil {
		return nil, nil
	}

	// Get active processes with connection counts
	snapshot := ac.connectionMapper.lastSnapshot
	if len(snapshot) == 0 {
		return []AppDelta{}, nil
	}

	// Sum total traffic and total connections
	var totalBytesIn, totalBytesOut uint64
	var timestamp int64

	for _, delta := range interfaceDeltas {
		totalBytesIn += delta.BytesIn
		totalBytesOut += delta.BytesOut
		timestamp = delta.Timestamp
	}

	// If no traffic, return empty
	if totalBytesIn == 0 && totalBytesOut == 0 {
		return []AppDelta{}, nil
	}

	// Calculate total connections and aggregate by app
	appConnections := make(map[string]int)
	for _, procInfo := range snapshot {
		appConnections[procInfo.AppName] += procInfo.Connections
	}

	totalConnections := 0
	for _, count := range appConnections {
		totalConnections += count
	}

	if totalConnections == 0 {
		return []AppDelta{}, nil
	}

	// Distribute traffic proportionally based on connection count
	appDeltas := make([]AppDelta, 0, len(appConnections))

	for appName, connections := range appConnections {
		weight := float64(connections) / float64(totalConnections)
		bytesIn := uint64(float64(totalBytesIn) * weight)
		bytesOut := uint64(float64(totalBytesOut) * weight)

		appDeltas = append(appDeltas, AppDelta{
			AppName:   appName,
			BytesIn:   bytesIn,
			BytesOut:  bytesOut,
			Timestamp: timestamp,
		})
	}

	return appDeltas, nil
}

