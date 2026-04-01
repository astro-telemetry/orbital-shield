package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

type InterfaceReport struct {
	HostID        string   `json:"hostname"`
	Timestamp     int64    `json:"timestamp"`
	InterfaceName string   `json:"interface_name"`
	IsActive      bool     `json:"is_active"`
	MTU           int      `json:"mtu"`
	IPs           []string `json:"ip_list"`
	ErrorLog      string   `json:"error_log,omitempty"`
}

var sanitizeRegex = regexp.MustCompile("[^a-z0-9]+")

func sanitize(input string) string {
	output := strings.ToLower(input)
	output = sanitizeRegex.ReplaceAllString(output, "-")
	return strings.Trim(output, "-")
}

func collectTelemetry() ([]InterfaceReport, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Unix()
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	interfaceReports := make([]InterfaceReport, 0)

	for _, iface := range interfaces {
		// Bitwise check: Is the 'Up' flag set?
		isActive := iface.Flags&net.FlagUp != 0
		report := InterfaceReport{
			HostID:        hostname,
			Timestamp:     timestamp,
			InterfaceName: iface.Name,
			IsActive:      isActive,
		}

		if isActive {
			report.MTU = iface.MTU
			addrs, err := iface.Addrs()
			if err != nil {
				report.ErrorLog = err.Error()
			} else {
				ips := []string{}
				for _, addr := range addrs {
					ips = append(ips, addr.String())
				}
				report.IPs = ips
			}
		}
		interfaceReports = append(interfaceReports, report)
	}
	return interfaceReports, nil
}

func runTelemetryCycle() {
	start := time.Now()
	// collect data
	reports, err := collectTelemetry()
	if err != nil {
		slog.Error("Collector Error", "error_detail", err.Error())
		return
	}

	// ensure list isn't empty
	if len(reports) == 0 {
		slog.Warn("No telemetry data captured")
		return
	}

	// marshal to json
	formatted_reports, err := json.MarshalIndent(reports, "", "  ")
	if err != nil {
		slog.Error("JSON Marshal Error", "error_detail", err.Error())
		return
	}

	// create filename
	host := sanitize(reports[0].HostID)
	ts := reports[0].Timestamp
	filename := fmt.Sprintf("telemetry_logs/telemetry_%s_%d.json", host, ts)

	// write file
	err = os.WriteFile(filename, formatted_reports, 0644)
	if err != nil {
		slog.Error("Disk Write Error", "file", filename, "error_detail", err.Error())
		return
	}

	// Calculate execution time for performance monitoring
	duration := time.Since(start)

	// Output a structured success log with high-value metadata
	slog.Info("Successfully saved telemetry",
		"file", filename,
		"interfaces_scanned", len(reports),
		"duration_ms", duration.Milliseconds(),
	)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Orbital Shield Agent Booting", "version", "1.0.0", "target", "Security ML Engine")
	// create directory
	err := os.MkdirAll("telemetry_logs", 0755)
	if err != nil {
		slog.Error("Failed to create telemetry directory", "error_detail", err.Error())
		os.Exit(1)
	}

	for {
		runTelemetryCycle()

		time.Sleep(60 * time.Second)
	}

}
