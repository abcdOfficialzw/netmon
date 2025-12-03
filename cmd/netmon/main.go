package main

import (
	"bufio"
	"flag"
	"fmt"
	"netmon/internal/db"
	"netmon/internal/stats"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Version information
const (
	Version = "0.2.0"
	AppName = "netmon"
)

func main() {
	// If no arguments, show apps by default
	if len(os.Args) < 2 {
		// Parse flags and show default view
		fs := flag.NewFlagSet("netmon", flag.ExitOnError)
		dbPath := fs.String("db", getDefaultDBPath(), "Path to SQLite database file")
		fs.Parse(os.Args[1:])

		database, err := db.Open(*dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
			os.Exit(1)
		}
		defer database.Close()

		showStatsApps(database)
		return
	}

	command := os.Args[1]

	// Parse global flags
	fs := flag.NewFlagSet("netmon", flag.ExitOnError)
	dbPath := fs.String("db", getDefaultDBPath(), "Path to SQLite database file")

	// Skip the command name when parsing flags
	fs.Parse(os.Args[2:])

	// Open database
	database, err := db.Open(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Execute command
	switch command {
	case "setup":
		handleSetup()
		return
	case "version":
		showVersion()
		return
	case "stats":
		// If "stats" with no subcommand, default to apps
		if fs.NArg() < 1 {
			showStatsApps(database)
			return
		}
		handleStats(database, fs.Arg(0))
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleStats(database *db.DB, subcommand string) {
	switch subcommand {
	case "today":
		showStatsToday(database)
	case "week":
		showStatsWeek(database)
	case "month":
		showStatsMonth(database)
	case "all":
		showStatsAll(database)
	case "interfaces":
		showStatsInterfaces(database)
	case "apps":
		showStatsApps(database)
	default:
		fmt.Fprintf(os.Stderr, "Unknown stats subcommand: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

func showStatsToday(database *db.DB) {
	startTime := db.GetStartOfDay()
	endTime := time.Now().Unix()

	logs, err := database.GetLogsInRange(startTime, endTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching logs: %v\n", err)
		os.Exit(1)
	}

	if len(logs) == 0 {
		fmt.Println("No data available for today")
		return
	}

	summary := stats.ComputeSummary(logs)

	fmt.Println("Stats for today")
	fmt.Println()
	fmt.Println("Overall Totals:")
	fmt.Printf("  Downloaded: %s\n", stats.FormatBytes(summary.TotalBytesIn))
	fmt.Printf("  Uploaded:   %s\n", stats.FormatBytes(summary.TotalBytesOut))
	fmt.Printf("  Total:      %s\n", stats.FormatBytes(summary.TotalBytesIn+summary.TotalBytesOut))
	fmt.Println()
	fmt.Printf("Peak Down:  %s\n", stats.FormatBytesPerSec(summary.PeakBytesIn))
	fmt.Printf("Peak Up:    %s\n", stats.FormatBytesPerSec(summary.PeakBytesOut))
}

func showStatsWeek(database *db.DB) {
	startTime := db.GetStartOfWeek()
	endTime := time.Now().Unix()

	logs, err := database.GetLogsInRange(startTime, endTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching logs: %v\n", err)
		os.Exit(1)
	}

	if len(logs) == 0 {
		fmt.Println("No data available for this week")
		return
	}

	summary := stats.ComputeSummary(logs)

	fmt.Println("Stats for this week (Monday - now)")
	fmt.Println()
	fmt.Println("Overall Totals:")
	fmt.Printf("  Downloaded: %s\n", stats.FormatBytes(summary.TotalBytesIn))
	fmt.Printf("  Uploaded:   %s\n", stats.FormatBytes(summary.TotalBytesOut))
	fmt.Printf("  Total:      %s\n", stats.FormatBytes(summary.TotalBytesIn+summary.TotalBytesOut))
	fmt.Println()
	fmt.Printf("Peak Down:  %s\n", stats.FormatBytesPerSec(summary.PeakBytesIn))
	fmt.Printf("Peak Up:    %s\n", stats.FormatBytesPerSec(summary.PeakBytesOut))
}

func showStatsMonth(database *db.DB) {
	startTime := db.GetStartOfMonth()
	endTime := time.Now().Unix()

	logs, err := database.GetLogsInRange(startTime, endTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching logs: %v\n", err)
		os.Exit(1)
	}

	if len(logs) == 0 {
		fmt.Println("No data available for this month")
		return
	}

	summary := stats.ComputeSummary(logs)

	now := time.Now()
	fmt.Printf("Stats for %s\n", now.Format("January 2006"))
	fmt.Println()
	fmt.Println("Overall Totals:")
	fmt.Printf("  Downloaded: %s\n", stats.FormatBytes(summary.TotalBytesIn))
	fmt.Printf("  Uploaded:   %s\n", stats.FormatBytes(summary.TotalBytesOut))
	fmt.Printf("  Total:      %s\n", stats.FormatBytes(summary.TotalBytesIn+summary.TotalBytesOut))
	fmt.Println()
	fmt.Printf("Peak Down:  %s\n", stats.FormatBytesPerSec(summary.PeakBytesIn))
	fmt.Printf("Peak Up:    %s\n", stats.FormatBytesPerSec(summary.PeakBytesOut))
}

func showStatsAll(database *db.DB) {
	startTime := db.GetStartOfAllTime()
	endTime := time.Now().Unix()

	logs, err := database.GetLogsInRange(startTime, endTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching logs: %v\n", err)
		os.Exit(1)
	}

	if len(logs) == 0 {
		fmt.Println("No data available")
		return
	}

	summary := stats.ComputeSummary(logs)

	fmt.Println("Stats for all time")
	fmt.Println()
	fmt.Println("Overall Totals:")
	fmt.Printf("  Downloaded: %s\n", stats.FormatBytes(summary.TotalBytesIn))
	fmt.Printf("  Uploaded:   %s\n", stats.FormatBytes(summary.TotalBytesOut))
	fmt.Printf("  Total:      %s\n", stats.FormatBytes(summary.TotalBytesIn+summary.TotalBytesOut))
	fmt.Println()
	fmt.Printf("Peak Down:  %s\n", stats.FormatBytesPerSec(summary.PeakBytesIn))
	fmt.Printf("Peak Up:    %s\n", stats.FormatBytesPerSec(summary.PeakBytesOut))
}

func showStatsInterfaces(database *db.DB) {
	startTime := db.GetStartOfDay()
	endTime := time.Now().Unix()

	logsByInterface, err := database.GetLogsByInterface(startTime, endTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching logs: %v\n", err)
		os.Exit(1)
	}

	if len(logsByInterface) == 0 {
		fmt.Println("No data available for today")
		return
	}

	summaries := stats.ComputeByInterface(logsByInterface)

	// Calculate overall totals
	var totalIn, totalOut uint64
	for _, summary := range summaries {
		totalIn += summary.TotalBytesIn
		totalOut += summary.TotalBytesOut
	}

	fmt.Println("Stats by interface (today)")
	fmt.Println()
	fmt.Println("Overall Totals:")
	fmt.Printf("  Downloaded: %s\n", stats.FormatBytes(totalIn))
	fmt.Printf("  Uploaded:   %s\n", stats.FormatBytes(totalOut))
	fmt.Printf("  Total:      %s\n", stats.FormatBytes(totalIn+totalOut))
	fmt.Println()
	fmt.Printf("%-20s %-15s %-15s %-15s\n", "Interface", "Downloaded", "Uploaded", "Total")
	fmt.Println("-------------------------------------------------------------------")

	for _, summary := range summaries {
		total := summary.TotalBytesIn + summary.TotalBytesOut
		fmt.Printf("%-20s %-15s %-15s %-15s\n",
			summary.Interface,
			stats.FormatBytes(summary.TotalBytesIn),
			stats.FormatBytes(summary.TotalBytesOut),
			stats.FormatBytes(total))
	}
}

func showStatsApps(database *db.DB) {
	startTime := db.GetStartOfDay()
	endTime := time.Now().Unix()

	logsByApp, err := database.GetAppLogsByName(startTime, endTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching app logs: %v\n", err)
		os.Exit(1)
	}

	if len(logsByApp) == 0 {
		fmt.Println("No application data available for today")
		fmt.Println("Make sure netmon-service is running")
		return
	}

	summaries := stats.ComputeByApp(logsByApp)

	// Sort by total traffic (downloaded + uploaded)
	sortAppSummaries(summaries)

	// Calculate overall totals
	var totalIn, totalOut uint64
	for _, summary := range summaries {
		totalIn += summary.TotalBytesIn
		totalOut += summary.TotalBytesOut
	}

	fmt.Println("Stats by application (today)")
	fmt.Println()
	fmt.Println("Overall Totals:")
	fmt.Printf("  Downloaded: %s\n", stats.FormatBytes(totalIn))
	fmt.Printf("  Uploaded:   %s\n", stats.FormatBytes(totalOut))
	fmt.Printf("  Total:      %s\n", stats.FormatBytes(totalIn+totalOut))
	fmt.Println()
	fmt.Printf("%-30s %-15s %-15s %-15s\n", "Application", "Downloaded", "Uploaded", "Total")
	fmt.Println("------------------------------------------------------------------------")

	for _, summary := range summaries {
		total := summary.TotalBytesIn + summary.TotalBytesOut
		fmt.Printf("%-30s %-15s %-15s %-15s\n",
			summary.AppName,
			stats.FormatBytes(summary.TotalBytesIn),
			stats.FormatBytes(summary.TotalBytesOut),
			stats.FormatBytes(total))
	}
}

func sortAppSummaries(summaries []stats.AppSummary) {
	// Simple bubble sort by total traffic (descending)
	for i := 0; i < len(summaries); i++ {
		for j := i + 1; j < len(summaries); j++ {
			totalI := summaries[i].TotalBytesIn + summaries[i].TotalBytesOut
			totalJ := summaries[j].TotalBytesIn + summaries[j].TotalBytesOut
			if totalJ > totalI {
				summaries[i], summaries[j] = summaries[j], summaries[i]
			}
		}
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  netmon setup              Set up background service (run this first!)")
	fmt.Println("  netmon version            Show version information")
	fmt.Println("  netmon                    Show today's usage by application (default)")
	fmt.Println("  netmon stats              Show today's usage by application (same as above)")
	fmt.Println("  netmon stats apps         Show today's usage by application")
	fmt.Println("  netmon stats today        Show today's total network usage")
	fmt.Println("  netmon stats week         Show this week's total network usage")
	fmt.Println("  netmon stats month        Show this month's total network usage")
	fmt.Println("  netmon stats all          Show all-time total network usage")
	fmt.Println("  netmon stats interfaces   Show today's usage by interface")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -db <path>               Path to SQLite database (default: ~/.netmon/netmon.db)")
}

// showVersion displays version information
func showVersion() {
	fmt.Printf("%s version %s\n", AppName, Version)
	fmt.Println()
	fmt.Println("macOS Network Usage Monitor")
	fmt.Println("Track network usage by interface and application")
	fmt.Println()
	fmt.Printf("Homepage: https://github.com/abcdOfficialzw/netmon\n")
}

func getDefaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./netmon.db"
	}
	return filepath.Join(home, ".netmon", "netmon.db")
}

// handleSetup runs the interactive setup wizard
func handleSetup() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                ║")
	fmt.Println("║                    NETMON SETUP WIZARD                         ║")
	fmt.Println("║                                                                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not determine executable path: %v\n", err)
		os.Exit(1)
	}

	// Resolve symlinks
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not resolve executable path: %v\n", err)
		os.Exit(1)
	}

	// Get the service binary path (assuming it's in the same directory)
	serviceExePath := filepath.Join(filepath.Dir(exePath), "netmon-service")

	// Check if service binary exists
	if _, err := os.Stat(serviceExePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: netmon-service not found at %s\n", serviceExePath)
		fmt.Println("\nMake sure both netmon and netmon-service are in the same directory.")
		os.Exit(1)
	}

	fmt.Printf("Found netmon-service at: %s\n\n", serviceExePath)

	// Check if already installed
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not determine home directory: %v\n", err)
		os.Exit(1)
	}

	plistPath := filepath.Join(home, "Library", "LaunchAgents", "com.netmon.service.plist")
	
	if _, err := os.Stat(plistPath); err == nil {
		// Already installed
		fmt.Println("⚠️  netmon service is already installed!")
		fmt.Println()
		fmt.Print("Do you want to reinstall/update it? (yes/no): ")
		
		if !promptYesNo() {
			fmt.Println("\nSetup cancelled.")
			return
		}
		
		// Unload existing service
		fmt.Println("\nUnloading existing service...")
		cmd := exec.Command("launchctl", "unload", plistPath)
		cmd.Run() // Ignore errors, it might not be loaded
	}

	// Ask user if they want persistent service
	fmt.Println("Do you want netmon-service to run automatically in the background?")
	fmt.Println("This will:")
	fmt.Println("  ✅ Start monitoring on boot")
	fmt.Println("  ✅ Keep running after restarts")
	fmt.Println("  ✅ Track network usage 24/7")
	fmt.Println()
	fmt.Print("Enable background service? (yes/no): ")

	if !promptYesNo() {
		fmt.Println("\nBackground service not enabled.")
		fmt.Println("You can run the service manually with: ./netmon-service")
		return
	}

	// Create LaunchAgents directory if it doesn't exist
	launchAgentsDir := filepath.Join(home, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not create LaunchAgents directory: %v\n", err)
		os.Exit(1)
	}

	// Create plist content
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.netmon.service</string>
    
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
    </array>
    
    <key>RunAtLoad</key>
    <true/>
    
    <key>KeepAlive</key>
    <true/>
    
    <key>StandardOutPath</key>
    <string>/tmp/netmon-service.log</string>
    
    <key>StandardErrorPath</key>
    <string>/tmp/netmon-service.error.log</string>
    
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
</dict>
</plist>
`, serviceExePath)

	// Write plist file
	fmt.Println("\nCreating service configuration...")
	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not create plist file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Created: %s\n", plistPath)

	// Load the service
	fmt.Println("\nLoading and starting service...")
	cmd := exec.Command("launchctl", "load", plistPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading service: %v\n%s\n", err, string(output))
		os.Exit(1)
	}

	// Wait a moment for service to start
	time.Sleep(1 * time.Second)

	// Verify it's running
	cmd = exec.Command("launchctl", "list", "com.netmon.service")
	if err := cmd.Run(); err != nil {
		fmt.Println("⚠️  Service loaded but may not be running properly.")
		fmt.Println("Check logs: tail -f /tmp/netmon-service.log")
	} else {
		fmt.Println("✓ Service loaded and started successfully!")
	}

	// Success message
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                ║")
	fmt.Println("║                    ✓ SETUP COMPLETE!                           ║")
	fmt.Println("║                                                                ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("netmon-service is now running in the background!")
	fmt.Println()
	fmt.Println("What's next:")
	fmt.Println("  • View your network usage: netmon")
	fmt.Println("  • Check service logs:      tail -f /tmp/netmon-service.log")
	fmt.Println("  • View monthly stats:      netmon stats month")
	fmt.Println()
	fmt.Println("Management commands:")
	fmt.Println("  • Stop service:   launchctl stop com.netmon.service")
	fmt.Println("  • Start service:  launchctl start com.netmon.service")
	fmt.Println("  • Uninstall:      launchctl unload ~/Library/LaunchAgents/com.netmon.service.plist")
	fmt.Println()
	fmt.Println("The service will automatically start on boot. Enjoy! 🚀")
}

// promptYesNo prompts the user for a yes/no answer
func promptYesNo() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		
		input = strings.TrimSpace(strings.ToLower(input))
		
		if input == "yes" || input == "y" {
			return true
		}
		if input == "no" || input == "n" {
			return false
		}
		
		fmt.Print("Please enter 'yes' or 'no': ")
	}
}

