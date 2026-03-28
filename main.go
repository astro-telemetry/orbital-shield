package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("--- Orbital Shield: Network Status Monitor ---")

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("Critical Error: %v\n", err)
		return
	}

	for _, iface := range interfaces {
		// Bitwise check: Is the 'Up' flag set?
		isActive := iface.Flags&net.FlagUp != 0

		if isActive {
			fmt.Printf("[ACTIVE]  Interface: %-10s (MTU: %d)\n", iface.Name, iface.MTU)

			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Printf("  ! Error: %v\n", err)
				continue
			}
			for _, addr := range addrs {
				fmt.Printf("  -> Addr: %s\n", addr.String())
			}
		} else {
			// Documenting the 'Down' cards
			fmt.Printf("[OFFLINE] Interface: %-10s (Hardware Check Recommended)\n", iface.Name)
		}
		fmt.Println("----------------------------------------------")
	}
}
