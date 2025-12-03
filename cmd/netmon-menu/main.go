package main

import (
	"flag"
	"fmt"
	"netmon/internal/db"
	"netmon/internal/stats"
	"os"
	"path/filepath"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	var dbPath string
	flag.StringVar(&dbPath, "db", getDefaultDBPath(), "Path to SQLite database file")
	flag.Parse()

	// Open database
	database, err := db.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Run the menu bar app
	systray.Run(func() {
		onReady(database)
	}, onExit)
}

func onReady(database *db.DB) {
	// Set initial title and tooltip
	systray.SetTitle("NetMon")
	systray.SetTooltip("Network Usage Monitor")

	// Create menu items
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Start a goroutine to update the menu bar title periodically
	go func() {
		ticker := time.NewTicker(5 * time.Second) // Update every 5 seconds
		defer ticker.Stop()

		updateTitle(database)

		for {
			select {
			case <-ticker.C:
				updateTitle(database)
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func updateTitle(database *db.DB) {
	startTime := db.GetStartOfDay()
	endTime := time.Now().Unix()

	logs, err := database.GetLogsInRange(startTime, endTime)
	if err != nil {
		systray.SetTitle("NetMon: Error")
		systray.SetTooltip(fmt.Sprintf("Error: %v", err))
		return
	}

	if len(logs) == 0 {
		systray.SetTitle("NetMon: 0 B")
		systray.SetTooltip("No data available for today")
		return
	}

	summary := stats.ComputeSummary(logs)
	totalBytes := summary.TotalBytesIn + summary.TotalBytesOut

	// Format for menu bar (keep it short)
	formatted := formatBytesShort(totalBytes)
	systray.SetTitle(fmt.Sprintf("🌐 %s", formatted))
	
	// Detailed tooltip
	tooltip := fmt.Sprintf("Today's Network Usage\n\nDownloaded: %s\nUploaded: %s\nTotal: %s",
		stats.FormatBytes(summary.TotalBytesIn),
		stats.FormatBytes(summary.TotalBytesOut),
		stats.FormatBytes(totalBytes))
	systray.SetTooltip(tooltip)
}

func formatBytesShort(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.1fT", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.1fG", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1fM", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1fK", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%dB", bytes)
	}
}

func onExit() {
	// Cleanup code if needed
}

func getDefaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./netmon.db"
	}
	return filepath.Join(home, ".netmon", "netmon.db")
}

