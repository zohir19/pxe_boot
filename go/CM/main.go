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

    // Write configuration files
    dnsmasqConfig := `
    port=0
    dhcp-hostsfile=/etc/dnsmasq.d/01-test.hosts
    interface=enp107s0
    dhcp-range=192.168.0.100,192.168.0.150,12h
    dhcp-boot=grubnetx64.efi.signed,linuxhint-s20,192.168.0.1
    enable-tftp
    tftp-root=/srv/tftp
    `
    err = config.WriteConfig("/etc/dnsmasq.d/00-header.conf", dnsmasqConfig)
    if err != nil {
        log.Fatalf("Error writing DNSMASQ config: %v", err)
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


}

