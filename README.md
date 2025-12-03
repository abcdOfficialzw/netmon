# netmon - macOS Network Usage Monitor

A lightweight network monitoring solution for macOS that tracks network usage statistics in real-time.

## Components

- **netmon-service** - Background daemon that collects network statistics every second
  - Tracks interface-level traffic
  - Tracks per-application traffic
- **netmon** - CLI tool for viewing network usage statistics
  - Default view: stats by application
  - View stats by time period (today, week)
  - View stats by interface

## Requirements

- macOS (tested on macOS 10.15+)
- Go 1.22 or later

## Installation

### 1. Build the binaries

```bash
# Install dependencies
go mod download

# Build both executables
go build -o bin/netmon-service ./cmd/netmon-service
go build -o bin/netmon ./cmd/netmon

# Optional: Install to system PATH
sudo cp bin/netmon-service /usr/local/bin/
sudo cp bin/netmon /usr/local/bin/
```

### 2. Quick Setup (Recommended)

Run the interactive setup wizard to configure background monitoring:

```bash
./bin/netmon setup
```

This will:
- Ask if you want the service to run automatically
- Create the background service configuration
- Start the service immediately
- Ensure it runs on boot and after restarts

**That's it!** You're ready to use netmon.

### Installation via Homebrew

```bash
# Add the tap
brew tap abcdofficialzw/netmon

# Install netmon
brew install netmon

# Run setup wizard
netmon setup
```

### Updating netmon

If you installed via Homebrew:

```bash
# Update to the latest version
brew upgrade netmon

# Check your current version
netmon version
```

For detailed update instructions, see [UPDATE_GUIDE.md](UPDATE_GUIDE.md).

## Usage

### Option A: Automatic Setup (Easiest)

```bash
# Run the setup wizard
./bin/netmon setup

# Answer 'yes' when prompted
# Service is now running in background!

# Check your stats
./bin/netmon
```

### Option B: Running the Service Manually

```bash
# Run in foreground (for testing)
./bin/netmon-service

# Run with custom database path
./bin/netmon-service -db /path/to/custom.db

# The default database location is ~/.netmon/netmon.db
```

The service will:
- Collect network statistics every second
- Store interface-level data in SQLite database
- Track per-application usage automatically
- Handle graceful shutdown on SIGINT/SIGTERM (Ctrl+C)

### Using the CLI Tool

```bash
# Default view: application statistics (today)
./bin/netmon

# Or explicitly:
./bin/netmon stats apps

# View today's total statistics
./bin/netmon stats today

# View this week's statistics (Monday to now)
./bin/netmon stats week

# View this month's statistics
./bin/netmon stats month

# View all-time statistics
./bin/netmon stats all

# View statistics by network interface (today)
./bin/netmon stats interfaces

# Use custom database path
./bin/netmon -db /path/to/custom.db
```

#### Example Output

```
Stats for today

Overall Totals:
  Downloaded: 437.72 MB
  Uploaded:   126.66 MB
  Total:      564.37 MB

Peak Down:  11.61 MB/s
Peak Up:    632.48 KB/s
```

```
Stats by interface (today)

Overall Totals:
  Downloaded: 438.99 MB
  Uploaded:   127.31 MB
  Total:      566.31 MB

Interface            Downloaded      Uploaded        Total          
-------------------------------------------------------------------
en0                  322.38 MB       10.70 MB        333.08 MB      
lo0                  116.61 MB       116.61 MB       233.22 MB      
```

```
Stats by application (today)

Overall Totals:
  Downloaded: 329.51 MB
  Uploaded:   92.33 MB
  Total:      421.84 MB

Application                    Downloaded      Uploaded        Total          
------------------------------------------------------------------------
Google Chrome                  77.77 MB        22.17 MB        99.94 MB       
Android Studio                 29.49 MB        8.28 MB         37.77 MB       
Cursor                         14.40 MB        4.00 MB         18.40 MB       
Figma                          9.60 MB         2.49 MB         12.10 MB       
```

## Running as a Background Service (launchd)

### Easy Way: Use Setup Command

```bash
./bin/netmon setup
```

The setup wizard handles everything automatically!

### Manual Way: Create launchd Configuration

