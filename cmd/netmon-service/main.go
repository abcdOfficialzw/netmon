package main

import (
	"flag"
	"log"
	"netmon/internal/collector"
	"netmon/internal/db"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	var dbPath string
	flag.StringVar(&dbPath, "db", getDefaultDBPath(), "Path to SQLite database file")
	flag.Parse()

	log.Println("Starting netmon-service...")
	log.Printf("Database path: %s", dbPath)
	log.Println("Application tracking: enabled")

	// Ensure database directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Open database
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	log.Println("Database initialized successfully")

	// Initialize collectors
	col := collector.NewCollector()
	appCol := collector.NewAppCollector()

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Create ticker for 1-second intervals
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	log.Println("Collection started (1-second intervals)")

	for {
		select {
		case <-ticker.C:
			if err := collectAndStore(col, database); err != nil {
				log.Printf("Collection error: %v", err)
			}

			if err := collectAndStoreApps(appCol, database); err != nil {
				log.Printf("App collection error: %v", err)
			}

		case sig := <-stop:
			log.Printf("Received signal: %v", sig)
			log.Println("Shutting down gracefully...")
			return
		}
	}
}

// collectAndStore collects network stats and stores them in the database.
func collectAndStore(col *collector.Collector, database *db.DB) error {
	deltas, err := col.Collect()
	if err != nil {
		return err
	}

	// First collection returns nil deltas
	if deltas == nil {
		return nil
	}

	for _, delta := range deltas {
		log := db.TrafficLog{
			Timestamp: delta.Timestamp,
			Interface: delta.Interface,
			BytesIn:   delta.BytesIn,
			BytesOut:  delta.BytesOut,
		}

		if err := database.InsertTrafficLog(log); err != nil {
			return err
		}
	}

	return nil
}

// collectAndStoreApps collects per-app network stats and stores them in the database.
func collectAndStoreApps(appCol *collector.AppCollector, database *db.DB) error {
	appDeltas, err := appCol.CollectWithWeighting()
	if err != nil {
		return err
	}

	// First collection returns nil deltas
	if appDeltas == nil {
		return nil
	}

	for _, delta := range appDeltas {
		log := db.AppTrafficLog{
			Timestamp: delta.Timestamp,
			AppName:   delta.AppName,
			BytesIn:   delta.BytesIn,
			BytesOut:  delta.BytesOut,
		}

		if err := database.InsertAppTrafficLog(log); err != nil {
			return err
		}
	}

	return nil
}

// getDefaultDBPath returns the default database path in user's home directory.
func getDefaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./netmon.db"
	}
	return filepath.Join(home, ".netmon", "netmon.db")
}

