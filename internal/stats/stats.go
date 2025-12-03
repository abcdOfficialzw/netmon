package stats

import (
	"fmt"
	"netmon/internal/db"
)

// Summary represents aggregated network traffic statistics.
type Summary struct {
	TotalBytesIn  uint64
	TotalBytesOut uint64
	PeakBytesIn   uint64 // bytes per second
	PeakBytesOut  uint64 // bytes per second
}

// ComputeSummary calculates traffic summary from traffic logs.
func ComputeSummary(logs []db.TrafficLog) Summary {
	var s Summary

	for _, log := range logs {
		s.TotalBytesIn += log.BytesIn
		s.TotalBytesOut += log.BytesOut

		if log.BytesIn > s.PeakBytesIn {
			s.PeakBytesIn = log.BytesIn
		}

		if log.BytesOut > s.PeakBytesOut {
			s.PeakBytesOut = log.BytesOut
		}
	}

	return s
}

// InterfaceSummary represents traffic summary for a single interface.
type InterfaceSummary struct {
	Interface     string
	TotalBytesIn  uint64
	TotalBytesOut uint64
}

// ComputeByInterface calculates traffic summary per interface.
func ComputeByInterface(logsByInterface map[string][]db.TrafficLog) []InterfaceSummary {
	summaries := make([]InterfaceSummary, 0, len(logsByInterface))

	for iface, logs := range logsByInterface {
		var totalIn, totalOut uint64
		for _, log := range logs {
			totalIn += log.BytesIn
			totalOut += log.BytesOut
		}

		summaries = append(summaries, InterfaceSummary{
			Interface:     iface,
			TotalBytesIn:  totalIn,
			TotalBytesOut: totalOut,
		})
	}

	return summaries
}

// FormatBytes converts bytes to a human-readable string.
func FormatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// FormatBytesPerSec converts bytes per second to a human-readable throughput string.
func FormatBytesPerSec(bytesPerSec uint64) string {
	return FormatBytes(bytesPerSec) + "/s"
}

// AppSummary represents traffic summary for a single application.
type AppSummary struct {
	AppName       string
	TotalBytesIn  uint64
	TotalBytesOut uint64
}

// ComputeByApp calculates traffic summary per application.
func ComputeByApp(logsByApp map[string][]db.AppTrafficLog) []AppSummary {
	summaries := make([]AppSummary, 0, len(logsByApp))

	for app, logs := range logsByApp {
		var totalIn, totalOut uint64
		for _, log := range logs {
			totalIn += log.BytesIn
			totalOut += log.BytesOut
		}

		summaries = append(summaries, AppSummary{
			AppName:       app,
			TotalBytesIn:  totalIn,
			TotalBytesOut: totalOut,
		})
	}

	return summaries
}

