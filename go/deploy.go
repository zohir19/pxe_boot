package main

import (
        "bufio"
        "fmt"
        "log"
        "os"
        "os/exec"
)

// InstallPackages installs a list of packages using apt-get
func InstallPackages(packages []string) error {
        for _, pkg := range packages {
                // Construct the install command
                cmd := exec.Command("sudo", "apt", "install", "-y", pkg)

                // Run the command and capture the output
                output, err := cmd.CombinedOutput()
                if err != nil {
                        return fmt.Errorf("failed to install package %s: %v\nOutput: %s", pkg, err, output)
                }
                fmt.Printf("Successfully installed package: %s\n", pkg)
        }
        return nil
}

// createDirectories creates multiple directories from a given list
func createDirectories(dirs []string) error {
        for _, dir := range dirs {
                err := os.MkdirAll(dir, 0755)
                if err != nil {
                        return fmt.Errorf("failed to create directory %s: %v", dir, err)
                }
                fmt.Printf("Directory %s created successfully\n", dir)
        }
        return nil
}

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
        // List of packages to install
        packages := []string{"dnsmasq", "vim", "nfs-kernel-server", "debootstrap", "grub-efi-amd64-signed"}

        // Call the InstallPackages function
        err := InstallPackages(packages)
        if err != nil {
                log.Fatalf("Error: %v", err)
        } else {
                fmt.Println("All packages installed successfully!")
        }

        // List of directories to create
        dirs := []string{"/srv/tftp", "/srv/nfs"} // Add as many directories as needed

        // Call the createDirectories function to create the directories
        err = createDirectories(dirs)
        if err != nil {
                log.Fatalf("Error: %v", err)
        }

        // Prepare the content for /etc/dnsmasq.d/00-header.conf
        confContent := `
port=0
dhcp-hostsfile=/etc/dnsmasq.d/01-test.hosts
interface=enp107s0                # Use the appropriate network interface
dhcp-range=192.168.0.100,192.168.0.150,12h
dhcp-boot=grubnetx64.efi.signed,linuxhint-s20,192.168.0.1
enable-tftp
tftp-root=/srv/tftp
`

        // Open the file in append mode
        confFile, err := os.OpenFile("/etc/dnsmasq.d/00-header.conf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
                log.Fatalf("Failed to open file: %v", err)
        }
        defer confFile.Close()

        // Write the content to the file
        _, err = confFile.WriteString(confContent)
        if err != nil {
                log.Fatalf("Failed to write to file: %v", err)
        }
        fmt.Println("Configuration written to /etc/dnsmasq.d/00-header.conf")

        // Example: Add host (manually or interactively)
        var mac, hostname, ip string
        scanner := bufio.NewScanner(os.Stdin)

        // Collect user input
        fmt.Print("Enter MAC address: ")
        scanner.Scan()
        mac = scanner.Text()

        fmt.Print("Enter hostname: ")
        scanner.Scan()
        hostname = scanner.Text()

        fmt.Print("Enter IP address: ")
        scanner.Scan()
        ip = scanner.Text()

        // Add host to the file
        err = addHost(mac, hostname, ip)
        if err != nil {
                log.Fatalf("Error adding host: %v", err)
        }
}
