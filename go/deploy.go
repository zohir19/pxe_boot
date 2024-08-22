package main

import (
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

	// Create the /srv/tftp directory
	err = os.MkdirAll("/srv/tftp", 0755)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
	fmt.Println("Directory /srv/tftp created successfully")

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
}