If you prefer manual setup, create the file `~/Library/LaunchAgents/com.netmon.service.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.netmon.service</string>
    
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/netmon-service</string>
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
```

### 2. Load and start the service

```bash
# Load the service (starts automatically)
launchctl load ~/Library/LaunchAgents/com.netmon.service.plist

# Check if it's running
launchctl list | grep netmon
```

### 3. Managing the service

```bash
# Stop the service
launchctl stop com.netmon.service

# Start the service
launchctl start com.netmon.service

# Unload the service (stops and disables it)
launchctl unload ~/Library/LaunchAgents/com.netmon.service.plist

# View logs
tail -f /tmp/netmon-service.log
tail -f /tmp/netmon-service.error.log
```

## Database Schema

The SQLite database contains two main tables:

```sql
CREATE TABLE traffic_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,
    interface TEXT NOT NULL,
    bytes_in INTEGER NOT NULL,
    bytes_out INTEGER NOT NULL
);

CREATE TABLE app_traffic_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp INTEGER NOT NULL,
    app_name TEXT NOT NULL,
    bytes_in INTEGER NOT NULL,
    bytes_out INTEGER NOT NULL
);
```

**traffic_logs:**
- **timestamp**: Unix timestamp (seconds)
- **interface**: Network interface name (e.g., "en0", "en1")
- **bytes_in**: Bytes received in the last second (delta)
- **bytes_out**: Bytes sent in the last second (delta)

**app_traffic_logs:**
- **timestamp**: Unix timestamp (seconds)
- **app_name**: Application name (e.g., "Google Chrome", "Slack")
- **bytes_in**: Bytes received by this app in the last second (estimated)
- **bytes_out**: Bytes sent by this app in the last second (estimated)

## Architecture

```
netmon/
├── cmd/
│   ├── netmon-service/    # Background daemon
│   └── netmon/            # CLI tool
├── internal/
│   ├── collector/         # Network stats collection
│   ├── db/                # Database operations
│   └── stats/             # Statistics computation
├── go.mod
└── README.md
```

## How It Works

### Interface Tracking
1. **Collection**: Uses `gopsutil` to read network interface counters from the OS
2. **Delta Computation**: Calculates per-second deltas by comparing consecutive readings
3. **Storage**: Stores deltas in SQLite with timestamps
4. **Aggregation**: CLI tool queries the database and computes totals and peaks

### Application Tracking
1. **Process Discovery**: Identifies processes with active network connections
2. **Application Mapping**: Maps processes to user-friendly application names
   - Detects .app bundles on macOS (e.g., "Google Chrome.app" → "Google Chrome")
   - Handles helper processes and daemons
3. **Traffic Attribution**: Distributes interface-level traffic among active applications
   - Uses connection-count weighting for more accurate attribution
   - Apps with more connections receive proportionally more traffic
4. **Note**: Per-app traffic is estimated using heuristics. For 100% accurate tracking, 
   a kernel extension or network extension would be required (needs special entitlements)

## Troubleshooting

### Service won't start

```bash
# Check permissions
ls -la /usr/local/bin/netmon-service

# Check logs
tail -f /tmp/netmon-service.error.log

# Verify database directory exists
ls -la ~/.netmon/
```

### No data showing in CLI

```bash
# Check if service is running
ps aux | grep netmon-service

# Verify database has data
sqlite3 ~/.netmon/netmon.db "SELECT COUNT(*) FROM traffic_logs;"

# Check database path matches
./bin/netmon-service -db ~/.netmon/netmon.db  # service
./bin/netmon -db ~/.netmon/netmon.db stats today  # CLI
```

### Permission denied on network interfaces

The service needs to run with sufficient permissions to read network statistics. On macOS, this typically works without special privileges, but if you encounter issues:

```bash
# Run with sudo (not recommended for production)
sudo ./bin/netmon-service
```

## Performance

- **CPU**: Minimal (~0.1-0.5%)
- **Memory**: ~10-20 MB
- **Disk**: ~1-2 MB per day of data (varies by network activity)
- **I/O**: One SQLite write per second per active interface

## License

MIT

## Author

Built with Go 1.22+ using:
- `modernc.org/sqlite` - Pure Go SQLite driver
- `github.com/shirou/gopsutil` - Cross-platform system statistics

