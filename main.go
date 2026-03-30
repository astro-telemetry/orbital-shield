package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
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

func CollectTelemetry() ([]InterfaceReport, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	interface_reports := make([]InterfaceReport, 0)
	report := InterfaceReport{
		HostID:    hostname,
		Timestamp: time.Now().Unix(),
	}

	for _, iface := range interfaces {
		// Bitwise check: Is the 'Up' flag set?
		isActive := iface.Flags&net.FlagUp != 0
		report.InterfaceName = iface.Name
		report.IsActive = isActive

		if isActive {
			report.MTU = iface.MTU
			ips := []string{}
			addrs, err := iface.Addrs()
			if err != nil {
				report.ErrorLog = err.Error()
				interface_reports = append(interface_reports, report)
				continue
			}
			for _, addr := range addrs {
				ips = append(ips, addr.String())
			}
			report.IPs = ips
		}
		interface_reports = append(interface_reports, report)
	}
	return interface_reports, err
}

func main() {

	formatted_reports, err := json.MarshalIndent(CollectTelemetry, "", "  ")

	if err != nil {
		fmt.Printf("	! Error: %v\n", err)
	} else {
		fmt.Println(string(formatted_reports))
	}
}
