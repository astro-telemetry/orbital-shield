package main

import (
	"encoding/json"
	"fmt"
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

func Sanitize(input string) string {
	output := strings.ToLower(input)
	output = sanitizeRegex.ReplaceAllString(output, "-")
	return strings.Trim(output, "-")
}

func CollectTelemetry() ([]InterfaceReport, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Unix()
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	interface_reports := make([]InterfaceReport, 0)

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
		interface_reports = append(interface_reports, report)
	}
	return interface_reports, nil
}

func main() {
	// create directory
	os.MkdirAll("telemetry_logs", 0755)

	// collect data
	collect_telemetry, err := CollectTelemetry()
	if err != nil {
		fmt.Printf("Collector Error: %v\n", err)
		os.Exit(1)
	}

	// ensure list isn't empty
	if len(collect_telemetry) == 0 {
		fmt.Println("No telemetry data captured.")
		return
	}

	// marshal to json
	formatted_reports, err := json.MarshalIndent(collect_telemetry, "", "  ")
	if err != nil {
		fmt.Printf("	! Error: %v\n", err)
		return
	}

	// create filename
	host := Sanitize(collect_telemetry[0].HostID)
	ts := collect_telemetry[0].Timestamp
	filename := fmt.Sprintf("telemetry_logs/telemetry_%s_%d.json", host, ts)

	// write file
	err = os.WriteFile(filename, formatted_reports, 0644)
	if err != nil {
		fmt.Printf("Disk Write Error: %v\n", err)
		return
	}

	fmt.Printf("Successfully saved telemetry to: %s\n", filename)

}
