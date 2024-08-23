package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// InstallPackages installs a list of packages using apt-get
func InstallPackages(packages []string) error {
	for _, pkg := range packages {
		cmd := exec.Command("sudo", "apt", "install", "-y", pkg)
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

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file from %s to %s: %v", src, dst, err)
	}

	fmt.Printf("Successfully copied file from %s to %s\n", src, dst)
	return nil
}

// copyFiles copies multiple files from predefined source to destination paths
func copyFiles() error {
	fileCopies := map[string]string{
		"/usr/lib/grub/x86_64-efi-signed/grubnetx64.efi.signed": "/srv/tftp/grubnetx64.efi.signed",
		"/usr/lib/shim/shimx64.efi.signed":                      "/srv/tftp/shimx64.efi.signed",
	}

	for src, dst := range fileCopies {
		err := copyFile(src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}
func copyDirectory(src, dst string) error {
	cmd := exec.Command("cp", "-r", src, dst)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to copy directory from %s to %s: %v\nOutput: %s", src, dst, err, output)
	}
	fmt.Printf("Successfully copied directory from %s to %s\n", src, dst)
	return nil
}

// runDebootstrap executes the debootstrap command
func runDebootstrap() error {
	cmd := exec.Command("sudo", "debootstrap", "jammy", "/srv/nfs/jammy")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute debootstrap: %v\nOutput: %s", err, output)
	}
	fmt.Println("Debootstrap completed successfully!")
	return nil
}

func main() {
	// List of packages to install
	packages := []string{"dnsmasq", "vim", "nfs-kernel-server", "debootstrap", "grub-efi-amd64-signed", "shim-signed"}

	// Call the InstallPackages function
	err := InstallPackages(packages)
	if err != nil {
		log.Fatalf("Error: %v", err)
	} else {
		fmt.Println("All packages installed successfully!")
	}

	// List of directories to create
	dirs := []string{"/srv/tftp", "/srv/nfs", "/srv/tftp/grub"}

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

	// Open the file and write the content
	confFile, err := os.OpenFile("/etc/dnsmasq.d/00-header.conf", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer confFile.Close()

	_, err = confFile.WriteString(confContent)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
	fmt.Println("Configuration written to /etc/dnsmasq.d/00-header.conf")

	// Prepare the content for /srv/tftp/grub/grub.conf
	confContent1 := `
set timeout=5
timeout_style=menu
#debug=all
set net_default_server=192.168.0.1

menuentry 'DB overlay' {
    linux /jammy/vmlinuz root=/dev/nfs nfsroot=192.168.0.1:/srv/nfs/db_overlay rw BOOTIF=01-$net_default_mac BOOTIP=$net_default_ip console=tty0 console=ttyS0,115200 earlyprintk=ttyS0,115200
    initrd /jammy/initrd.img
}

menuentry 'Ubuntu 22.04' {
    linux /jammy/vmlinuz root=/dev/nfs nfsroot=192.168.0.1:/srv/nfs/jammy rw BOOTIF=01-$net_default_mac BOOTIP=$net_default_ip console=tty0 console=ttyS1,115200 earlyprintk=ttyS1,115200
    initrd /jammy/initrd.img
}
`

	// Open the file and write the content
	confFile, err = os.OpenFile("/srv/tftp/grub/grub.cfg", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer confFile.Close()

	_, err = confFile.WriteString(confContent1)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
	fmt.Println("Configuration written to /srv/tftp/grub/grub.cfg")

	// Collect user input for adding a DHCP host
	var mac, hostname, ip string
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

	// Call the addHost function to add the host entry
	err = addHost(mac, hostname, ip)
	if err != nil {
		log.Fatalf("Error adding host: %v", err)
	}

	// Perform the file copy operations
	err = copyFiles()
	if err != nil {
		log.Fatalf("Error copying files: %v", err)
	}
        err = copyFile("/usr/lib/shim/shimx64.efi.signed","/srv/tftp/shimx64.efi.signed")
	if err != nil {
		log.Fatalf("Error copying files: %v", err)
	}

	// Run debootstrap command
//	err = runDebootstrap()
//	if err != nil {
//		log.Fatalf("Error running debootstrap: %v", err)
//	}
        // Open the file and write the content
	confFile, err := os.OpenFile("/etc/dnsmasq.d/00-header.conf", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer confFile.Close()

	_, err = confFile.WriteString(confContent)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
	fmt.Println("Configuration written to /etc/dnsmasq.d/00-header.conf")
        // Prepare the content for /srv/tftp/grub/grub.conf
	confContent2 := `/srv/nfs/jammy *(rw,sync,no_subtree_check,no_root_squash)`
        // Open the file and write the content
	confFile, err = os.OpenFile("/etc/exports", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer confFile.Close()

	_, err = confFile.WriteString(confContent2)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
	fmt.Println("Configuration written to /etc/exports")
        // Copy the directory /boot/grub/x86_64-efi/ to /srv/tftp/grub/
	err = copyDirectory("/boot/grub/x86_64-efi/", "/srv/tftp/grub/")
	if err != nil {
		log.Fatalf("Error copying directory: %v", err)
	}


}
