package collector

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// ProcessNetInfo represents network activity for a specific process.
type ProcessNetInfo struct {
	PID         int32
	ProcessName string
	AppName     string // User-friendly application name
	Connections int    // Number of active connections
}

// GetActiveProcesses returns information about processes with active network connections.
func GetActiveProcesses() ([]ProcessNetInfo, error) {
	// Get all network connections
	connections, err := net.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("get connections: %w", err)
	}

	// Map to track unique PIDs
	pidMap := make(map[int32]*ProcessNetInfo)

	for _, conn := range connections {
		if conn.Pid == 0 {
			continue
		}

		// If we've already seen this PID, just increment connection count
		if info, exists := pidMap[conn.Pid]; exists {
			info.Connections++
			continue
		}

		// Get process information
		proc, err := process.NewProcess(conn.Pid)
		if err != nil {
			continue // Process may have terminated
		}

		name, err := proc.Name()
		if err != nil {
			continue
		}

		// Try to get a more user-friendly app name
		appName := extractAppName(proc, name)

		pidMap[conn.Pid] = &ProcessNetInfo{
			PID:         conn.Pid,
			ProcessName: name,
			AppName:     appName,
			Connections: 1,
		}
	}

	// Convert map to slice
	result := make([]ProcessNetInfo, 0, len(pidMap))
	for _, info := range pidMap {
		result = append(result, *info)
	}

	return result, nil
}

// extractAppName attempts to get a user-friendly application name.
func extractAppName(proc *process.Process, processName string) string {
	// Try to get the executable path
	exe, err := proc.Exe()
	if err != nil {
		return cleanProcessName(processName)
	}

	// On macOS, apps are in .app bundles like:
	// /Applications/Google Chrome.app/Contents/MacOS/Google Chrome
	if strings.Contains(exe, ".app/") {
		parts := strings.Split(exe, ".app/")
		if len(parts) > 0 {
			appPath := parts[0] + ".app"
			appName := filepath.Base(appPath)
			// Remove .app extension
			return strings.TrimSuffix(appName, ".app")
		}
	}

	return cleanProcessName(processName)
}

// cleanProcessName removes common suffixes and cleans up process names.
func cleanProcessName(name string) string {
	// Remove common helper/daemon suffixes
	name = strings.TrimSuffix(name, "Helper")
	name = strings.TrimSuffix(name, "d") // daemon suffix
	name = strings.TrimSpace(name)

	// Capitalize first letter
	if len(name) > 0 {
		return strings.ToUpper(name[:1]) + name[1:]
	}

	return name
}

// ConnectionMapper tracks the mapping between network interfaces and processes.
type ConnectionMapper struct {
	lastSnapshot map[int32]ProcessNetInfo
}

// NewConnectionMapper creates a new connection mapper.
func NewConnectionMapper() *ConnectionMapper {
	return &ConnectionMapper{
		lastSnapshot: make(map[int32]ProcessNetInfo),
	}
}

// Update refreshes the process-to-connection mapping.
func (cm *ConnectionMapper) Update() error {
	processes, err := GetActiveProcesses()
	if err != nil {
		return err
	}

	// Update snapshot
	newSnapshot := make(map[int32]ProcessNetInfo)
	for _, p := range processes {
		newSnapshot[p.PID] = p
	}

	cm.lastSnapshot = newSnapshot
	return nil
}

// GetActiveApps returns a list of application names currently using the network.
func (cm *ConnectionMapper) GetActiveApps() []string {
	appMap := make(map[string]bool)

	for _, info := range cm.lastSnapshot {
		appMap[info.AppName] = true
	}

	apps := make([]string, 0, len(appMap))
	for app := range appMap {
		apps = append(apps, app)
	}

	return apps
}

