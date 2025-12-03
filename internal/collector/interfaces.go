package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/net"
)

// InterfaceStats represents network interface statistics at a point in time.
type InterfaceStats struct {
	Name     string
	BytesIn  uint64
	BytesOut uint64
}

// ReadInterfaces reads current network interface statistics from the system.
func ReadInterfaces() ([]InterfaceStats, error) {
	ioCounters, err := net.IOCounters(true) // true = per interface
	if err != nil {
		return nil, fmt.Errorf("read io counters: %w", err)
	}

	stats := make([]InterfaceStats, 0, len(ioCounters))
	for _, counter := range ioCounters {
		// Skip interfaces with no traffic
		if counter.BytesRecv == 0 && counter.BytesSent == 0 {
			continue
		}

		stats = append(stats, InterfaceStats{
			Name:     counter.Name,
			BytesIn:  counter.BytesRecv,
			BytesOut: counter.BytesSent,
		})
	}

	return stats, nil
}

