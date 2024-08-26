package network

import (
    "fmt"
    "os"
)

// AddHost appends a DHCP host entry to the /etc/dnsmasq.d/01-test.hosts file
func AddHost(mac, hostname, ip string) error {
    confFile, err := os.OpenFile("/etc/dnsmasq.d/01-test.hosts", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("failed to open file: %v", err)
    }
    defer confFile.Close()

    entry := fmt.Sprintf("dhcp-host=%s,%s,%s,3600\n", mac, hostname, ip)
    _, err = confFile.WriteString(entry)
    if err != nil {
        return fmt.Errorf("failed to write to file: %v", err)
    }

    fmt.Printf("Added entry for MAC: %s, Hostname: %s, IP: %s\n", mac, hostname, ip)
    return nil
}
