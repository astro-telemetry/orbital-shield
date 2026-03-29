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
}

func main() {
	fmt.Println("--- Orbital Shield: Network Status Monitor ---")

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Critical Error: %v\n", err)
		return
	}

	interface_reports := make([]InterfaceReport, 5)

	for _, iface := range interfaces {
		// Bitwise check: Is the 'Up' flag set?
		isActive := iface.Flags&net.FlagUp != 0
		report := InterfaceReport{
			Timestamp:     time.Now().Unix(),
			InterfaceName: iface.Name,
			IsActive:      isActive}

		if isActive {
			report.MTU = iface.MTU
			fmt.Printf("[ACTIVE]  Interface: %-10s (MTU: %d)\n", iface.Name, iface.MTU)
			ips := []string{}
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Printf("  ! Error: %v\n", err)
				continue
			}
			for _, addr := range addrs {
				fmt.Printf("  -> Addr: %s\n", addr.String())
				ips = append(ips, addr.String())
			}
			report.IPs = ips
		} else {
			// Documenting the 'Down' cards
			fmt.Printf("[OFFLINE] Interface: %-10s (Hardware Check Recommended)\n", iface.Name)
		}
		fmt.Println("----------------------------------------------")
		interface_reports = append(interface_reports, report)

	}
	formatted_reports, err := json.MarshalIndent(interface_reports, "", "  ")
	if err != nil {
		fmt.Println(string(formatted_reports))
	} else {
		fmt.Printf("	! Error: %v\n", err)
	}
}
