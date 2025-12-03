package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS traffic_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,
    interface TEXT NOT NULL,
    bytes_in INTEGER NOT NULL,
    bytes_out INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS app_traffic_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,
    app_name TEXT NOT NULL,
    bytes_in INTEGER NOT NULL,
    bytes_out INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_timestamp ON traffic_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_interface ON traffic_logs(interface);
CREATE INDEX IF NOT EXISTS idx_app_timestamp ON app_traffic_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_app_name ON app_traffic_logs(app_name);
`

// DB wraps a sql.DB connection with application-specific methods.
type DB struct {
	conn *sql.DB
}

// Open creates a new database connection and runs migrations.
func Open(path string) (*DB, error) {
	// Ensure the directory exists before trying to open the database
	dbDir := filepath.Dir(path)
	if dbDir != "." && dbDir != "" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("create database directory %s: %w", dbDir, err)
		}
	}

	// Check if the directory is writable
	if dbDir != "." && dbDir != "" {
		if info, err := os.Stat(dbDir); err != nil {
			return nil, fmt.Errorf("database directory %s does not exist and could not be created: %w", dbDir, err)
		} else if !info.IsDir() {
			return nil, fmt.Errorf("database path %s exists but is not a directory", dbDir)
		}
	}

	// Check if the database file exists and is readable/writable
	if _, err := os.Stat(path); err == nil {
		// File exists, check if it's writable
		if info, err := os.Stat(path); err == nil {
			if info.Mode().Perm()&0200 == 0 {
				return nil, fmt.Errorf("database file %s exists but is not writable (permissions: %s)", path, info.Mode().String())
			}
		}
	} else if !os.IsNotExist(err) {
		// Some other error checking the file
		return nil, fmt.Errorf("cannot access database file %s: %w", path, err)
	}

	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		// Provide more helpful error message for common SQLite errors
		if err.Error() == "unable to open database file: out of memory (14)" {
			return nil, fmt.Errorf("cannot open database file %s: %w\nHint: This usually means the directory doesn't exist, there are permission issues, or the disk is full. Check that %s is writable.", path, err, dbDir)
		}
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// migrate runs database schema migrations.
func (db *DB) migrate() error {
	_, err := db.conn.Exec(schema)
	return err
}

// TrafficLog represents a single network traffic measurement.
type TrafficLog struct {
	ID        int64
	Timestamp int64
	Interface string
	BytesIn   uint64
	BytesOut  uint64
}

// AppTrafficLog represents network traffic for a specific application.
type AppTrafficLog struct {
	ID        int64
	Timestamp int64
	AppName   string
	BytesIn   uint64
	BytesOut  uint64
}

// InsertTrafficLog inserts a new traffic log entry.
func (db *DB) InsertTrafficLog(log TrafficLog) error {
	query := `INSERT INTO traffic_logs (timestamp, interface, bytes_in, bytes_out) VALUES (?, ?, ?, ?)`
	_, err := db.conn.Exec(query, log.Timestamp, log.Interface, log.BytesIn, log.BytesOut)
	return err
}

// GetLogsInRange retrieves all traffic logs within a time range.
func (db *DB) GetLogsInRange(startTime, endTime int64) ([]TrafficLog, error) {
	query := `SELECT id, timestamp, interface, bytes_in, bytes_out 
	          FROM traffic_logs 
	          WHERE timestamp >= ? AND timestamp <= ? 
	          ORDER BY timestamp ASC`

	rows, err := db.conn.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []TrafficLog
	for rows.Next() {
		var log TrafficLog
		if err := rows.Scan(&log.ID, &log.Timestamp, &log.Interface, &log.BytesIn, &log.BytesOut); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// GetLogsByInterface retrieves traffic logs grouped by interface within a time range.
func (db *DB) GetLogsByInterface(startTime, endTime int64) (map[string][]TrafficLog, error) {
	logs, err := db.GetLogsInRange(startTime, endTime)
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]TrafficLog)
	for _, log := range logs {
		grouped[log.Interface] = append(grouped[log.Interface], log)
	}

	return grouped, nil
}

// GetStartOfDay returns the Unix timestamp for the start of today (midnight).
func GetStartOfDay() int64 {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return midnight.Unix()
}

// GetStartOfWeek returns the Unix timestamp for the start of this week (Monday).
func GetStartOfWeek() int64 {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	daysToMonday := weekday - 1
	monday := now.AddDate(0, 0, -daysToMonday)
	startOfWeek := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
	return startOfWeek.Unix()
}

// GetStartOfMonth returns the Unix timestamp for the start of this month.
func GetStartOfMonth() int64 {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return startOfMonth.Unix()
}

// GetStartOfAllTime returns the Unix timestamp for the earliest possible time (effectively 0).
func GetStartOfAllTime() int64 {
	return 0
}

// InsertAppTrafficLog inserts a new application traffic log entry.
func (db *DB) InsertAppTrafficLog(log AppTrafficLog) error {
	query := `INSERT INTO app_traffic_logs (timestamp, app_name, bytes_in, bytes_out) VALUES (?, ?, ?, ?)`
	_, err := db.conn.Exec(query, log.Timestamp, log.AppName, log.BytesIn, log.BytesOut)
	return err
}

// GetAppLogsInRange retrieves all app traffic logs within a time range.
func (db *DB) GetAppLogsInRange(startTime, endTime int64) ([]AppTrafficLog, error) {
	query := `SELECT id, timestamp, app_name, bytes_in, bytes_out 
	          FROM app_traffic_logs 
	          WHERE timestamp >= ? AND timestamp <= ? 
	          ORDER BY timestamp ASC`

	rows, err := db.conn.Query(query, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []AppTrafficLog
	for rows.Next() {
		var log AppTrafficLog
		if err := rows.Scan(&log.ID, &log.Timestamp, &log.AppName, &log.BytesIn, &log.BytesOut); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// GetAppLogsByName retrieves traffic logs grouped by app name within a time range.
func (db *DB) GetAppLogsByName(startTime, endTime int64) (map[string][]AppTrafficLog, error) {
	logs, err := db.GetAppLogsInRange(startTime, endTime)
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]AppTrafficLog)
	for _, log := range logs {
		grouped[log.AppName] = append(grouped[log.AppName], log)
	}

	return grouped, nil
}

