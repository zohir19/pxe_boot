package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "CM/install"
    "CM/fileops"
    "CM/network"
    "CM/config"
)

func main() {
    // List of packages to install
    packages := []string{"dnsmasq", "vim", "nfs-kernel-server", "debootstrap", "grub-efi-amd64-signed", "shim-signed"}

    // Install packages
    err := install.InstallPackages(packages)
    if err != nil {
        log.Fatalf("Error installing packages: %v", err)
    }

    // Create necessary directories
    dirs := []string{"/srv/tftp", "/srv/nfs", "/srv/tftp/grub"}
    err = fileops.CreateDirectories(dirs)
    if err != nil {
        log.Fatalf("Error creating directories: %v", err)
    }
    var ip_range_start, ip_range_end, int_dhcp, server_ip string
	scanner := bufio.NewScanner(os.Stdin)

	// Get DHCP range start
	fmt.Print("Enter DHCP range start: ")
	scanner.Scan()
	ip_range_start = scanner.Text()

	// Get DHCP range end
	fmt.Print("Enter DHCP range end: ")
	scanner.Scan()
	ip_range_end = scanner.Text()
    fmt.Print("Enter DHCP interface: ")
	scanner.Scan()
	int_dhcp = scanner.Text()
    fmt.Print("Enter server IP: ")
	scanner.Scan()
	server_ip = scanner.Text()

	// Write configuration files dynamically using the user input
	dnsmasqConfig := fmt.Sprintf(`
port=0
dhcp-hostsfile=/etc/dnsmasq.d/01-test.hosts
interface=enp107s0
dhcp-range=%s,%s,12h
dhcp-boot=grubnetx64.efi.signed,linuxhint-s20,192.168.0.1
enable-tftp
tftp-root=/srv/tftp
`, ip_range_start, ip_range_end)

	// Print or save the generated config
	fmt.Println("Generated dnsmasq config:")
	fmt.Println(dnsmasqConfig)

	// Optionally, you can write it to a file
	err := os.WriteFile("/etc/dnsmasq.d/00-header.conf", []byte(dnsmasqConfig), 0644)
	if err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
	} else {
		fmt.Println("Configuration written to /etc/dnsmasq.d/00-header.conf")
	}
}
    // Collect user input and add DHCP host entry
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

    err = network.AddHost(mac, hostname, ip)
    if err != nil {
        log.Fatalf("Error adding DHCP host: %v", err)
    }

    // Copy files
    err = fileops.CopyFile("/usr/lib/shim/shimx64.efi.signed", "/srv/tftp/shimx64.efi.signed")
    if err != nil {
        log.Fatalf("Error copying files: %v", err)
    }
    // Run debootstrap command
    err = fileops.RunDebootstrap()
    if err != nil {
      log.Fatalf("Error running debootstrap: %v", err)
   }

    err = fileops.CopyDirectory("/usr/lib/grub/x86_64-efi/", "/srv/tftp/grub/")
    if err != nil {
        log.Fatalf("Error copying directory: %v", err)
    }
    // Restarting dnsmasq service
    err = fileops.RestartService("dnsmasq")
    if err != nil {
        log.Fatalf("Error restarting service: %v", err)
    }
    // Restarting nfs-kernel service
    err = fileops.RestartService("nfs-kernel-server")
    if err != nil {
        log.Fatalf("Error restarting service: %v", err)
    }
    // Call bindMount for each directory pair
    err = fileops.BindMount("/dev", "/srv/nfs/jammy/dev")
    if err != nil {
        fmt.Println(err)
    }
    err = fileops.BindMount("/proc", "/srv/nfs/jammy/proc")
    if err != nil {
        fmt.Println(err)
    }
    err = fileops.BindMount("/sys", "/srv/nfs/jammy/sys")
    if err != nil {
        fmt.Println(err)
    }
	// update the sources list
	sourceslist := `
deb http://archive.ubuntu.com/ubuntu/ jammy main restricted universe multiverse
# deb-src http://archive.ubuntu.com/ubuntu/ jammy main restricted universe multiverse

deb http://archive.ubuntu.com/ubuntu/ jammy-updates main restricted universe multiverse
# deb-src http://archive.ubuntu.com/ubuntu/ jammy-updates main restricted universe multiverse

deb http://archive.ubuntu.com/ubuntu/ jammy-security main restricted universe multiverse
# deb-src http://archive.ubuntu.com/ubuntu/ jammy-security main restricted universe multiverse

deb http://archive.ubuntu.com/ubuntu/ jammy-backports main restricted universe multiverse
# deb-src http://archive.ubuntu.com/ubuntu/ jammy-backports main restricted universe multiverse

deb http://archive.canonical.com/ubuntu/ jammy partner
# deb-src http://archive.canonical.com/ubuntu/ jammy partner
`
	err = config.WriteConfig("/srv/nfs/jammy/etc/apt/sources.list", sourceslist)
    if err != nil {
        log.Fatalf("Error writing sources list config: %v", err)
    }
	err = fileops.CopyFile("initial_setup.sh", "/srv/nfs/jammy/usr/local/bin")
    if err != nil {
        log.Fatalf("Error copying files: %v", err)
    }
    err = fileops.CopyFile("initial_setup.service", "/srv/nfs/jammy/etc/systemd/system")
    if err != nil {
        log.Fatalf("Error copying files: %v", err)
    }
// chroot Part
    rootDir := "/srv/nfs/jammy/"
    // update 
    err = fileops.RunInChroot(rootDir, "apt", "update")
	if err != nil {
		fmt.Println(err)
	}
    // Install packages inside chroot
	packages = []string{
		"linux-image-generic", "vim", "parted", "dosfstools", "rsync", "nfs-common", "grub-pc-lib", "grub-pc-bin",
	}
    for _, pkg := range packages {
		err = fileops.RunInChroot(rootDir, "apt", "install", "-y", pkg)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Set password
	err = fileops.RunInChroot(rootDir, "passwd")
	if err != nil {
		fmt.Println(err)
	}

	// Enable initial setup service
	err = fileops.RunInChroot(rootDir, "systemctl", "enable", "initial_setup.service")
	if err != nil {
		fmt.Println(err)
	}

	// Exit the chroot environment
	fmt.Println("All commands executed in chroot environment.")
    err = fileops.CopyFile("/srv/nfs/jammy/boot/vmlinuz", "/srv/tftp/jammy/vmlinuz")
    if err != nil {
        log.Fatalf("Error copying files: %v", err)
    }
    err = fileops.CopyFile("/srv/nfs/jammy/boot/initrd.img", "/srv/tftp/jammy/initrd.img")
    if err != nil {
        log.Fatalf("Error copying files: %v", err)
    }

}

