package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type InterfaceReport struct {
	Timestamp     int64    `json:"timestamp"`
	InterfaceName string   `json:"interface_name"`
	IsActive      bool     `json:"is_active"`
	MTU           int      `json:"mtu"`
	IPs           []string `json:"ip_list"`
	ErrorLog      string   `json:"error_log,omitempty"`
}

func main() {

	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	interface_reports := make([]InterfaceReport, 0)

	for _, iface := range interfaces {
		// Bitwise check: Is the 'Up' flag set?
		isActive := iface.Flags&net.FlagUp != 0
		report := InterfaceReport{
			Timestamp:     time.Now().Unix(),
			InterfaceName: iface.Name,
			IsActive:      isActive}

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

	formatted_reports, err := json.MarshalIndent(interface_reports, "", "  ")

	if err != nil {
		fmt.Printf("	! Error: %v\n", err)
	} else {
		fmt.Println(string(formatted_reports))
	}
}
