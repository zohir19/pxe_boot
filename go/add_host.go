package main

import (
        "bufio"
        "fmt"
        "log"
        "os"
)

// addHost appends a DHCP host entry to the /etc/dnsmasq.d/01-test.hosts file
func addHost(mac string, hostname string, ip string) error {
        // Open the file in append mode
        confFile, err := os.OpenFile("/etc/dnsmasq.d/01-test.hosts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
                return fmt.Errorf("failed to open file: %v", err)
        }
        defer confFile.Close()

        // Format the dhcp-host entry
        entry := fmt.Sprintf("dhcp-host=%s,%s,%s,3600\n", mac, hostname, ip)

        // Write the entry to the file
        _, err = confFile.WriteString(entry)
        if err != nil {
                return fmt.Errorf("failed to write to file: %v", err)
        }

        fmt.Printf("Added entry for MAC: %s, Hostname: %s, IP: %s\n", mac, hostname, ip)
        return nil
}

func main() {
        // Example: Manually add a host (this could be interactive or pre-defined)
        var mac, hostname, ip string

        // Collect user input
        scanner := bufio.NewScanner(os.Stdin)
        fmt.Print("Enter MAC address: ")
        scanner.Scan()
        mac = scanner.Text()

        fmt.Print("Enter hostname: ")
        scanner.Scan()
        hostname = scanner.Text()

        fmt.Print("Enter IP address: ")
        scanner.Scan()
        ip = scanner.Text()

        // Call the addHost function to add the entry to the file
        err := addHost(mac, hostname, ip)
        if err != nil {
                log.Fatalf("Error: %v", err)
        }
}
